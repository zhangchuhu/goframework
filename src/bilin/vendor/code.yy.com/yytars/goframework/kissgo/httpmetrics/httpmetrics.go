// @author kordenlu
// @创建时间 2017/05/08 10:20
// 功能描述:

package httpmetrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"sync"
	"io"
	"io/ioutil"
)

const (
	metricsapi         = "http://127.0.0.1:10039/metrics_api"
	ReportNumThreshold = 20 // 20条聚合上报一次
	MaxReportBuf       = 3000
	MaxReportInterval  = 5 * time.Second
)

var (
	//error
	ReportBufFullErr = errors.New("Report too many,drop")
)

type DefModel struct {
	Topic     string  `json:"topic"`
	Uri       string  `json:"uri"`
	UriTag    string  `json:"uri_tag"`
	Duration  int64   `json:"duration"`
	Code      string  `json:"code"`
	IsSuccess string  `json:"isSuccess,omitempty"`
	Scale     []int64 `json:"scale,omitempty"`
}

type Counter struct {
	Topic string `json:"topic"`
	Uri   string `json:"uri"`
	Val   int64  `json:"val"`
}

type Gauge struct {
	Topic string `json:"topic"`
	Uri   string `json:"uri"`
	Val   int64  `json:"val"`
}

type Histo struct {
	Topic string  `json:"topic"`
	Uri   string  `json:"uri"`
	Val   int64   `json:"val"`
	Scale []int64 `json:"scale"`
}

type RetCode struct {
	Topic string `json:"topic"`
	Uri   string `json:"uri"`
	Code  int64  `json:"code"`
}
type MetricsReport struct {
	AppName       string      `json:"app_name"`
	AppVer        string      `json:"app_ver"`
	ServiceName   string      `json:"service_name"`
	Step          int64       `json:"step"`
	Ver           string      `json:"ver"`
	ServerId      int64       `json:"server_id,omitempty"`
	IdcId         int64       `json:"idc_id,omitempty"`
	ISP           string      `json:"isp,omitempty"`
	AreaId        int64       `json:"area_id,omitempty"`
	Skip1stPeriod bool        `json:"skip_1st_period"`
	Defmodel      []*DefModel `json:"defmodel,omitempty"`
	Counter       []*Counter  `json:"counter,omitempty"`
	Gauge         []*Gauge    `json:"gauge,omitempty"`
	Histo         []*Histo    `json:"histo,omitempty"`
	Retcode       []*RetCode  `json:"retcode,omitempty"`
}

type MetricsRepMannager struct {
	Enabled      bool
	DefModelChan chan *DefModel
	CounterChan  chan *Counter
	GaugeChan    chan *Gauge
	HistoChan    chan *Histo
	RetCodeChan  chan *RetCode
	Reports      MetricsReport
}

var (
	once = sync.Once{}
	GMetricsRepMannager  = MetricsRepMannager{
		DefModelChan: make(chan *DefModel, MaxReportBuf),
		CounterChan:  make(chan *Counter, MaxReportBuf),
		GaugeChan:    make(chan *Gauge, MaxReportBuf),
		HistoChan:    make(chan *Histo, MaxReportBuf),
		RetCodeChan:  make(chan *RetCode, MaxReportBuf),
		Reports: MetricsReport{
			AppVer: "1.0",
			Ver:    "0.1",
			Step:   60,
		},
	}
)

func EnableMetrics(appname, svcname string) {
	once.Do(
		func() {
			GMetricsRepMannager.Enabled = true
			GMetricsRepMannager.Reports.AppName = appname
			GMetricsRepMannager.Reports.ServiceName = svcname
			go GMetricsRepMannager.reportloop()
		})
}

type SuccessFun func (code int64)bool

func DefaultSuccessFun(code int64) bool {
	return code == 0
}

func DefReport(uri string, code int64, tbegin time.Time,succfun SuccessFun) error {
	if GMetricsRepMannager.Enabled {
		model := &DefModel{
			Uri:       uri,
			Duration:  int64(time.Since(tbegin)) / 1000,
			Code:      strconv.FormatInt(code, 10),
			UriTag:    "s",
			IsSuccess: isSuccess(succfun(code)),
		}
		return GMetricsRepMannager.defaultreportModel(model)
	}
	return nil
}

func CounterMetric(uri string, val int64) error {
	if GMetricsRepMannager.Enabled {
		model := &Counter{
			Topic: "counter",
			Uri:   uri,
			Val:   val,
		}
		select {
		case GMetricsRepMannager.CounterChan <- model:
			return nil
		default:
			fmt.Printf("CounterChan chan full,drop model:%v\n", model)
			return ReportBufFullErr
		}
	}

	return nil
}

func GaugeMetric(uri string, val int64) error {
	if GMetricsRepMannager.Enabled {
		model := &Gauge{
			Topic: "gauge",
			Uri:   uri,
			Val:   val,
		}
		select {
		case GMetricsRepMannager.GaugeChan <- model:
			return nil
		default:
			fmt.Printf("GaugeChan chan full,drop model:%v\n", model)
			return ReportBufFullErr
		}
	}
	return nil
}

func HistoMetric(uri string, val int64) error {
	if GMetricsRepMannager.Enabled {
		model := &Histo{
			Topic: "histo",
			Uri:   uri,
			Val:   val,
		}
		select {
		case GMetricsRepMannager.HistoChan <- model:
			return nil
		default:
			fmt.Printf("HistoChan chan full,drop model:%v\n", model)
			return ReportBufFullErr
		}
	}

	return nil
}

func RetCodeMetric(uri string, code int64) error {
	if GMetricsRepMannager.Enabled {
		model := &RetCode{
			Topic: "retcode",
			Uri:   uri,
			Code:  code,
		}
		select {
		case GMetricsRepMannager.RetCodeChan <- model:
			return nil
		default:
			fmt.Printf("RetCodeChan chan full,drop model:%v\n", model)
			return ReportBufFullErr
		}
	}

	return nil
}

func (self *MetricsRepMannager) reportloop() {
	ticker := time.NewTicker(MaxReportInterval)
	defer ticker.Stop()

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for {
		if self.overReportThreshold() {
			self.report(client)
		}
		select {
		case elem := <-self.DefModelChan:
			self.Reports.Defmodel = append(self.Reports.Defmodel, elem)
		case elem := <-self.CounterChan:
			self.Reports.Counter = append(self.Reports.Counter, elem)
		case elem := <-self.GaugeChan:
			self.Reports.Gauge = append(self.Reports.Gauge, elem)
		case elem := <-self.HistoChan:
			self.Reports.Histo = append(self.Reports.Histo, elem)
		case elem := <-self.RetCodeChan:
			self.Reports.Retcode = append(self.Reports.Retcode, elem)
		case <-ticker.C:
			self.report(client)
		}
	}
}

func (self *MetricsRepMannager) overReportThreshold() bool {
	return (len(self.Reports.Defmodel) > ReportNumThreshold || len(self.Reports.Counter) > ReportNumThreshold ||
		len(self.Reports.Gauge) > ReportNumThreshold || len(self.Reports.Histo) > ReportNumThreshold ||
		len(self.Reports.Retcode) > ReportNumThreshold)
}

func (self *MetricsRepMannager) report(client *http.Client) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic info:%v\n", err)
		}
	}()
	if len(self.Reports.Defmodel) == 0 &&
		len(self.Reports.Counter) == 0 &&
		len(self.Reports.Gauge) == 0 &&
		len(self.Reports.Histo) == 0 &&
		len(self.Reports.Retcode) == 0 {
		return
	}
	data, err := json.Marshal(self.Reports)
	if err != nil {
		fmt.Printf("marshel json err:%v\n", err)
		return
	}

	body := bytes.NewReader(data)
	resp, err := client.Post(metricsapi, "application/x-www-form-urlencoded", body)
	if err != nil {
		fmt.Printf("failed to report,err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	//to reuse
	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		fmt.Printf("ioCopy failed,err:%v\n",err)
	}
	self.cleanupreported()
}

func (self *MetricsRepMannager) cleanupreported() {
	self.Reports.Defmodel = nil
	self.Reports.Counter = nil
	self.Reports.Gauge = nil
	self.Reports.Histo = nil
	self.Reports.Retcode = nil
}

func (self *MetricsRepMannager) defaultreportModel(model *DefModel) error {
	select {
	case self.DefModelChan <- model:
		return nil
	default:
		fmt.Printf("ReportModel chan full,drop model:%v\n", model)
		return ReportBufFullErr
	}
}

func isSuccess(code bool) string {
	if code {
		return "y"
	}

	return "n"
}
