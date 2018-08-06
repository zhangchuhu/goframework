package dao

import "strconv"

type CommonCountry struct {
	ID       int64  `gorm:"AUTO_INCREMENT;primary_key;column:ID"`
	Name     string `gorm:"column:name"`
	Alias    string `gorm:"column:alias"`
	TypeID   int32  `gorm:"column:type"`
	ParentID int64  `gorm:"column:parentid"`
	HasChild int32  `gorm:"column:haschild"`
	Path     string `gorm:"column:path"`
}

//var db *gorm.DB
//
//func init() {
//	var err error
//	db, err = gorm.Open("mysql", "bilin_admin:avYsLkYwQ@tcp(58.215.143.9:6307)/Hujiao?charset=utf8&parseTime=True&loc=Local")
//	if err != nil {
//		fmt.Println("failed", err)
//		os.Exit(-1)
//	}
//}

func GetCommonCountry(id int64) (*CommonCountry, error) {
	var ret CommonCountry
	condition := "ID = " + strconv.FormatInt(id, 10)
	db_ := UserDB.Table("COMMON_COUNTRY").First(&ret, condition)
	if db_.RecordNotFound() {
		return nil, nil
	}
	return &ret, db_.Error
}
