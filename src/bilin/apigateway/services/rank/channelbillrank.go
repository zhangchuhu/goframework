package rank

import (
	"bilin/thrift/gen-go/rank"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"context"
	"errors"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"time"
)

func GetChannelBillRankList() (error, []*rank.TRank) {
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
	code := GetChannelBillRankCode()
	rank_list, err := client.QueryRank(ctx, code, "incr", "day", true, 100, 0)
	if err != nil {
		connectionPool.ReportErrorConnection(connection)
		return err, nil
	}
	connectionPool.ReturnConnection(connection)
	return nil, rank_list
}

func GetChannelBillRankCode() string {
	t := time.Now()
	minute := t.Minute()
	hour := t.Hour()
	minute = minute - minute%5
	return fmt.Sprintf("bilinChannel5MinuteBillboard_%02d%02d", hour, minute)
}

func GetChannelBillRankInfo() (*map[uint64]ChanlBillRank, error) {

	err, cont_rank_list := GetChannelBillRankList()
	if err != nil {
		appzaplog.Error("GetChannelBillRankList err", zap.Error(err))
		return nil, err
	}

	rank_map := make(map[uint64]ChanlBillRank)
	for i := 0; i < len(cont_rank_list); i++ {
		rank_map[uint64(cont_rank_list[i].UID)] = ChanlBillRank{
			UID:   uint64(cont_rank_list[i].UID),
			Value: cont_rank_list[i].Value,
			Rank:  cont_rank_list[i].Rank,
		}
	}

	return &rank_map, nil
}
