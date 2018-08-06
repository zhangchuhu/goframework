package handler

import (
	"bilin/clientcenter"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/pbkdf2"
	"strconv"
	"strings"
	"time"
)

var (
	successHttp             = &HttpError{"success", 0}
	authHttpErr             = &HttpError{"auth failed", 1}
	ratelimitHttpErr        = &HttpError{"rate limited", 2}
	parseFormHttpErr        = &HttpError{"parseform failed", 3}
	daoGetHttpErr           = &HttpError{"dao get failed", 4}
	jsonMarshalHttpErr      = &HttpError{"json marshal failed", 5}
	notSupportMethodHttpErr = &HttpError{"not support method", 6}
	parseURLHttpErr         = &HttpError{"parse url failed", 7}
	daoPutHttpErr           = &HttpError{"dao put failed", 8}
)

type HandlerFuncWithErro func(*gin.Context) *HttpError

type HttpError struct {
	desc string
	code int64
}

func AuthMiddleWare(fn HandlerFuncWithErro) HandlerFuncWithErro {
	return func(c *gin.Context) *HttpError {
		//auth
		if ok, err := auth(c); !ok || err != nil {
			appzaplog.Warn("user auth failed")
			c.JSON(401, gin.H{
				"code": 1,
				"desc": "auth failed",
			})
			return authHttpErr
		}
		return fn(c)
	}
}

func MetricsMiddleWare(uri string, fn HandlerFuncWithErro) func(c *gin.Context) {
	return func(c *gin.Context) {
		var err *HttpError
		defer func(now time.Time) {
			httpmetrics.DefReport(uri, err.code, now, httpmetrics.DefaultSuccessFun)
		}(time.Now())
		err = fn(c)
	}
}

func auth(c *gin.Context) (bool, error) {
	cookie, err := c.Cookie("bilinanchortoken_")
	if err != nil {
		appzaplog.Error("auth Cookie err", zap.Error(err))
		return authToken(c)
	}
	return authCookie(cookie, c)
}

const (
	uidContextKey = "uid"
)

func authToken(c *gin.Context) (bool, error) {
	uid := c.Query("userId")
	if uid == "" {
		appzaplog.Warn("auth no userId")
		return false, fmt.Errorf("no userId")
	}
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		appzaplog.Error("auth ParseInt err", zap.Error(err))
		return false, err
	}
	c.Set(uidContextKey, uidInt)

	token := c.Query("accessToken")
	if token == "" {
		appzaplog.Warn("auth no accessToken")
		return false, fmt.Errorf("auth no accessToken")
	}
	return clientcenter.VerifyAccessToken(token, uid)
}

const aeskey = "smkldospdosldcaa"

var iv = []byte{0xA, 1, 0xB, 5, 4, 0xF, 7, 9, 0x17, 3, 1, 6, 8, 0xC, 0xD, 91}

func authCookie(cookie string, c *gin.Context) (bool, error) {
	ciphertext, err := base64.RawURLEncoding.DecodeString(cookie)
	if err != nil {
		appzaplog.Error("authCookie base64decode err", zap.Error(err), zap.String("cookie", cookie))
		return false, err
	}
	origData, err := decrpt(aeskey, ciphertext)
	if err != nil {
		appzaplog.Error("authCookie decrpt err", zap.Error(err), zap.String("cookie", cookie))
		return false, err
	}
	info := strings.Split(string(origData), "_")
	if len(info) >= 2 {
		uid, err := strconv.ParseInt(info[0], 10, 64)
		if err != nil {
			appzaplog.Error("authCookie ParseInt err", zap.Error(err))
			return false, err
		}
		expired, err := strconv.ParseInt(info[1], 10, 64)
		if err != nil {
			return false, err
		}
		if time.Now().Unix() > expired {
			return false, nil
		}
		c.Set(uidContextKey, uid)
		return true, nil
	} else {
		return false, nil
	}
}

func getKeyBytes(key string) []byte {
	salt := []byte{0, 7, 2, 3, 4, 5, 6, 7, 8, 1, 0xA, 0xB, 0xE, 0xD, 0xE, 0xF}
	dk := pbkdf2.Key([]byte(key), salt, 10000, 16, sha1.New)
	return []byte(dk)
}

func decrpt(key string, crypted []byte) ([]byte, error) {
	iv := []byte{0xA, 1, 0xB, 5, 4, 0xF, 7, 9, 0x17, 3, 1, 6, 8, 0xC, 0xD, 91}
	keyBytes := getKeyBytes(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroUnPadding(origData []byte) []byte {
	return origData
	//return bytes.TrimFunc(origData,
	//	func(r rune) bool {
	//		return r == rune(0)
	//	})
}
