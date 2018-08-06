package service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"bilin/bcserver/bccommon"
	"bilin/bcserver/config"
	"bilin/bcserver/domain/entity"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"
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

func MysqlStorageRoomInfo(room *entity.Room) (err error) {
	const prefix = "MysqlStorageRoomInfo "
	var ReplaceSql *sql.Stmt
	ReplaceSql, err = mysqldb.Prepare("REPLACE INTO bc_roomlist (roomid,owner,status,roomtype,linkstatus,title,roomType2,roomCategoryID,roomPendantLevel,hostBilinID,starttime,endtime,autolink,maixuswitch) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	defer ReplaceSql.Close()
	var res sql.Result
	res, err = ReplaceSql.Exec(room.Roomid, room.Owner, room.Status, room.RoomType, room.LinkStatus, room.Title,
		room.RoomType2, room.RoomCategoryID, room.RoomPendantLevel, room.HostBilinID, room.StartTime, room.EndTime, room.AutoLink, room.Maixuswitch)
	if err != nil {
		log.Error(prefix, zap.Any("err", err), zap.Any("room", room))
		return
	}

	log.Info(prefix, zap.Any("res", res), zap.Any("room", room))
	return
}

func MysqlGetRoomInfo(roomid uint64) (room *entity.Room, err error) {
	const prefix = "MysqlGetRoomInfo "

	defer func(now time.Time) {
		httpmetrics.DefReport("MysqlGetRoomInfo", 0, now, bccommon.SuccessOrFailedFun)
	}(time.Now())

	room = &entity.Room{}
	row := mysqldb.QueryRow("select * from bc_roomlist WHERE roomid = ?", roomid)
	err = row.Scan(&room.Roomid,
		&room.Owner,
		&room.Status,
		&room.RoomType,
		&room.LinkStatus,
		&room.Title,
		&room.RoomType2,
		&room.RoomCategoryID,
		&room.RoomPendantLevel,
		&room.HostBilinID,
		&room.StartTime,
		&room.EndTime,
		&room.AutoLink,
		&room.Maixuswitch,
	)
	if err != nil {
		log.Warn(prefix, zap.Any("mysql scan error", err), zap.Any("roomid", roomid))
		return nil, err
	}

	log.Debug(prefix, zap.Any("room", room))
	return
}

//主播开播，关播流水记录
//每天开播流水记录
func MysqlCreateDailyLivingRecordTable(table string) {
	const prefix = "MysqlCreateDailyLivingRecordTable "
	createSql := `CREATE TABLE IF NOT EXISTS ` + table + `
			(
				id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
				roomid BIGINT(20) UNSIGNED NOT NULL COMMENT '房间ID',
				owner BIGINT(20) UNSIGNED NOT NULL COMMENT '主持人信息',
				title VARCHAR(256) COLLATE utf8_unicode_ci NOT NULL COMMENT '房间名称',
				roomType2 INT(11) NOT NULL COMMENT '1: 官频  2: PGC   3:UGC',
				starttime DATETIME NOT NULL COMMENT '房间开始直播时间',
				endtime DATETIME NOT NULL COMMENT '房间结束直播时间',
				livingTime INT(11) NOT NULL DEFAULT '0' COMMENT '直播时长,单位:秒 ,默认为0',
				PRIMARY KEY (id)
			) ENGINE=MYISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci
	`

	smt, err := mysqldb.Prepare(createSql)
	if err != nil {
		log.Error(prefix, zap.Any("err", err))
	}

	defer smt.Close()
	smt.Exec()
	log.Info(prefix+"create table success !", zap.Any("table", table))
}

//主播开播时写流水
func MysqlInsertLivingData(room *entity.Room) (err error) {
	const prefix = "MysqlInsertLivingData "
	//get current date
	tableName := "daily_living_record_" + strings.Replace(time.Now().String()[0:10], "-", "", -1)

	var stm *sql.Stmt
	for retry := 0; retry < 3; retry++ {
		stm, err = mysqldb.Prepare("INSERT " + tableName + " SET roomid=?, owner=?, title=?, roomType2=?, starttime=now()")
		if err != nil {
			if retry == 0 {
				//check table if exists
				if _, deserr := mysqldb.Exec("DESCRIBE " + tableName); deserr != nil {
					// MySQL error 1146 is "table does not exist"
					if mErr, ok := deserr.(*mysql.MySQLError); ok && mErr.Number == 1146 {
						MysqlCreateDailyLivingRecordTable(tableName)
						continue
					}
					// Unknown error
					log.Error(prefix+"mysqldb.DESCRIBE failed", zap.Any("room", room), zap.Any("err", err))
					return
				}
			}

			log.Error(prefix+"mysqldb.Prepare failed", zap.Any("room", room), zap.Any("err", err))
			return
		} else {
			break
		}
	}

	defer stm.Close()

	res, err := stm.Exec(room.Roomid, room.Owner, room.Title, room.RoomType2)
	if err != nil {
		return
	}

	room.UniqueId, err = res.LastInsertId()

	log.Info(prefix+"create table success !", zap.Any("room", room))

	return
}

//主播关播时更新流水，考虑到夸天的情况，如果在当前日期的表里查不到数据，直接插入一条新的数据
func MysqlUpdateLivingData(room *entity.Room) (err error) {
	const prefix = "MysqlUpdateLivingData "

	//get current date
	tableName := "daily_living_record_" + strings.Replace(time.Now().String()[0:10], "-", "", -1)

	var ReplaceStmt *sql.Stmt
	ReplaceStmt, err = mysqldb.Prepare("REPLACE INTO " + tableName + " (id, roomid, owner, title, roomType2, starttime, endtime, livingTime) " +
		"VALUES (?,?,?,?,?,from_unixtime(?), NOW(), ?)")
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	defer ReplaceStmt.Close()
	var res sql.Result
	res, err = ReplaceStmt.Exec(room.UniqueId, room.Roomid, room.Owner, room.Title, room.RoomType2, room.StartTime, time.Now().Unix()-int64(room.StartTime))
	if err != nil {
		log.Error(prefix, zap.Any("err", err), zap.Any("room", room))
		return
	}

	log.Info(prefix, zap.Any("res", res), zap.Any("room", room))
	return
}
