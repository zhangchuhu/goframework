package dao

import (
	"bilin/official/service"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Contract struct {
	gorm.Model
	GuildID              int64 `gorm:"not null" sql:"index"`
	HostUid              int64 `gorm:"not null" sql:"index"`
	ContractStartTime    time.Time
	ContractEndTime      time.Time
	GuildSharePercentage int64
	HostSharePercentage  int64
	ContractState        int32
}

func (Contract) TableName() string {
	return "contract"
}

const contracttableoptions = "ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='签约管理'"

func (c *Contract) Create() error {
	if !hujiaoRecDB.HasTable(&Contract{}) {
		if err := hujiaoRecDB.Set("gorm:table_options", contracttableoptions).CreateTable(&Contract{}).Error; err != nil {
			appzaplog.Error("Create Contract Table err", zap.Error(err))
			return err
		}
	}

	if db_ := hujiaoRecDB.Create(c); db_.Error != nil {
		appzaplog.Error("Create Guild err", zap.Error(db_.Error))
		return db_.Error
	}

	return nil
}

func (c *Contract) Get() ([]Contract, error) {
	var ret []Contract
	if err := hujiaoRecDB.Where(c).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return ret, nil
		}
		appzaplog.Error("Contract.Get err", zap.Error(err), zap.Any("filter", c))
		return ret, err
	}
	return ret, nil
}

func (c *Contract) Update() error {
	if err := hujiaoRecDB.Model(c).Updates(c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil
		}
		appzaplog.Error("Contract.Updates err", zap.Error(err), zap.Any("filter", c))
		return err
	}
	return nil
}

func (c *Contract) Delete() error {
	if err := hujiaoRecDB.Where(c).Delete(c).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return nil
		}
		appzaplog.Error("Contract.Delete err", zap.Error(err), zap.Any("filter", c))
		return err
	}
	return nil
}

func putContractToTurnOver(c *Contract) error {
	info, err := GetByGuildID(c.GuildID)
	if err != nil {
		appzaplog.Error("GetByGuildID err", zap.Error(err), zap.Any("contract", c))
		return err
	}
	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()
	ret, err := service.AddContractInfoExternal(ctx, int64(c.HostUid), int64(c.GuildID), int64(info.OW), int32(c.GuildSharePercentage))
	if err != nil {
		appzaplog.Error("AddContractInfoExternal err", zap.Error(err), zap.Any("contract", c))
		return err
	}
	if ret != 1 {
		appzaplog.Warn("AddContractInfoExternal failed", zap.Int32("ret", ret))
		return errors.New("AddContractInfoExternal failed")
	}
	return nil
}

func GetContractByHostUid(hostuid uint64) (*Contract, error) {
	c := &Contract{}
	cond := fmt.Sprintf("host_uid = %d", hostuid)
	if err := hujiaoRecDB.First(c, cond).Error; err != nil {
		appzaplog.Error("Get Guild err", zap.Error(err), zap.Uint64("hostuid", hostuid))
		return c, err
	}
	appzaplog.Debug("Get Guild", zap.Any("resp", c))
	return c, nil
}

func GetContractsByGuildID(guildid uint64) ([]Contract, error) {
	c := []Contract{}
	cond := fmt.Sprintf("guild_id = %d", guildid)
	if err := hujiaoRecDB.Find(&c, cond).Error; err != nil {
		appzaplog.Error("Get Guild err", zap.Error(err), zap.Uint64("guildid", guildid))
		return c, err
	}
	appzaplog.Debug("Get Guild", zap.Any("resp", c))
	return c, nil
}

func GetAllContract() ([]Contract, error) {
	var ret []Contract
	if err := hujiaoRecDB.Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) || IsTableNotExistErr(err) {
			return ret, nil
		}
		appzaplog.Error("GetAllContract err", zap.Error(err))
		return ret, err
	}
	return ret, nil
}

func MigrateContract() error {
	return hujiaoRecDB.AutoMigrate(&Contract{}).Error
}
