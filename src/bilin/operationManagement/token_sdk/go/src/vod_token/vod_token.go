/*
 * filename      : vod_token.go
 * descriptor    :
 * author        : chenyuhui
 * create time   : 2014-11-28 17:09
 * modify list   :
 * +----------------+---------------+---------------------------+
 * | date           | who           | modify summary            |
 * +----------------+---------------+---------------------------+
 */

package vod_token

/*
   #cgo CFLAGS: -I../../../c/
   #cgo LDFLAGS: -L../../../lib -lvod_token_c -lycloud_token
   #include "vod_token.h"
*/
import "C"

import (
	"unsafe"
)

const (
	TYPE_LIVE = iota
	TYPE_VOD
)

type TokenInfo struct {
	Appid uint32
	Ttl   uint64
	Uid   uint32

	//点播字段
	Auth    uint32
	Context string

	//直播字段
	Sid               uint32
	Audio_send_expire uint32
	Audio_recv_expire uint32
	Video_send_expire uint32
	Video_recv_expire uint32
	Text_send_expire  uint32
	Text_recv_expire  uint32
}

func GenToken(secretvsn uint16, secretkey string, info *TokenInfo, tokentype int) (token string, succ bool) {
	var cinfo C.struct_TokenInfo
	cinfo.appid = C.uint32_t(info.Appid)
	cinfo.uid = C.uint32_t(info.Uid)
	cinfo.ttl = C.uint64_t(info.Ttl)
	cinfo.auth = C.uint32_t(info.Auth)
	cinfo.context = C.CString(info.Context)
	defer C.free(unsafe.Pointer(cinfo.context))
	cinfo.sid = C.uint32_t(info.Sid)
	cinfo.audio_send_expire = C.uint32_t(info.Audio_send_expire)
	cinfo.audio_recv_expire = C.uint32_t(info.Audio_recv_expire)
	cinfo.video_send_expire = C.uint32_t(info.Video_send_expire)
	cinfo.video_recv_expire = C.uint32_t(info.Video_recv_expire)
	cinfo.text_send_expire = C.uint32_t(info.Text_send_expire)
	cinfo.text_recv_expire = C.uint32_t(info.Text_recv_expire)
	var ctoken C.struct_Token
	ckey := C.CString(secretkey)
	defer C.free(unsafe.Pointer(ckey))
	var ctype C.TokenType
	if tokentype == TYPE_LIVE {
		ctype = C.TYPE_LIVE
	} else {
		ctype = C.TYPE_VOD
	}
	cvsn := C.uint16_t(secretvsn)
	ret := C.GenToken(cvsn, ckey, &cinfo, ctype, &ctoken)
	if ret == 0 {
		succ = false
	} else {
		succ = true
	}
	defer C.free(unsafe.Pointer(ctoken.key))
	token = C.GoStringN(ctoken.key, ctoken.keylen)
	return
}

func ValidateToken(token string, secretkey string) (succ bool) {
	tokenstr := C.CString(token)
	defer C.free(unsafe.Pointer(tokenstr))
	ctoken := C.struct_Token{
		key:    tokenstr,
		keylen: C.int(len(token)),
	}
	ckey := C.CString(secretkey)
	defer C.free(unsafe.Pointer(ckey))
	ret := C.ValidateToken(&ctoken, ckey)
	if ret == 0 {
		succ = false
	} else {
		succ = true
	}
	return
}

func GetProperty(token string) (succ bool, info *TokenInfo) {
	tokenstr := C.CString(token)
	defer C.free(unsafe.Pointer(tokenstr))
	ctoken := C.struct_Token{
		key:    tokenstr,
		keylen: C.int(len(token)),
	}
	var ctokeninfo C.struct_TokenInfo
	ret := C.GetProperty(&ctoken, &ctokeninfo)
	defer C.free(unsafe.Pointer(ctokeninfo.context))
	if ret == 0 {
		succ = false
	} else {
		succ = true
	}
	info = &TokenInfo{
		Appid:             uint32(ctokeninfo.appid),
		Ttl:               uint64(ctokeninfo.ttl),
		Uid:               uint32(ctokeninfo.uid),
		Auth:              uint32(ctokeninfo.auth),
		Context:           C.GoString(ctokeninfo.context),
		Sid:               uint32(ctokeninfo.sid),
		Audio_send_expire: uint32(ctokeninfo.audio_send_expire),
		Audio_recv_expire: uint32(ctokeninfo.audio_recv_expire),
		Video_send_expire: uint32(ctokeninfo.video_send_expire),
		Video_recv_expire: uint32(ctokeninfo.video_recv_expire),
		Text_send_expire:  uint32(ctokeninfo.text_send_expire),
		Text_recv_expire:  uint32(ctokeninfo.text_recv_expire),
	}
	return
}
