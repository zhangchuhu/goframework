package service

import (
	"bilin/bizroomcenter/config"
	"bilin/protocol"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
		log.Error(prefix, zap.Any("connect mysql error : dataSourceName", dataSourceName), zap.Any("err", err))
		panic("MysqlInit failed")
	}

	log.Info(prefix, zap.Any("connect mysql success : dataSourceName", dataSourceName))
}

func MysqlSetBizRoomInfo(roominfo *bilin.BizRoomInfo) (err error) {
	const prefix = "MysqlSetBizRoomInfo "
	var ReplaceStmt *sql.Stmt
	ReplaceStmt, err = mysqldb.Prepare("REPLACE INTO bizlockedrooms (roomid, lockstatus, password) " + "VALUES (?,?,?)")
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	defer ReplaceStmt.Close()
	var res sql.Result
	res, err = ReplaceStmt.Exec(roominfo.Roomid, roominfo.Lockstatus, roominfo.Roompwd)
	if err != nil {
		log.Error(prefix, zap.Any("err", err), zap.Any("roominfo", roominfo))
		return
	}

	log.Info(prefix, zap.Any("res", res), zap.Any("roominfo", roominfo))
	return
}

func MysqlDelBizRoomInfo(roomid uint64) (err error) {
	const prefix = "MysqlDelBizRoomInfo "

	var DeleteStmt *sql.Stmt
	DeleteStmt, err = mysqldb.Prepare("delete from bizlockedrooms where roomid=?")
	defer DeleteStmt.Close()
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	var res sql.Result
	res, err = DeleteStmt.Exec(roomid)
	if err != nil {
		log.Error(prefix, zap.Any("err", err), zap.Any("roomid", roomid))
		return
	}

	log.Info(prefix, zap.Any("res", res), zap.Any("roomid", roomid))
	return
}

func MysqlGetBizRoomInfo(roomid uint64) (roominfo *bilin.BizRoomInfo, err error) {
	const prefix = "MysqlGetBizRoomInfo "

	roominfo = &bilin.BizRoomInfo{}
	row := mysqldb.QueryRow("select * from bizlockedrooms WHERE roomid = ?", roomid)
	err = row.Scan(&roominfo.Roomid, &roominfo.Lockstatus, &roominfo.Roompwd)
	if err != nil {
		log.Warn(prefix, zap.Any("mysql scan error", err))
		return nil, err
	}

	log.Info(prefix, zap.Any("roominfo", roominfo))
	return
}

func MysqlGetAllBizRoomInfos() (infos []*bilin.BizRoomInfo, err error) {
	const prefix = "MysqlGetAllVipUsers "

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select * from bizlockedrooms")
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err))
		return
	}

	for retRows.Next() {
		roominfo := &bilin.BizRoomInfo{}
		err = retRows.Scan(&roominfo.Roomid, &roominfo.Lockstatus, &roominfo.Roompwd)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		infos = append(infos, roominfo)
	}

	log.Info(prefix, zap.Any("infos", infos))
	return
}
