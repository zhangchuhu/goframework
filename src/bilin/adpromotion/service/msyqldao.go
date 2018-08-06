package service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"bilin/adpromotion/config"
	"bilin/adpromotion/entity"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
)

var (
	mysqldb *sql.DB
)

func MysqlInit() {
	const prefix = "MysqlInit "
	dataSourceName := config.GetAppConfig().MysqlAddr
	var err error
	mysqldb, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Error(prefix, zap.Any("connect mysql error : dataSourceName", dataSourceName))
		panic("MysqlInit failed")
	}

	log.Info(prefix, zap.Any("connect mysql success : dataSourceName", dataSourceName))
}

func MysqlStorageQihu360Clicks(data map[string][]string, clientIP string) (err error) {
	const prefix = "MysqlStorageQihu360Clicks "
	log.Debug(prefix + "begin")
	var insertSql *sql.Stmt
	insertSql, err = mysqldb.Prepare("INSERT ad_click_events SET UniqueID=?, clicktime=now(), storagetime=now(), IP=?, OS=?, devicetype=?, imei_md5=?, IDFA=?, MAC_MD5=?, callback_url=?, `from`=?")
	defer insertSql.Close()
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	var res sql.Result
	res, err = insertSql.Exec(data["UniqueID"][0], clientIP, data["OS"][0], data["devicetype"][0], data["imei_md5"][0], "111", data["MAC_MD5"][0], data["callback_url"][0], "qihu360")
	if err != nil {
		log.Error(prefix, zap.Any("err", err))
		return
	}

	log.Info(prefix, zap.Any("res", res), zap.Any("UniqueID", data["UniqueID"][0]))
	return
}

func UpdateQihu360CallBackResult(clickInfo *entity.ClickInfo, response string) (err error) {
	const prefix = "UpdateQihu360CallBackResult "
	log.Debug(prefix + "begin")

	//更新数据
	var updateSql *sql.Stmt
	updateSql, err = mysqldb.Prepare("update ad_click_events set callback_rsp=?, callback_time=now() where UniqueID=?")
	defer updateSql.Close()
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	var res sql.Result
	res, err = updateSql.Exec(response, clickInfo.UniqueID)
	if err != nil {
		log.Error(prefix, zap.Any("err", err))
		return
	}

	log.Info(prefix+"success", zap.Any("response", response), zap.Any("UniqueID", clickInfo.UniqueID), zap.Any("res", res))
	return
}

func MysqlSelectClickInfoByImei(imei string) (clickInfo *entity.ClickInfo, err error) {
	const prefix = "MysqlSelectClickInfoByImei "

	clickInfo = &entity.ClickInfo{}
	row := mysqldb.QueryRow("select UniqueID,clicktime,IP,OS,devicetype,imei_md5,MAC_MD5,callback_url,storagetime from ad_click_events WHERE imei_md5 = ? and callback_rsp is NULL limit 1", imei)
	err = row.Scan(&clickInfo.UniqueID,
		&clickInfo.Clicktime,
		&clickInfo.IP,
		&clickInfo.OS,
		&clickInfo.Devicetype,
		&clickInfo.Imei_md5,
		&clickInfo.MAC_MD5,
		&clickInfo.Callback_url,
		&clickInfo.Storagetime,
	)
	if err != nil {
		log.Warn(prefix, zap.Any("mysql scan error", err), zap.Any("imei", imei))
		return
	}

	log.Debug(prefix, zap.Any("clickInfo", clickInfo))
	return
}
