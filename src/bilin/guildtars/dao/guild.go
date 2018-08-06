package dao

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

type Guild struct {
	gorm.Model
	OW        int64 `gorm:"not null" sql:"index"`
	Title     string
	Mobile    string
	Describle string
	GuildLogo string
}

func (Guild) TableName() string {
	return "guild"
}

const guildtableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='工会管理'"

func (g *Guild) Create() error {
	if !hujiaoRecDB.HasTable(&Guild{}) {
		if err := hujiaoRecDB.Set("gorm:table_options", guildtableoptions).CreateTable(&Guild{}).Error; err != nil {
			appzaplog.Error("Create Guild Table err", zap.Error(err))
			return err
		}
	}
	if db_ := hujiaoRecDB.Create(g); db_.Error != nil {
		appzaplog.Error("Create Guild err", zap.Error(db_.Error))
		return db_.Error
	}
	return nil
}

func (c *Guild) Get() ([]Guild, error) {
	var ret []Guild
	if err := hujiaoRecDB.Where(c).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return ret, nil
		}
		appzaplog.Error("Guild.Get err", zap.Error(err), zap.Any("filter", c))
		return ret, err
	}
	return ret, nil
}

func (c *Guild) Update() error {
	if err := hujiaoRecDB.Model(c).Updates(c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil
		}
		appzaplog.Error("Guild.Updates err", zap.Error(err), zap.Any("filter", c))
		return err
	}
	return nil
}

func (c *Guild) Delete() error {
	if err := hujiaoRecDB.Where(c).Delete(c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil
		}
		appzaplog.Error("Guild.Delete err", zap.Error(err), zap.Any("filter", c))
		return err
	}
	return nil
}

func GetAllGuild() ([]Guild, error) {
	var ret []Guild
	if err := hujiaoRecDB.Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return ret, nil
		}
		appzaplog.Error("GetAllGuild err", zap.Error(err))
		return ret, err
	}
	return ret, nil
}

func GetByGuildID(id int64) (guild *Guild, err error) {
	guild = &Guild{}
	cond := fmt.Sprintf("id = %d", id)
	if err = hujiaoRecDB.First(guild, cond).Error; err != nil {
		appzaplog.Error("Get Guild err", zap.Error(err), zap.Int64("guildid", id))
		return
	}
	appzaplog.Debug("Get Guild", zap.Any("resp", guild))
	return
}

func (g *Guild) UpdateByOw(ow uint64) error {
	var fields []string
	if g.Mobile != "" {
		fields = append(fields, fmt.Sprintf("mobile=\"%s\"", g.Mobile))
	}
	if g.GuildLogo != "" {
		fields = append(fields, fmt.Sprintf("guild_logo=\"%s\"", g.GuildLogo))
	}
	if g.Title != "" {
		fields = append(fields, fmt.Sprintf("title=\"%s\"", g.Title))
	}
	if g.Describle != "" {
		fields = append(fields, fmt.Sprintf("describle=\"%s\"", g.Describle))
	}
	rawsql := fmt.Sprintf("UPDATE guild SET %s where ow=%d",
		strings.Join(fields, ","), ow,
	)
	appzaplog.Info("UpdateByOw", zap.String("rawsql", rawsql))
	if err := hujiaoRecDB.Exec(rawsql).Error; err != nil {
		appzaplog.Error("UpdateByOw Guild err", zap.Error(err), zap.Any("req", g))
		return err
	}
	if err := hujiaoRecDB.First(g, fmt.Sprintf("ow=%d", ow)).Error; err != nil {
		appzaplog.Error("UpdateByOw First err", zap.Error(err), zap.Any("req", g))
		return err
	}
	return nil
}
