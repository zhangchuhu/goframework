package service

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"bilin/relationlist/config"
	"bilin/relationlist/entity"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

const (
	PREFIX_DAILY_RELATION_TABLE  = "daily_relationlist_value_"         //实时更新，每天一张表
	PREFIX_WEEKLY_RELATION_TABLE = "weekly_relationlist_value_by_hido" //hido每天晚上会更新一次
	PREFIX_TOTAL_RELATION_TABLE  = "total_relationlist_value"          //实时更新，一张总表

	PREFIX_MONTHLY_MEDAL_TABLE = "monthly_relation_medals_record_" //勋章流水，每月一张表
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

//每个主播对应的各用户亲密值
func MysqlCreateDailyRelationTable(table string) {
	const prefix = "MysqlCreateDailyRelationTable "
	createSql := `CREATE TABLE IF NOT EXISTS ` + table + `
			(
				owner bigint(20) unsigned NOT NULL COMMENT '主播UID',
				guest_uid bigint(20) unsigned NOT NULL COMMENT '嘉宾uid（麦上用户或者送礼用户）',
				relation_value float unsigned NOT NULL DEFAULT '0' COMMENT '亲密值（上麦+送礼总和）',
				mike_relation bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '上麦亲密值',
				gift_relation float unsigned NOT NULL DEFAULT '0' COMMENT '送礼亲密值',
				lastupdatetime int(11) DEFAULT NULL COMMENT '最后一次更新时间',
				PRIMARY KEY (owner,guest_uid)
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

//亲密值更新时需要更新mysql,每天一张表。等海度定时任务来计算周榜
func MysqlUpdateUserDailyRelationValue(owner uint64, guest_id uint64, mike_relation_val uint64, gift_relation_val float32) (err error) {
	const prefix = "MysqlUpdateUserDailyRelationValue "

	//get current date
	tableName := PREFIX_DAILY_RELATION_TABLE + strings.Replace(time.Now().String()[0:10], "-", "", -1)

	var stm *sql.Stmt
	for retry := 0; retry < 3; retry++ {
		stm, err = mysqldb.Prepare("INSERT into " + tableName + " (owner, guest_uid, relation_value, mike_relation, gift_relation,lastupdatetime) VALUES (?,?,?,?,?,UNIX_TIMESTAMP())" +
			"ON DUPLICATE KEY UPDATE relation_value=relation_value+?, mike_relation=mike_relation+?, gift_relation=gift_relation+?, lastupdatetime=UNIX_TIMESTAMP()")
		if err != nil {
			if retry == 0 {
				//check table if exists
				if _, deserr := mysqldb.Exec("DESCRIBE " + tableName); deserr != nil {
					// MySQL error 1146 is "table does not exist"
					if mErr, ok := deserr.(*mysql.MySQLError); ok && mErr.Number == 1146 {
						MysqlCreateDailyRelationTable(tableName)
						continue
					}
					// Unknown error
					log.Error(prefix+"mysqldb.DESCRIBE failed", zap.Any("owner", owner), zap.Any("err", err))
					return
				}
			}

			log.Error(prefix+"mysqldb.Prepare failed", zap.Any("owner", owner), zap.Any("err", err))
			return
		} else {
			break
		}
	}

	defer stm.Close()
	_, err = stm.Exec(owner, guest_id, float32(mike_relation_val)+gift_relation_val, mike_relation_val, gift_relation_val,
		float32(mike_relation_val)+gift_relation_val, mike_relation_val, gift_relation_val)
	if err != nil {
		log.Error(prefix+"mysqldb.Exec failed", zap.Any("owner", owner), zap.Any("guest_id", guest_id), zap.Any("err", err))
		return
	}

	log.Info(prefix+"done!", zap.Any("owner", owner), zap.Any("guest_id", guest_id), zap.Any("mike_relation_val", mike_relation_val), zap.Any("gift_relation_val", gift_relation_val))
	return
}

//总榜实时更新，每次亲密值改变时需要更新该表的数据
func MysqlUpdateTotalRelationValue(owner uint64, guest_id uint64, mike_relation_val uint64, gift_relation_val float32) (err error) {
	const prefix = "MysqlUpdateTotalRelationValue "

	//get current date
	tableName := PREFIX_TOTAL_RELATION_TABLE

	var stm *sql.Stmt
	stm, err = mysqldb.Prepare("INSERT into " + tableName + " (owner, guest_uid, relation_value, mike_relation, gift_relation, lastupdatetime) VALUES (?,?,?,?,?,UNIX_TIMESTAMP())" +
		"ON DUPLICATE KEY UPDATE relation_value=relation_value+?, mike_relation=mike_relation+?, gift_relation=gift_relation+?,lastupdatetime=UNIX_TIMESTAMP()")
	if err != nil {
		log.Error(prefix+"mysqldb.Prepare failed", zap.Any("owner", owner), zap.Any("err", err))
		return
	}

	defer stm.Close()
	_, err = stm.Exec(owner, guest_id, float32(mike_relation_val)+gift_relation_val, mike_relation_val, gift_relation_val,
		float32(mike_relation_val)+gift_relation_val, mike_relation_val, gift_relation_val)
	if err != nil {
		log.Error(prefix+"mysqldb.Exec failed", zap.Any("owner", owner), zap.Any("guest_id", guest_id), zap.Any("err", err))
		return
	}

	log.Info(prefix+"done!", zap.Any("owner", owner), zap.Any("guest_id", guest_id), zap.Any("mike_relation_val", mike_relation_val), zap.Any("gift_relation_val", gift_relation_val))
	return
}

func MysqlGetDailyStatisticsRelationList(owner uint64) (result *entity.RelationStatistics, err error) {
	const prefix = "MysqlGetDailyStatisticsRelationList "
	result = &entity.RelationStatistics{AnchorInfo: &entity.UserRelationInfo{UserID: owner, RelationVal: 0}}

	tableName := PREFIX_DAILY_RELATION_TABLE + strings.Replace(time.Now().String()[0:10], "-", "", -1)

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select guest_uid,TRUNCATE(relation_value,0) AS rv from "+tableName+"  WHERE owner = ? and TRUNCATE(relation_value,0) > 0 ORDER BY rv DESC,lastupdatetime ASC LIMIT 50", owner)
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err), zap.Any("owner", owner))
		return
	}

	for retRows.Next() {
		item := &entity.UserRelationInfo{}
		err = retRows.Scan(&item.UserID, &item.RelationVal)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		result.RelationList = append(result.RelationList, item)
		result.AnchorInfo.RelationVal += int64(item.RelationVal)
	}

	log.Info(prefix+"done!", zap.Any("owner", owner), zap.Any("result", result))
	return
}

func MysqlGetWeeklyStatisticsRelationList(owner uint64) (result *entity.RelationStatistics, err error) {
	const prefix = "MysqlGetWeeklyStatisticsRelationList "
	result = &entity.RelationStatistics{AnchorInfo: &entity.UserRelationInfo{UserID: owner, RelationVal: 0}}

	tableName := PREFIX_WEEKLY_RELATION_TABLE

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select guest_uid,TRUNCATE(relation_value,0) AS rv from "+tableName+"  WHERE owner = ? and TRUNCATE(relation_value,0) > 0 GROUP BY owner, guest_uid ORDER BY rv DESC,lastupdatetime ASC LIMIT 50", owner)
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err), zap.Any("owner", owner))
		return
	}

	for retRows.Next() {
		item := &entity.UserRelationInfo{}
		err = retRows.Scan(&item.UserID, &item.RelationVal)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		result.RelationList = append(result.RelationList, item)
		result.AnchorInfo.RelationVal += int64(item.RelationVal)
	}

	log.Info(prefix+"done!", zap.Any("owner", owner), zap.Any("result", result))
	return
}

func MysqlGetTotalStatisticsRelationList(owner uint64) (result *entity.RelationStatistics, err error) {
	const prefix = "MysqlGetTotalStatisticsRelationList "
	result = &entity.RelationStatistics{AnchorInfo: &entity.UserRelationInfo{UserID: owner, RelationVal: 0}}

	tableName := PREFIX_TOTAL_RELATION_TABLE

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select guest_uid,TRUNCATE(relation_value,0) AS rv from "+tableName+"  WHERE owner = ? and TRUNCATE(relation_value,0) > 0 ORDER BY rv DESC,lastupdatetime ASC LIMIT 50", owner)
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err), zap.Any("owner", owner))
		return
	}

	for retRows.Next() {
		item := &entity.UserRelationInfo{}
		err = retRows.Scan(&item.UserID, &item.RelationVal)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		result.RelationList = append(result.RelationList, item)
		result.AnchorInfo.RelationVal += int64(item.RelationVal)
	}

	log.Info(prefix+"done!", zap.Any("owner", owner), zap.Any("result", result))
	return
}

//勋章相关, 每月一张表
func MysqlCreateMonthlyMedalTable(table string) {
	const prefix = "MysqlCreateMonthlyMedalTable "
	createSql := `CREATE TABLE IF NOT EXISTS ` + table + `
			(
				date varchar(20) NOT NULL COMMENT '插入日期',
  				owner bigint(20) NOT NULL COMMENT '主播uid',
  				guest_uid bigint(20) NOT NULL COMMENT '嘉宾uid',
  				medalid int(11) DEFAULT NULL COMMENT '勋章id',
  				PRIMARY KEY (date,owner,guest_uid)
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

func MysqlInsertUserMedalInfo(owner uint64, guest_id uint64, medalid int32) (err error) {
	const prefix = "MysqlInsertUserMedalInfo "

	//get current month
	tableName := PREFIX_MONTHLY_MEDAL_TABLE + strings.Replace(time.Now().String()[0:7], "-", "", -1)

	var stm *sql.Stmt
	for retry := 0; retry < 3; retry++ {
		stm, err = mysqldb.Prepare("INSERT into " + tableName + " (date, owner, guest_uid, medalid) VALUES (?,?,?,?)")
		if err != nil {
			if retry == 0 {
				//check table if exists
				if _, deserr := mysqldb.Exec("DESCRIBE " + tableName); deserr != nil {
					// MySQL error 1146 is "table does not exist"
					if mErr, ok := deserr.(*mysql.MySQLError); ok && mErr.Number == 1146 {
						MysqlCreateMonthlyMedalTable(tableName)
						continue
					}
					// Unknown error
					log.Error(prefix+"mysqldb.DESCRIBE failed", zap.Any("owner", owner), zap.Any("err", err))
					return
				}
			}

			log.Error(prefix+"mysqldb.Prepare failed", zap.Any("owner", owner), zap.Any("err", err))
			return
		} else {
			break
		}
	}

	defer stm.Close()
	_, err = stm.Exec(getCurrentDate(), owner, guest_id, medalid)
	if err != nil {
		log.Error(prefix+"mysqldb.Exec failed", zap.Any("owner", owner), zap.Any("guest_id", guest_id), zap.Any("err", err))
		return
	}

	log.Info(prefix+"done!", zap.Any("owner", owner), zap.Any("guest_id", guest_id), zap.Any("medalid", medalid))
	return
}

//获取主播名下所有勋章用户列表
func MysqlGetOwnerMedalsInfo(owner uint64) (result map[uint64]int32, err error) {
	const prefix = "MysqlGetOwnerMedalsInfo "
	tableName := PREFIX_MONTHLY_MEDAL_TABLE + strings.Replace(time.Now().String()[0:7], "-", "", -1)

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select guest_uid, medalid from "+tableName+"  WHERE owner = ? and date = ? ", owner, getCurrentDate())
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err), zap.Any("owner", owner))
		return
	}

	result = make(map[uint64]int32)
	for retRows.Next() {
		var guest_uid uint64
		var medal_id int32
		err = retRows.Scan(&guest_uid, &medal_id)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		result[guest_uid] = medal_id
	}

	log.Info(prefix+"done!", zap.Any("owner", owner), zap.Any("result", result))
	return
}

//初始化勋章配置
func MysqlGetMedalsConfig() (result map[int32]entity.MedalInfo, err error) {
	const prefix = "MysqlGetMedalsConfig "

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select id, medalName, medalUrl from relation_medals_config")
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err))
		return
	}

	result = make(map[int32]entity.MedalInfo)
	for retRows.Next() {
		item := entity.MedalInfo{}
		err = retRows.Scan(&item.Id, &item.MedalName, &item.MedalUrl)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		result[item.Id] = item
	}

	log.Info(prefix+"done!", zap.Any("result", result))
	return
}

//查询hido任务是否完成
func MysqlGetHidoTaskStatus() (status int32, err error) {
	const prefix = "MysqlGetHidoTaskStatus "

	row := mysqldb.QueryRow("select status from hido_task_status WHERE date = ? and opt = \"hidotask\" ", getCurrentDate())
	err = row.Scan(&status)
	if err != nil {
		log.Warn(prefix, zap.Any("mysql scan error", err))
		return -1, err
	}

	log.Info(prefix+"done!", zap.Any("status", status))
	return
}

//查询勋章分配是否完成
func MysqlGetMedalsTaskStatus() (status int32, err error) {
	const prefix = "MysqlGetMedalsTaskStatus "

	row := mysqldb.QueryRow("select status from hido_task_status WHERE date = ? and opt = ? ", getCurrentDate(), "dispatchmedal")
	err = row.Scan(&status)
	if err != nil {
		log.Warn(prefix, zap.Any("mysql scan error", err))
		return -1, err
	}

	log.Info(prefix+"done!", zap.Any("status", status))
	return
}

func MysqlUpdateMedalsTaskStatus(status int32) (err error) {
	const prefix = "MysqlUpdateMedalsTaskStatus "

	var ReplaceStmt *sql.Stmt
	ReplaceStmt, err = mysqldb.Prepare("REPLACE INTO hido_task_status (date, opt, status) " + "VALUES (?,?,?)")
	if err != nil {
		log.Error(prefix, zap.Any("mysqldb.Prepare error :", err))
		return
	}

	defer ReplaceStmt.Close()
	_, err = ReplaceStmt.Exec(getCurrentDate(), "dispatchmedal", status)
	if err != nil {
		log.Error(prefix, zap.Any("err", err), zap.Any("status", status))
		return
	}

	log.Info(prefix+"done!", zap.Any("status", status))
	return
}

func MysqlGetWeeklyOwners() (owners []uint64, err error) {
	const prefix = "MysqlGetWeeklyOwners "

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select DISTINCT(owner) from weekly_relationlist_value_by_hido where date=?", getCurrentDate())
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err))
		return
	}

	for retRows.Next() {
		var owner uint64
		err = retRows.Scan(&owner)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		owners = append(owners, owner)
	}

	log.Info(prefix+"done!", zap.Any("owners", owners))
	return
}

//从数据库中获取前三名的用户要发放勋章
func MysqlGetTop3RelationByOwner(owner uint64) (result []*entity.UserRelationInfo, err error) {
	const prefix = "MysqlGetTop3RelationByOwner "

	var retRows *sql.Rows
	retRows, err = mysqldb.Query("select guest_uid,TRUNCATE(relation_value,0) AS rv from "+PREFIX_WEEKLY_RELATION_TABLE+"  WHERE owner = ? and TRUNCATE(relation_value,0) > 0 and date = ? ORDER BY rv DESC,lastupdatetime ASC limit 3",
		owner, getCurrentDate())
	if err != nil {
		log.Error(prefix, zap.Any("select error :", err), zap.Any("owner", owner))
		return
	}

	for retRows.Next() {
		item := &entity.UserRelationInfo{}
		err = retRows.Scan(&item.UserID, &item.RelationVal)
		if err != nil {
			log.Error(prefix, zap.Any("Scan error :", err))
			continue
		}

		result = append(result, item)
	}

	log.Info(prefix+"done!", zap.Any("result", result))
	return
}
