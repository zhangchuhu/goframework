package servant

import (
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"code.yy.com/yytars/goframework/tars/servant/protocol"
	"code.yy.com/yytars/goframework/tars/util/conf"
	"errors"
	"strings"
	"sync"
)

var (
	emptyError   = errors.New("empty command")
	noParamError = errors.New("not enoguht param")
)

var (
	pushconfigchan = make(chan *pbtaf.ConfigPushNotice, 10)
)

const (
	TARS_CMD_LOAD_CONFIG         = "tars.loadconfig"    //从配置中心, 拉取配置下来: tars.loadconfig filename
	TARS_CMD_SET_LOG_LEVEL       = "tars.setloglevel"   //设置滚动日志的等级: tars.setloglevel [NONE, ERROR, WARN, DEBUG]
	TARS_CMD_VIEW_STATUS         = "tars.viewstatus"    //查看服务状态
	TARS_CMD_VIEW_VERSION        = "tars.viewversion"   //查看服务采用TARS的版本
	TARS_CMD_CONNECTIONS         = "tars.connection"    //查看当前链接情况
	TARS_CMD_LOAD_PROPERTY       = "tars.loadproperty"  //使配置文件的property信息生效
	TARS_CMD_VIEW_ADMIN_COMMANDS = "tars.help"          //帮助查看服务支持的管理命令
	TARS_CMD_SET_DYEING          = "tars.setdyeing"     //设置染色信息: tars.setdyeing key servant [interface]
	TARS_CMD_CLOSE_COUT          = "tars.closecout"     //设置是否启关闭cout\cin\cerr: tars.openthreadcontext yes/NO 服务重启才生效
	TARS_CMD_SET_DAYLOG_LEVEL    = "tars.enabledaylog"  //设置按天日志是否输出: tars.enabledaylog [remote|local]|[logname]|[true|false]
	TARS_CMD_CLOSE_CORE          = "tars.closecore"     //设置服务的core limit:  tars.setlimit [yes|no]
	TARS_CMD_RELOAD_LOCATOR      = "tars.reloadlocator" //重新加载locator的配置信息
	TARS_CMD_SET_CALL_STATUS     = "tars.setcallstatus" //设置rpc调用的状态
	TARS_CMD_GET_CALL_STATUS     = "tars.getcallstatus" //获取设置的rpc调用状态
)

type Admin struct {
}

func (a *Admin) Shutdown() error {
	for obj, s := range tarsAndPbSvrs {
		appzaplog.Debug("shutdown", zap.String("obj", obj))
		//TODO
		go s.Shutdown()
	}
	shutdown <- true
	return nil
}

func (a *Admin) Notify(command string) (string, error) {
	appzaplog.Debug("recv Notify command", zap.String("command", command))

	cm, params, err := parseCommand(command)
	if err != nil {
		return "parseCommand failed", err
	}

	switch cm {
	case TARS_CMD_VIEW_VERSION:
		return GetServerConfig().Version, nil
	case TARS_CMD_SET_LOG_LEVEL:
		return setLogLevelCMD(params)
	case TARS_CMD_CONNECTIONS:
		return connctionCMD()
	case TARS_CMD_LOAD_CONFIG:
		return loadconfigCMD(params)
	case TARS_CMD_VIEW_STATUS:
		return viewstatusCMD()
	case TARS_CMD_RELOAD_LOCATOR:
		return reloadLocatorCMD(params)
	case TARS_CMD_SET_CALL_STATUS:
		return setCallStatus(params)
	case TARS_CMD_GET_CALL_STATUS:
		return getCallStatus()
	default:
		return "what are you say?", nil
	}
}

func parseCommand(command string) (cm string, params []string, err error) {
	commandfields := strings.Fields(command)
	switch len(commandfields) {
	case 0:
		err = emptyError
	case 1:
		cm = commandfields[0]
	default:
		cm = commandfields[0]
		params = commandfields[1:]
	}
	return
}

//todo
func connctionCMD() (string, error) {
	return "not support yet", nil
}

var (
	lock            sync.Mutex
	disabledcallmap map[string]struct{} = make(map[string]struct{})
)

func getCallStatus() (string, error) {
	lock.Lock()
	disablemap := disabledcallmap
	lock.Unlock()
	var disablevector []string
	for k, _ := range disablemap {
		disablevector = append(disablevector, k)
	}
	return strings.Join(disablevector, ":"), nil
}

func setCallStatus(params []string) (string, error) {
	if len(params) == 2 && params[0] != "" {
		switch params[1] {
		case "1":
			lock.Lock()
			delete(disabledcallmap, params[0])
			lock.Unlock()
		case "0":
			lock.Lock()
			disabledcallmap[params[0]] = struct{}{}
			lock.Unlock()
		default:
			return "illegal params", noParamError
		}
		return "set call status " + params[1] + " success", nil
	} else {
		return "illegal params", noParamError
	}
}

func callDisabled(objfunname string)bool  {
	lock.Lock()
	if _,exist := disabledcallmap[objfunname];exist{
		lock.Unlock()
		return true
	}
	lock.Unlock()
	return false
}

func pushconfignotice(notice *pbtaf.ConfigPushNotice) {
	select {
	case pushconfigchan <- notice:
	default:
		appzaplog.Warn("no config consumer or too slow consumer,notice droped", zap.Any("notice", notice))
	}
}

func SubConfigPush() <-chan *pbtaf.ConfigPushNotice {
	return pushconfigchan
}

func loadconfigCMD(params []string) (string, error) {
	var (
		notifystr string = "load config success"
		err       error
		notice    = &pbtaf.ConfigPushNotice{
			Retcode: pbtaf.ConfigPushNotice_SUCCESS,
		}
	)
	defer func() {
		reportNotifyInfo(notifystr)
	}()

	if len(params) == 0 {
		notifystr = "no file name in cmd"
		return notifystr, errors.New(notifystr)
	}
	filename := params[0]
	notice.Filename = filename
	// todo ,support bAppConfigOnly
	notifystr, err = addConfig(filename, false)
	if err != nil {
		// fake it not err
		notifystr = err.Error()
		notice.Retcode = pbtaf.ConfigPushNotice_FAILED
	}
	pushconfignotice(notice)
	return notifystr, nil
}

func viewstatusCMD() (string, error) {
	return "not support yet", nil
}

// 当 locator失效后,可以通过命令来手工切换到其他的locator
// 原来的locator现在还没有释放机制
func reloadLocatorCMD(params []string) (string, error) {
	if len(params) >= 1 && params[0] == "reload" {
		c := conf.NewConf(configFile)
		locator := c.GetString("/tars/application/client<locator>")
		if locator == "" {
			return "locator empty", nil
		}
		// todo remove not used locator
		startFrameWorkComm().SetProperty(clientKeyLocator, locator)
		return "load locator:" + locator, nil
	} else {
		return "please input right paras", noParamError
	}
}

func setLogLevelCMD(params []string) (string, error) {
	if len(params) >= 1 {
		if err := appzaplog.SetLogLevel(strings.ToLower(params[0])); err != nil {
			return "SetLogLevel failed", err
		}
		return "set level to " + params[0], nil
	} else {
		return "not enoguht param", noParamError
	}
}
