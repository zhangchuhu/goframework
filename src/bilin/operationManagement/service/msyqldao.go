package service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"bilin/bcserver/config"
	"bilin/operationManagement/entity"
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

func MysqlAddVipUser(headgearinfo *entity.HeadgearInfo) (err error) {
	const prefix = "MysqlAddVipUser "
	var ReplaceStmt *sql.Stmt
	ReplaceStmt, err = mysqldb.Prepare("INSERT INTO vip_headgear (uid, headgear_url, effecttime, expiretime, id) " + "VALUES (?,?,?,?,?)")
	defer ReplaceStmt.Close()
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	var res sql.Result
	res, err = ReplaceStmt.Exec(headgearinfo.Uid, headgearinfo.Headgear, headgearinfo.EffectTime, headgearinfo.ExpireTime, headgearinfo.Id)
	if err != nil {
		log.Error(prefix, zap.Any("err", err), zap.Any("headgearinfo", headgearinfo))
		return
	}

	log.Info(prefix, zap.Any("res", res), zap.Any("headgearinfo", headgearinfo))
	return
}

func MysqlUpdateVipUser(headgearinfo *entity.HeadgearInfo) (err error) {
	const prefix = "MysqlUpdateVipUser "
	var ReplaceStmt *sql.Stmt
	ReplaceStmt, err = mysqldb.Prepare("update vip_headgear set headgear_url=?, effecttime=?, expiretime=?, id=? where uid=?")
	defer ReplaceStmt.Close()
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	var res sql.Result
	res, err = ReplaceStmt.Exec(headgearinfo.Headgear, headgearinfo.EffectTime, headgearinfo.ExpireTime, headgearinfo.Id, headgearinfo.Uid)
	if err != nil {
		log.Error(prefix, zap.Any("err", err), zap.Any("headgearinfo", headgearinfo))
		return
	}

	log.Info(prefix, zap.Any("res", res), zap.Any("headgearinfo", headgearinfo))
	return
}

func MysqlReplaceVipUser(headgearinfo *entity.HeadgearInfo) (err error) {
	const prefix = "MysqlReplaceVipUser "
	var ReplaceStmt *sql.Stmt
	ReplaceStmt, err = mysqldb.Prepare("REPLACE INTO vip_headgear (uid, headgear_url, effecttime, expiretime, id) " + "VALUES (?,?,?,?,?)")
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	defer ReplaceStmt.Close()
	var res sql.Result
	res, err = ReplaceStmt.Exec(headgearinfo.Uid, headgearinfo.Headgear, headgearinfo.EffectTime, headgearinfo.ExpireTime, headgearinfo.Id)
	if err != nil {
		log.Error(prefix, zap.Any("err", err), zap.Any("headgearinfo", headgearinfo))
		return
	}

	log.Info(prefix, zap.Any("res", res), zap.Any("headgearinfo", headgearinfo))
	return
}

func MysqlDelVipUser(uid int64) (err error) {
	const prefix = "MysqlDelVipUser "
	var DeleteStmt *sql.Stmt
	DeleteStmt, err = mysqldb.Prepare("delete from vip_headgear where uid=?")
	defer DeleteStmt.Close()
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	var res sql.Result
	res, err = DeleteStmt.Exec(uid)
	if err != nil {
		log.Error(prefix, zap.Any("err", err), zap.Any("uid", uid))
		return
	}

	log.Info(prefix, zap.Any("res", res), zap.Any("uid", uid))
	return
}

func MysqlGetVipUser(uid int64) (headgearInfo *entity.HeadgearInfo, err error) {
	const prefix = "MysqlGetVipUser "

	headgearInfo = &entity.HeadgearInfo{}
	row := mysqldb.QueryRow("SELECT a.uid,b.headgear_url,a.effecttime,a.expiretime,a.id FROM vip_headgear a JOIN headgear_config b ON a.id=b.id WHERE uid = ?", uid)
	err = row.Scan(&headgearInfo.Uid, &headgearInfo.Headgear, &headgearInfo.EffectTime, &headgearInfo.ExpireTime, &headgearInfo.Id)
	if err != nil {
		log.Warn(prefix, zap.Any("mysql scan error", err))
		return nil, err
	}

	log.Debug(prefix, zap.Any("headgearInfo", headgearInfo))
	return
}

func MysqlGetAllVipUsers() (infos []*entity.HeadgearInfo, err error) {
	const prefix = "MysqlGetVipUsers "

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select a.uid,b.headgear_url,a.effecttime,a.expiretime,a.id FROM vip_headgear a JOIN headgear_config b ON a.id=b.id where expiretime >= (NOW()- INTERVAL 1 DAY)") //已经过期一天以上的数据没有必要取出来
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err))
		return
	}

	for retRows.Next() {
		headgearInfo := &entity.HeadgearInfo{}
		err = retRows.Scan(&headgearInfo.Uid, &headgearInfo.Headgear, &headgearInfo.EffectTime, &headgearInfo.ExpireTime, &headgearInfo.Id)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		infos = append(infos, headgearInfo)
	}

	log.Debug(prefix, zap.Any("infos", infos))
	return
}

//param: date: 20180611
func MysqlGetAllLivingRecord(date string) (infos []*entity.LivingRecordInfo, err error) {
	const prefix = "MysqlGetAllLivingRecord "

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select Roomid,Owner,Title,RoomType2,Starttime,Endtime,LivingTime from " + "daily_living_record_" + date + " order by LivingTime desc limit 1000")
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err))
		return
	}

	for retRows.Next() {
		item := &entity.LivingRecordInfo{}
		err = retRows.Scan(&item.Roomid, &item.Owner, &item.Title, &item.RoomType2, &item.Starttime, &item.Endtime, &item.LivingTime)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		infos = append(infos, item)
	}

	log.Debug(prefix, zap.Any("infos", infos))
	return
}
