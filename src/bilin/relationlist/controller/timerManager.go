package controller

import (
	"bilin/bcserver/bccommon"
	"bilin/relationlist/service"
	log "code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/kissgo/httpmetrics"
	"strings"
	"time"
)

const (
	DAILY_MIKE_RELATION_LIMIT        = 10     //每天通过上麦获取亲密度上限值
	EQUIVALENT_EXCHANGE_BY_MIKE_TIME = 5 * 60 //上麦5分钟，加1亲密度
	HandleSuccess                    = 0
)

func calculationRelationByMike(owners []uint64) {
	now := time.Now().Unix()
	for _, owner := range owners {
		var result map[uint64]int64 // uid --> timestamp
		result, _ = service.RedisGetUserMikeInfoByOwner(owner)
		for guest_uid, timestamp := range result {
			deltaTimestamp := now - timestamp
			deltaRelation := deltaTimestamp / EQUIVALENT_EXCHANGE_BY_MIKE_TIME
			if deltaRelation == 0 {
				continue
			}

			//只要deltaRelation>0，麦上用户都需要更新最后一次扫描的时间戳
			service.RedisUserOnMike(owner, guest_uid, now)

			//先查找当日亲密度是否达到最大值10
			curRelationVal, _ := service.RedisGetDailyRelationByMike(owner, guest_uid)
			if curRelationVal >= DAILY_MIKE_RELATION_LIMIT { //当日亲密度已经达到最大值
				continue
			}

			//计算增量，添加到日榜和总榜
			if curRelationVal+deltaRelation > DAILY_MIKE_RELATION_LIMIT {
				deltaRelation = DAILY_MIKE_RELATION_LIMIT - curRelationVal
			}

			if deltaRelation <= 0 {
				continue
			}

			//写redis，写mysql

			//更新当日亲密度
			service.RedisSetDailyRelationByMike(owner, guest_uid, curRelationVal+deltaRelation)

			//实时更新日榜  先不写redis，直接通过db计算
			//service.RedisAddDailyStaticsticsRelationList(owner, guest_uid, deltaRelation)
			service.MysqlUpdateUserDailyRelationValue(owner, guest_uid, uint64(deltaRelation), 0)

			//实时更新总榜
			{
				//service.RedisAddOwnerTotalRelation(owner, guest_uid, deltaRelation)
				service.MysqlUpdateTotalRelationValue(owner, guest_uid, uint64(deltaRelation), 0)
			}
		}
	}
}

type TimerManager struct {
}

func NewTimerManager() *TimerManager {
	return &TimerManager{}
}

func (this *TimerManager) Start(interval time.Duration) {
	const prefix = "TimerManager Start "

	for {
		log.Info(prefix + "begin")

		retCode := HandleSuccess
		defer func(now time.Time) {
			httpmetrics.DefReport(strings.TrimSpace(prefix), int64(retCode), now, bccommon.SuccessOrFailedFun)
		}(time.Now())

		if time.Now().Unix()%5 == 0 {
			go this.checkMikeUsers()
		}

		if time.Now().Unix()%60 == 0 {
			go this.DispatchMedaltoTop3()
		}

		time.Sleep(interval)
	}
}

//定时检查麦上用户，计算用户与主播的亲密度
func (this *TimerManager) checkMikeUsers() (err error) {
	const prefix = "checkMikeUsers "
	log.Info(prefix + "begin")

	//加锁，同一时间只能一个进程进行任务
	redisLock, err := service.RedisLock(service.RELATION_LIST_LOCK_KEY, service.RELATION_LIST_LOCK_TIME)
	if err != nil {
		log.Error(prefix+"redis.RedisLock", zap.Any("rediskey", service.RELATION_LIST_LOCK_KEY), zap.Any("err", err))
		return
	}
	defer service.RedisUnLock(redisLock, service.RELATION_LIST_LOCK_KEY)

	var cursor uint64 = 0
	var total_room_count int = 0
	for {
		var owners []uint64
		owners, cursor, err = service.RedisScanOwners(cursor)
		if err != nil {
			log.Error(prefix+"redis.RedisScanOwners", zap.Any("err", err), zap.Any("owners", owners))
			break
		}

		//开始计算亲密度  每5分钟=1亲密度   注意这里不要另起协程，因为我们需要锁住这个操作
		calculationRelationByMike(owners)

		total_room_count += len(owners)
		if cursor == 0 {
			log.Info(prefix+"finished loop, waitting next time", zap.Any("cursor", cursor), zap.Any("total_room_count", total_room_count))
			break
		}
	}

	log.Info(prefix + "end")
	return
}

//每分钟查一次，定时检查周榜，给前三名用户分配勋章并存流水
func (this *TimerManager) DispatchMedaltoTop3() {
	const prefix = "DispatchMedaltoTop3 "
	log.Info(prefix + "begin")

	//加锁，同一时间只能一个进程进行任务
	redisLock, err := service.RedisLock(service.DISPATCH_MEDAL_LOCK_KEY, service.DISPATCH_MEDAL_LOCK_TIME)
	if err != nil {
		log.Error(prefix+"redis.RedisLock", zap.Any("rediskey", service.DISPATCH_MEDAL_LOCK_KEY), zap.Any("err", err))
		return
	}
	defer service.RedisUnLock(redisLock, service.DISPATCH_MEDAL_LOCK_KEY)

	//先查询一下周榜数据是否生成
	status, err := service.MysqlGetHidoTaskStatus()
	if status == -1 { //没有生成,等待海度任务完成
		log.Info(prefix + "MysqlGetHidoTaskStatus failed,waitting hido finished task ")
		return
	}

	status, err = service.MysqlGetMedalsTaskStatus()
	if status != -1 { //进行中或者已完成
		log.Info(prefix+"MysqlGetMedalsTaskStatus dispatch medal already ing... ", zap.Any("status", status))
		return
	}

	//写一个状态数据到db中,防止其他进程同时进行任务
	err = service.MysqlUpdateMedalsTaskStatus(0)
	if err != nil {
		log.Error(prefix+"MysqlUpdateMedalsTaskStatus failed", zap.Any("err", err), zap.Any("status", 0))
		return
	}

	// todo 业务处理
	owners, err := service.MysqlGetWeeklyOwners()
	if err != nil {
		log.Error(prefix+"MysqlGetWeeklyOwners failed", zap.Any("err", err))
		return
	}

	dispatchNum := 0
	for _, owner := range owners {
		relationRet, err := service.MysqlGetTop3RelationByOwner(owner)
		if err != nil {
			log.Error(prefix+"MysqlGetTop3RelationByOwner failed", zap.Any("err", err), zap.Any("owner", owner))
			continue
		}

		//发放勋章
		for index, item := range relationRet {
			if index < 3 && item.RelationVal >= 10 {
				service.MysqlInsertUserMedalInfo(owner, item.UserID, int32(index+1))
			}
		}

		dispatchNum += 1
	}

	//业务处理完成之后再更新一次状态
	service.MysqlUpdateMedalsTaskStatus(1)
	if err != nil {
		log.Error(prefix+"MysqlUpdateMedalsTaskStatus failed", zap.Any("err", err), zap.Any("status", 1))
		return
	}

	log.Info(prefix+"end", zap.Any("owner length", len(owners)), zap.Any("dispatchNum", dispatchNum))
	return
}
