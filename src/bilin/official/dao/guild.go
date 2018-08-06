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
	OW        uint64 `gorm:"not null" sql:"index"`
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
	if !hujiaoCallDB.HasTable(&Guild{}) {
		if err := hujiaoCallDB.Set("gorm:table_options", guildtableoptions).CreateTable(&Guild{}).Error; err != nil {
			appzaplog.Error("Create Guild Table err", zap.Error(err))
			return err
		}
	}
	if db_ := hujiaoCallDB.Create(g); db_.Error != nil {
		appzaplog.Error("Create Guild err", zap.Error(db_.Error))
		return db_.Error
	}
	return nil
}

func GetByGuildID(id uint64) (guild *Guild, err error) {
	guild = &Guild{}
	cond := fmt.Sprintf("id = %d", id)
	if err = hujiaoCallDB.First(guild, cond).Error; err != nil {
		appzaplog.Error("Get Guild err", zap.Error(err), zap.Uint64("guildid", id))
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
	if err := hujiaoCallDB.Exec(rawsql).Error; err != nil {
		appzaplog.Error("UpdateByOw Guild err", zap.Error(err), zap.Any("req", g))
		return err
	}
	if err := hujiaoCallDB.First(g, fmt.Sprintf("ow=%d", ow)).Error; err != nil {
		appzaplog.Error("UpdateByOw First err", zap.Error(err), zap.Any("req", g))
		return err
	}
	return nil
}

func GetGuildByOW(ow uint64) (*Guild, error) {
	var g *Guild = &Guild{}
	cond := fmt.Sprintf("ow = %d", ow)
	if err := hujiaoCallDB.First(g, cond).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		appzaplog.Error("Get Guild err", zap.Error(err), zap.Any("req", g))
		return g, err
	}
	appzaplog.Debug("Get Guild", zap.Any("resp", g))
	return g, nil
}
