/*
用户徽章
*/
package dao

import (
	"strconv"
	"time"
)

type UserBadge struct {
	UserID   uint64 `gorm:"not null;column:userid" sql:"DEFAULT:0`
	BadgeUrl string `gorm:"not null;column:badgeurl" sql:"DEFAULT:''"`
}

func GetUserBadges() ([]UserBadge, error) {
	nowStr := strconv.FormatInt(time.Now().Unix()*1000, 10)
	filter := " and USER_LABEL.START_TIME < " + nowStr + " and USER_LABEL.END_TIME >" + nowStr
	condition := "select USER_LABEL_RELATION.USER_ID,USER_LABEL.ICON_URL from USER_LABEL_RELATION,USER_LABEL where USER_LABEL_RELATION.LABEL_ID = USER_LABEL.ID"
	db_ := hujiaoUserDB.Raw(condition + filter)
	if db_.Error != nil {
		if db_.RecordNotFound() {
			return []UserBadge{}, nil
		}
		return nil, db_.Error
	}
	rows, err := db_.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []UserBadge
	for rows.Next() {
		var badge UserBadge
		rows.Scan(&badge.UserID, &badge.BadgeUrl)
		ret = append(ret, badge)
	}
	return ret, nil
}

type USER_LABEL struct {
	ID         uint64 `gorm:"not null;primary_key;AUTO_INCREMENT;column:ID"`
	ICON_URL   string `gorm:"not null;column:ICON_URL"`
	NAME       string `gorm:"not null;column:NAME"`
	START_TIME uint64 `gorm:"column:START_TIME"`
	END_TIME   uint64 `gorm:"column:END_TIME"`
}

func (USER_LABEL) TableName() string {
	return "USER_LABEL"
}

type USER_LABEL_RELATION struct {
	ID       uint64 `gorm:"not null;primary_key;AUTO_INCREMENT;column:ID"`
	USER_ID  uint64 `gorm:"not null;column:USER_ID"`
	LABEL_ID uint64 `gorm:"not null;column:LABEL_ID"`
}

func (USER_LABEL_RELATION) TableName() string {
	return "USER_LABEL_RELATION"
}
