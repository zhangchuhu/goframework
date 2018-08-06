package rank

import (
	"bilin/apigateway/config"
	"bilin/common/appthrift"
	"bilin/thrift/gen-go/rank"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
)

const (
	CALL_TITMOUT = time.Second * 2
)

//var hosts = "221.228.110.138:6903;"

var connectionPool *appthrift.ConnectionPool

func InitThriftConnentPool(hosts string) {
	connectionPool = appthrift.NewConnectionPool(3, time.Minute*1, time.Second*3, 1000, hosts)
}

//从营收获取排行榜数据
func GetContributeRankList() (error, []*rank.TRank) {
	if connectionPool == nil {
		return errors.New("thrift not init"), nil
	}

	connection, err := connectionPool.GetConnection()
	if err != nil {
		return errors.New("GetConnection error"), nil
	}
	if connection == nil {
		return errors.New("GetConnection nil"), nil
	}

	var client *rank.TRankServiceClient
	protocolFactory := thrift.NewTCompactProtocolFactory()
	client = rank.NewTRankServiceClientFactory(connection.Transport, protocolFactory)
	ctx, cancel := context.WithTimeout(context.Background(), CALL_TITMOUT)
	defer cancel()
	rank_list, err := client.QueryRank(ctx, "bilinDayUserBillboard", "incr", "day", true, 10, 0)
	if err != nil {
		connectionPool.ReportErrorConnection(connection)
		return err, nil
	}
	connectionPool.ReturnConnection(connection)
	return nil, rank_list
}

func GetContributeRankInfo() (*RankInfo, error) {

	rank_info := GetDefaultContributeRankInfo()
	err, cont_rank_list := GetContributeRankList()
	if err != nil {
		appzaplog.Error("GetContributeRank err", zap.Error(err))
		return nil, err
	}

	uids := []uint64{}
	for i := 0; i < len(cont_rank_list) && i < TOP_NUM; i++ {
		uids = append(uids, uint64(cont_rank_list[i].UID))
	}

	users, err := GetRankUserInfo(uids)
	if err != nil {
		appzaplog.Error("GetRankUserInfo err", zap.Error(err))
		return nil, err
	}

	rank_info.Users = users
	return rank_info, nil
}

func GetDefaultContributeRankInfo() *RankInfo {
	rank_info := &RankInfo{}
	//rank_info.TargetURL = "http://" + config.GetAppConfig().RankTargetHost + CONTRIBUTE_RANK_TARGET_URL
	rank_info.TargetURL = config.GetAppConfig().ContributeRankTargetURL
	rank_info.Title = "今日神豪榜"
	rank_info.Icon = CONTRIBUTE_RANK_ICON_URL
	rank_info.FirstBadge = FIRST_CONTRIBUTE_BADGE_URL
	rank_info.SecondBadge = SECOND_CONTRIBUTE_BADGE_URL
	rank_info.ThirdBadge = THIRD_CONTRIBUTE_BADGE_URL
	return rank_info
}
