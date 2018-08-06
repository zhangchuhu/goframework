package dao

import "strings"

type SYS_DICT struct {
	ID         uint64 `gorm:"AUTO_INCREMENT;primary_key;column:ID"`
	DICT_KEY   string `gorm:"not null;column:DICT_KEY"`
	DICT_VALUE string `gorm:"not null;column:DICT_VALUE"`
	DEL_FLAG   int32  `gorm:"column:DEL_FLAG"`
}

func (SYS_DICT) TableName() string {
	return "SYS_DICT"
}
func GetAuditWorld() ([]string, error) {
	var ret SYS_DICT
	db_ := hujiaoDB.First(&ret, "DICT_KEY='AUDIT_WORD'")
	if db_.Error != nil {
		if db_.RecordNotFound() {
			return nil, nil
		}
		return nil, db_.Error
	}

	return strings.Split(ret.DICT_VALUE, ","), nil
}
