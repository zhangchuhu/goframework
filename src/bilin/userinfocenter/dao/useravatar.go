package dao

import (
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"strings"
)

type UserAvatarInfo struct {
	Id       uint64 `gorm:"primary_key;column:ID"`
	UserId   uint64 `gorm:"not null;column:USER_ID"`
	ServerId uint64 `gorm:"not null;column:SERVER_ID"`
	Dir      string `gorm:"not null;column:DIR"`
	FileName string `gorm:"not null;column:FILE_NAME"`
}

type DefaultAvatar struct {
	Id           uint64 `gorm:"primary_key;column:ID"`
	SmallHeadURL string `gorm:"not null;column:HEAD_URL_SMALL"`
	BigHeadURL   string `gorm:"not null;column:HEAD_URL_BIG"`
}

type ImageServer struct {
	Id  uint64 `gorm:"primary_key;column:id"`
	URL string `gorm:"column:URL"`
}

//func init() {
//	var err error
//	UserAvatarDB, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/hujiaoavatarts0?charset=utf8&parseTime=True&loc=Local")
//	if err != nil {
//		fmt.Println("failed", err)
//		os.Exit(-1)
//	}
//}

var (
	DefaultAvatarMap map[uint64]*DefaultAvatar
	ImageServerMap   map[uint64]*ImageServer
)

func InitAvatarCache() error {
	DefaultAvatarMap = make(map[uint64]*DefaultAvatar, 0)
	ImageServerMap = make(map[uint64]*ImageServer, 0)
	avatars, err := GetDefaultAvatar()
	if err != nil {
		log.Error("GetDefaultAvatar fail", zap.Error(err))
	}

	for _, v := range avatars {
		DefaultAvatarMap[v.Id] = v
	}
	log.Info("GetDefaultAvatar success", zap.Any("DefaultAvatarMap", DefaultAvatarMap))

	servers, err := GetImageServer()
	if err != nil {
		log.Error("GetImageServer fail", zap.Error(err))
	}

	for _, v := range servers {
		ImageServerMap[v.Id] = v
	}
	log.Info("GetImageServer success", zap.Any("ImageServerMap", ImageServerMap))

	return nil
}

func GetUserAvatarInfo(avatarId, uid uint64) (*UserAvatarInfo, error) {

	if UserAvatarDB == nil {
		log.Warn("no user avatar database connection available")
		return nil, NoAvailabelDB
	}

	avatar := UserAvatarInfo{}
	index := uid % 100
	table_name := fmt.Sprintf("USER_AVATAR_%d", index)
	condition := "ID = " + strconv.FormatUint(avatarId, 10) + " and USER_ID = " + strconv.FormatUint(uid, 10)
	db_ := UserAvatarDB.Table(table_name).First(&avatar, condition)
	if db_.RecordNotFound() {
		return nil, nil
	}

	if db_.Error != nil {
		log.Error("GetUserInfo fail", zap.Uint64("avatarId", avatarId), zap.Uint64("uid", uid), zap.Error(db_.Error))
		return nil, db_.Error
	}

	return &avatar, nil
}

func GetDefaultAvatar() ([]*DefaultAvatar, error) {

	if UserAvatarDB == nil {
		log.Warn("no user avatar database connection available")
		return nil, NoAvailabelDB
	}

	avatars := make([]*DefaultAvatar, 0)
	db_ := UserAvatarDB.Table("DEFAULT_AVATAR").Find(&avatars)
	if db_.RecordNotFound() {
		return nil, nil
	}

	if db_.Error != nil {
		log.Error("GetDefaultAvatar fail", zap.Error(db_.Error))
		return nil, db_.Error
	}

	log.Info("GetDefaultAvatar", zap.Any("avatars", avatars))
	return avatars, nil
}

func GetImageServer() ([]*ImageServer, error) {

	if UserDB == nil {
		log.Warn("no user avatar database connection available")
		return nil, NoAvailabelDB
	}

	servers := make([]*ImageServer, 0)
	db_ := UserDB.Table("IMAGE_SERVER").Find(&servers)
	if db_.RecordNotFound() {
		return nil, nil
	}

	if db_.Error != nil {
		log.Error("GetImageServer fail", zap.Error(db_.Error))
		return nil, db_.Error
	}

	log.Info("GetImageServer", zap.Any("servers", servers))
	return servers, nil
}

func GetUserAvatatrInfos(index uint64, count uint64) ([]UserAvatarInfo, error) {
	if UserAvatarDB == nil {
		log.Warn("no user avatar database connection available")
		return nil, NoAvailabelDB
	}
	if count > 60 {
		count = 60
	}
	avatars := make([]UserAvatarInfo, 0, count)
	index = index % 100
	table_name := fmt.Sprintf("USER_AVATAR_%d", index)
	db_ := UserAvatarDB.Table(table_name).Limit(count).Where("ID > 0 and USER_ID > 0").Find(&avatars)
	if db_.RecordNotFound() {
		return nil, nil
	}

	if db_.Error != nil {
		if IsTableNotExist(db_.Error) {
			log.Warn("query mysql fail", zap.Error(db_.Error))
			return nil, nil
		}
		log.Error("GetUserAvatatrInfos fail", zap.Uint64("index", index), zap.Error(db_.Error))
		return nil, db_.Error
	}

	return avatars, nil
}

func GetUserAvatarUrl(avatarId, uid uint64) (string, error) {
	avatar_info, err := GetUserAvatarInfo(avatarId, uid)
	if err != nil {
		return "", err
	}

	if avatar_info == nil {
		return "", nil
	}

	url := avatar_info.GetAvatar()

	url = HttpsAvatarURL(url)

	return url, nil
}

type ChainStrings string

func (cs ChainStrings) Replace(old, new string) ChainStrings {
	return ChainStrings(strings.Replace(string(cs), old, new, -1))
}
func (cs ChainStrings) ToString() string {
	return string(cs)
}

func HttpsAvatarURL(url string) string {
	chainstr := ChainStrings(url)
	return chainstr.Replace("http://img.onbilin", "https://img.inbilin").
		Replace("http://img2.hujiaozhuanyi.com/imgs/201108/defaultBoy.png", "https://img.inbilin.com/defaultBoy.png").
		Replace("http://img2.hujiaozhuanyi.com/imgs/201108/defaultGirl.png", "https://img.inbilin.com/defaultGirl.png").
		ToString()
}

func (a UserAvatarInfo) GetAvatar() string {
	if a.Dir == "bs2" {
		id, _ := strconv.Atoi(a.FileName)
		if v, ok := DefaultAvatarMap[uint64(id)]; ok {
			return v.SmallHeadURL
		} else {
			log.Warn("not found bs2 default avatar", zap.String("bs2 id", a.FileName))
			return ""
		}
	} else {
		if v, ok := ImageServerMap[a.ServerId]; ok {
			avatar_url := v.URL + "/" + a.Dir + "/" + a.FileName
			avatar_url = strings.Replace(avatar_url, "http://img-bilin.qiniudn.com/", "http://img.onbilin.com/", -1)
			avatar_url = strings.Replace(avatar_url, "http://7d9kvt.com1.z0.glb.clouddn.com/", "http://img.onbilin.com/", -1)
			avatar_url = tosmall(avatar_url)
			return avatar_url
		} else {
			log.Warn("not found server id for default avatar", zap.Uint64("ServerId", a.ServerId))
			return ""
		}
	}
}

func tosmall(url string) string {
	pos := strings.LastIndex(url, "-big")
	if pos >= 0 {
		return url[0:pos] + "-small"
	}
	return url
}
