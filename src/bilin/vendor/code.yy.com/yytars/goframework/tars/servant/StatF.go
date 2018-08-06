/*
todo replace list by channel
*/

package servant

import (
	"container/list"
	"fmt"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"strings"
	"sync"
	"time"
	"code.yy.com/yytars/goframework/jce/servant/taf"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
)

const (
	MAX_SERVERSTAT_NUM = 10000
	MAX_CLIENTSTAT_NUM = 10000
)

type httpStatInfo struct {
	reqAddr    string
	pattern    string
	statusCode int
	costTime   int64
}

type StatInfo struct {
	Head taf.StatMicMsgHead
	Body taf.StatMicMsgBody
}

type StatFHelper struct {
	Enable     bool // if stat report enable
	lStatInfo  *list.List
	lStatInfoFromServer *list.List

	mStatInfo  map[taf.StatMicMsgHead]taf.StatMicMsgBody
	mStatCount map[taf.StatMicMsgHead]int
	mlock      *sync.Mutex
	comm       *Communicator
	sf         *taf.StatF
	node       string
}

func (s *StatFHelper) Init(comm *Communicator, node string) {
	s.Enable = true
	s.node = node
	s.lStatInfo = list.New()
	s.lStatInfoFromServer = list.New()
	s.mlock = new(sync.Mutex)
	s.mStatInfo = make(map[taf.StatMicMsgHead]taf.StatMicMsgBody)
	s.mStatCount = make(map[taf.StatMicMsgHead]int)
	s.comm = comm
	s.sf = new(taf.StatF)
	s.sf.SetServant(s.comm.GetServantProxy(s.node))
}

func (s *StatFHelper) Run() {
	loop := time.NewTicker(10 * time.Second)
	for range loop.C {
		s.addUpMsg(s.lStatInfo, false)
		s.addUpMsg(s.lStatInfoFromServer, true)
	}
}

func (s *StatFHelper) pushBackMsg(stStatInfo StatInfo, fromServer bool) {
	defer s.mlock.Unlock()
	s.mlock.Lock()
	switch {
	case fromServer:
		if s.lStatInfoFromServer.Len() < MAX_SERVERSTAT_NUM {
			s.lStatInfoFromServer.PushFront(stStatInfo)
		} else {
			appzaplog.Warn("server stat report queue is full")
		}
	default:
		if s.lStatInfo.Len() < MAX_CLIENTSTAT_NUM {
			s.lStatInfo.PushFront(stStatInfo)
		} else {
			appzaplog.Warn("client stat report queue is full")
		}
	}
}

func (s *StatFHelper) addUpMsg(statList *list.List, fromServer bool) {
	defer s.mlock.Unlock()
	s.mlock.Lock()
	var n *list.Element

	for e := statList.Front(); e != nil; e = n {
		statInfo := e.Value.(StatInfo)
		bodyList := statInfo.Body

		if body, ok := s.mStatInfo[statInfo.Head]; ok {
			body.Count += statInfo.Body.Count
			body.TimeoutCount += statInfo.Body.TimeoutCount
			body.ExecCount += statInfo.Body.ExecCount
			body.TotalRspTime += statInfo.Body.TotalRspTime
			body.MaxRspTime += statInfo.Body.MaxRspTime
			body.MinRspTime += statInfo.Body.MinRspTime
			body.WeightValue += statInfo.Body.WeightValue
			body.WeightCount += statInfo.Body.WeightCount
			s.mStatInfo[statInfo.Head] = body
			s.mStatCount[statInfo.Head] += 1
		} else {
			headMap := statInfo.Head
			s.mStatInfo[headMap] = taf.StatMicMsgBody{
				Count:bodyList.TimeoutCount,
				TimeoutCount:bodyList.TimeoutCount,
				ExecCount: bodyList.ExecCount,
				TotalRspTime : bodyList.TotalRspTime,
				MaxRspTime: bodyList.MaxRspTime,
				MinRspTime: bodyList.MinRspTime,
				WeightValue: bodyList.WeightValue,
				WeightCount: bodyList.WeightCount,
			}
			s.mStatCount[headMap] = 1
		}

		n = e.Next()
		statList.Remove(e)
	}

	for k, v := range s.mStatInfo {
		c := int32(s.mStatCount[k])
		v.Count = v.Count / c
		v.TimeoutCount = v.TimeoutCount / c
		v.ExecCount = v.ExecCount / c
		v.TotalRspTime = v.TotalRspTime / int64(c)
		v.MaxRspTime = v.MaxRspTime / c
		v.MinRspTime = v.MinRspTime / c
		v.WeightValue = v.WeightValue / c
		v.WeightCount = v.WeightCount / c
		s.mStatInfo[k] = v
	}

	_, err := s.sf.ReportMicMsg(s.mStatInfo, !fromServer)
	if err != nil {
		appzaplog.Error("report err", zap.Error(err))
	}

	for m := range s.mStatInfo {
		delete(s.mStatInfo, m)
	}
}

func (s *StatFHelper) ReportMicMsg(stStatInfo StatInfo) {
	if s.Enable {
		go s.pushBackMsg(stStatInfo, false)
	}
}

var StatReport *StatFHelper

func initStatF(comm *Communicator,stat string) error {
	if comm == nil || stat == ""{
		appzaplog.Warn("initReport failed, nil comm or stat empty")
		return NilParamsErr
	}
	StatReport = &StatFHelper{
	}
	StatReport.Init(comm,stat)
	go StatReport.Run()

	return nil
}

func ReportStat(msg IMessage, succ int32, timeout int32, exec int32) {
	if StatReport == nil || !StatReport.Enable {
		return
	}
	var head taf.StatMicMsgHead
	if cfg := GetServerConfig(); cfg != nil {
		head.MasterName = fmt.Sprintf("%s.%s", cfg.App, cfg.Server)
		head.MasterIp = cfg.LocalIP
		head.TafVersion = cfg.Version
		head.SMasterContainer = cfg.Container
	} else {
		//TODO
		return
	}
	head.InterfaceName = msg.getFuncName()
	if adp := msg.getAdapterProxy(); adp != nil {
		head.SSlaveContainer = adp.GetPoint().ContainerName
		head.SlaveIp = adp.GetPoint().Host
		head.SlavePort = adp.GetPoint().Port
	}

	sNames := strings.Split(msg.getSServantName(), ".")
	if len(sNames) < 2 {
		appzaplog.Warn("No Stat Server found")
		return
	}
	head.SlaveName = fmt.Sprintf("%s.%s", sNames[0], sNames[1])
	//TODO set Resp
	head.ReturnValue = msg.getRespRet()

	info := StatInfo{
		Head: head,
		Body: taf.StatMicMsgBody{
			Count:        succ,
			TimeoutCount: timeout,
			ExecCount:    exec,
			TotalRspTime: msg.Cost(),
			MaxRspTime:   int32(msg.Cost()),
			MinRspTime:   int32(msg.Cost()),
		},
	}

	StatReport.ReportMicMsg(info)
}

func reportHttpStat(st *httpStatInfo) {
	if StatReport == nil || !StatReport.Enable {
		return
	}
	var cfg *serverConfig
	if cfg = GetServerConfig();cfg == nil{
		return
	}

	var statHead = taf.StatMicMsgHead{
		MasterName:"http_client",
		MasterIp:st.reqAddr,
		TafVersion:cfg.Version,
		SlaveName:fmt.Sprintf("%s.%s", cfg.App, cfg.Server),
		SlaveIp:cfg.LocalIP,
		SSlaveContainer:cfg.Container,
		InterfaceName:st.pattern,
	}

	var statBody = taf.StatMicMsgBody{}
	if st.statusCode >= 400 {
		statBody.ExecCount = 1 // 异常
	} else {
		statBody.Count = 1
		statBody.TotalRspTime = st.costTime
		statBody.MaxRspTime = int32(st.costTime)
		statBody.MinRspTime = int32(st.costTime)
	}

	info := StatInfo{}
	info.Head = statHead
	info.Body = statBody
	StatReport.ReportMicMsg(info)
}
