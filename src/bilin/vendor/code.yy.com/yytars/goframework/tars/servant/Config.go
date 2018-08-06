package servant

import (
	"fmt"
	"time"
	"code.yy.com/yytars/goframework/tars/tarsserver"
	"code.yy.com/yytars/goframework/tars/util/conf"
	"code.yy.com/yytars/goframework/tars/util/endpoint"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"strconv"
	"strings"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"errors"
)

const (
	pathspliter = "/"

	containerpath = "/tars/application<container>"

	serverpath = "/tars/application/server"
	clientpath = "/tars/application/client"

	// sever config keys,path: /tars/application/server/
	serKeyNode                = "node"
	serKeyApp                 = "app"
	serKeyServer              = "server"
	serKeyLocalip             = "localip"
	serKeyLocal               = "local"
	serKeyBasepath            = "basepath"
	serKeyDatapath            = "datapath"
	serKeyLogpath             = "logpath"
	serKeyLogsize             = "logsize"
	serKeyConfig              = "config"
	serKeyNotify              = "notify"
	serKeyLog                 = "log"
	serMetricsEnable          = "metricsenable"
	serKeyDeactivatingtimeout = "deactivating-timeout"
	serKeyLogLevel            = "logLevel"
	serKeyNetThread           = "netthread" // not config yet, default to 1?

	// Adapter key /tars/application/server/*/
	adapterKeyAllow        = "allow"
	adapterKeyEndPoint     = "endpoint"
	adapterKeyHandleGroup  = "handlegroup"
	adapterKeyMaxConns     = "maxconns"
	adapterKeyProtocol     = "protocol"
	adapterKeyQueueCap     = "queuecap"
	adapterKeyQueueTimeout = "queuetimeout"
	adapterKeyServant      = "servant"
	adapterKeyThreads      = "threads"

	// client config keys path: /tars/application/client/
	clientKeyLocator                 = "locator"
	clientKeySyncInvokeTimeout       = "sync-invoke-timeout"
	clientKeyAsyncInvokeTimeout      = "async-invoke-timeout"
	clientKeyRefreshEndpointInterval = "refresh-endpoint-interval"
	clientKeyStat                    = "stat"
	clientKeyProperty                = "property"
	clientKeyReportInterval          = "report-interval"
	clientKeySampleRate              = "sample-rate"
	clientKeyMaxSampleCount          = "max-sample-count"
	clientKeyAsyncThread             = "asyncthread"
	clientKeyNetConnectionNum               = "netconnectionnum"
	clientKeyModuleName              = "modulename"

	//
	SDK_VERSION = "1.0.0"
	PBSERVANT = "pbservant"

	// fulladminobjname和adminadaptername不能更改? node节点已经写死
	fulladminobjname = "AdminObj"
	adminadaptername = "AdminAdapter"
)

var (
	svrCfg *serverConfig
	cltCfg *clientConfig
)

// LoadConfig 加载配置信息
func loadConfig(c *conf.Conf) error {
	svrCfg = &serverConfig{
		Adapters: make(map[string]adapterConfig),
	}
	cltCfg = &clientConfig{}

	if err := svrCfg.loadserverConfig(c); err != nil {
		appzaplog.Error("loadserverConfig err", zap.Error(err))
		return err
	}

	// start  server log
	if err := appzaplog.SetLogLevel(strings.ToLower(svrCfg.LogLevel)); err != nil {
		appzaplog.Error("SetLogLevel failed", zap.Error(err))
		return err
	}

	if err := cltCfg.loadclientConfig(c); err != nil {
		appzaplog.Error("loadclientConfig err", zap.Error(err))
		return err
	}

	return nil
}

func GetServerConfig() *serverConfig {
	return svrCfg
}

func GetClientConfig() *clientConfig {
	return cltCfg
}

type adapterConfig struct {
	Endpoint endpoint.Endpoint
	Protocol string
	Obj      string
	Threads  int
	MaxConns int // 最大连接数目
	QueueCap int // 队列长度
}

type serverConfig struct {
	Node      string `json:"node"`
	App       string `json:"app"`
	Server    string `json:"server"`
	LocalIP   string `json:"localip"`
	Local     string
	BasePath  string
	DataPath  string
	Container string // not used yet, place it here ?
	LogPath   string
	LogSize   string
	LogLevel  string
	Version   string
	config    string
	notify    string
	log       string
	metricsenable string
	Adapters map[string]adapterConfig
}

func (this *serverConfig) loadserverConfig(c *conf.Conf) error {
	sMap := c.GetMap(serverpath)
	this.Node = sMap[serKeyNode]
	this.App = sMap[serKeyApp]
	this.Server = sMap[serKeyServer]
	this.LocalIP = sMap[serKeyLocalip]
	this.Local = sMap[serKeyLocal]
	this.Container = c.GetString(containerpath)
	//init log
	this.LogPath = sMap[serKeyLogpath]
	this.LogSize = sMap[serKeyLogsize]
	this.LogLevel = sMap[serKeyLogLevel]
	this.config = sMap[serKeyConfig]
	this.notify = sMap[serKeyNotify]
	this.BasePath = sMap[serKeyBasepath]
	this.DataPath = sMap[serKeyDatapath]

	this.log = sMap[serKeyLog]
	if metrics,ok := sMap[serMetricsEnable];ok{
		this.metricsenable = metrics
	}else {
		this.metricsenable = "yes"
	}
	//add version info
	this.Version = SDK_VERSION

	return this.loadServerAdapter(c)
}

func (this *serverConfig) loadServerAdapter(c *conf.Conf) error {
	serList := c.GetDomain(serverpath)

	for _, adapter := range serList {
		adapterpath := serverpath + pathspliter + adapter
		appzaplog.Debug("load conf adapterpath", zap.String("adapterpath", adapterpath))

		endString := c.GetString(adapterpath + genKeyPath(adapterKeyEndPoint))
		end := endpoint.Parse(endString)

		svrObj := c.GetString(adapterpath + genKeyPath(adapterKeyServant))

		this.Adapters[adapter] = adapterConfig{
			Endpoint: end,
			Protocol: c.GetString(adapterpath + genKeyPath(adapterKeyProtocol)),
			Obj:      svrObj,
			Threads:  c.GetInt(adapterpath + genKeyPath(adapterKeyThreads)),
			MaxConns: c.GetInt(adapterpath + genKeyPath(adapterKeyMaxConns)),
			QueueCap: c.GetInt(adapterpath + genKeyPath(adapterKeyQueueCap)),
		}

		host := end.Host
		if end.Bind != "" {
			host = end.Bind
		}
		servantConfig[svrObj] = &tarsserver.TarsServerConf{
			Proto:         end.Proto,
			Address:       fmt.Sprintf("%s:%d", host, end.Port),
			MaxAccept:     this.Adapters[adapter].MaxConns,
			MaxInvoke:     this.Adapters[adapter].QueueCap,
			AcceptTimeout: time.Millisecond * 500,
			ReadTimeout:   time.Millisecond * 100,
			WriteTimeout:  time.Millisecond * 100,
			HandleTimeout: time.Millisecond * 60000,
			IdleTimeout:   time.Millisecond * 600000,
		}
	}
	return this.addAdminAdapter()
}

func (this *serverConfig) addAdminAdapter() error {
	appzaplog.Debug("add admin adapter config", zap.Any("servantConfig", servantConfig))
	localpoint := endpoint.Parse(this.Local)

	// add admin adpater
	servantConfig[fulladminobjname] = &tarsserver.TarsServerConf{
		Proto:         "tcp",
		Address:       fmt.Sprintf("%s:%d", localpoint.Host, localpoint.Port),
		MaxAccept:     1000,
		MaxInvoke:     5000,
		AcceptTimeout: time.Millisecond * 500,
		ReadTimeout:   time.Millisecond * 100,
		WriteTimeout:  time.Millisecond * 100,
		HandleTimeout: time.Millisecond * 60000,
		IdleTimeout:   time.Millisecond * 600000,
	}

	this.Adapters[adminadaptername] = adapterConfig{
		Endpoint: localpoint,
		Protocol: "tcp",
		Obj:      fulladminobjname,
		Threads:  1,
		MaxConns: 1000,
		QueueCap: 5000,
	}
	return nil
}

type clientConfig struct {
	Locator                 string
	stat                    string
	property                string
	refreshEndpointInterval int
	reportInterval          int
	asyncthread             int
	netconnectionnum        int
	modulename              string
	syncInvokeTimeout       int
}

func (this *clientConfig) loadclientConfig(c *conf.Conf) error {
	cMap := c.GetMap(clientpath)
	this.Locator = cMap[clientKeyLocator]
	this.stat = cMap[clientKeyStat]

	refreshendpointinterval := cMap[clientKeyRefreshEndpointInterval]
	temp, err := strconv.Atoi(refreshendpointinterval)
	if err != nil {
		appzaplog.Error("ParseInt refreshEndpointInterval failed", zap.Error(err), zap.String("refreshendpointinterval", refreshendpointinterval))
		return err
	}
	this.refreshEndpointInterval = temp

	asyncthread := cMap[clientKeyAsyncThread]
	temp, err = strconv.Atoi(asyncthread)
	if err != nil {
		appzaplog.Error("ParseInt asyncthread failed", zap.Error(err), zap.String("asyncthread", asyncthread))
		return err
	}
	this.asyncthread = temp

	netconnectionnum := cMap[clientKeyNetConnectionNum]
	temp, err = strconv.Atoi(netconnectionnum)
	if err != nil {
		appzaplog.Warn("ParseInt netconnectionnum failed", zap.Error(err), zap.String("netconnectionnum", netconnectionnum))
		//return err
	}else {
		this.netconnectionnum = temp
	}

	syncinvoketimeout := cMap[clientKeySyncInvokeTimeout]
	temp, err = strconv.Atoi(syncinvoketimeout)
	if err != nil {
		appzaplog.Error("ParseInt syncinvoketimeout failed", zap.Error(err), zap.String("syncinvoketimeout", syncinvoketimeout))
		return err
	}
	this.syncInvokeTimeout = temp

	return nil
}

func genKeyPath(key string) string {
	return "<" + key + ">"
}

var (
	NilServerConfig = errors.New("nil server config")
	EmptyAppOrServerName = errors.New("empty app or server name")
)
func fullObjName(objname string) (string, error) {
	var fullobjname string
	pos := strings.Index(objname, ".")
	if pos > 0 {
		// if objname with ., take it as fullname
		fullobjname =  objname
	} else {
		switch {
		case GetServerConfig() == nil:
			return fullobjname, NilServerConfig
		case GetServerConfig().App == "" || GetServerConfig().Server == "":
			return fullobjname,EmptyAppOrServerName
		}
		fullobjname = strings.Join([]string{
			GetServerConfig().App,
			GetServerConfig().Server,
			objname,
		}, ".")
	}
	return fullobjname,nil
}

func takeclientsynctimeinvoketimeout()time.Duration  {
	clientconf := GetClientConfig()
	if clientconf == nil || clientconf.syncInvokeTimeout == 0{
		return 3*time.Second
	}
	return time.Millisecond*time.Duration(clientconf.syncInvokeTimeout)
}