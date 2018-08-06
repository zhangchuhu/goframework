// @author kordenlu
// @创建时间 2018/02/09 11:45
// 功能描述:

package servant

import (
	"code.yy.com/yytars/goframework/jce/config/taf"
	"errors"
	"code.yy.com/yytars/goframework/kissgo/appzaplog"
	"code.yy.com/yytars/goframework/kissgo/appzaplog/zap"
	"time"
	"strconv"
	"os"
	"io/ioutil"
	"bytes"
)

type TarConfig struct {
	configPrx *taf.Config
	appname string
	servername string
	basepath string
	setdivision string
	maxBakNum int
}

var defaultTarConfig *TarConfig

// 业务用来获取配置的基本路径
func GetConfBasePath()(string,error)  {
	var basepath string
	if defaultTarConfig == nil{
		return basepath,errors.New("nil tar config")
	}
	basepath = defaultTarConfig.basepath
	return basepath,nil
}

func initTarConfig(comm ICommunicator,conf *serverConfig,maxBakNum int)error  {
	if comm == nil || conf == nil{
		return errors.New("nil param")
	}
	defaultTarConfig = &TarConfig{
		appname:conf.App,
		servername:conf.Server,
		basepath:conf.BasePath,
		setdivision:"",
		maxBakNum:maxBakNum,
	}
	defaultTarConfig.configPrx = &taf.Config{}
	defaultTarConfig.configPrx.SetServant(comm.GetServantProxy(conf.config))
	return nil
}

// addConfig
// 1 take remote config to local
// 2 backup local config
func addConfig(filename string,bAppConfigOnly bool)(string,error)  {
	if defaultTarConfig == nil{
		return "tarconfig client not initialized",errors.New("tarconfig client not initialized")
	}
	return defaultTarConfig.addConfig(filename,bAppConfigOnly)
}

func pathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadConf(filename string) ([]byte,error) {
	basepath,err := GetConfBasePath()
	if err != nil {
		return nil,err
	}
	file_ := basepath+filename
	//if pathExist(file_){
	//	return ioutil.ReadFile(file_)
	//}
	if _,err := addConfig(filename,false);err != nil{
		return nil,err
	}
	return ioutil.ReadFile(file_)
}

func samecontent(newpath, oldpath string)(bool,error)  {
	newfilebyte,err := ioutil.ReadFile(newpath)
	if err != nil {
		return false,err
	}

	// check if old exist, not exist should return false directly
	if !pathExist(oldpath){
		return false,nil
	}

	oldfilebyte,err := ioutil.ReadFile(oldpath)
	if err != nil {
		return false,err
	}

	if bytes.Compare(newfilebyte, oldfilebyte) == 0{
		return true,nil
	}
	return false,nil
}

func index2file(src string, index int64) string {
	return src +"."+strconv.FormatInt(index,10)+".bak"
}

func (tc *TarConfig)addConfig(filename string,bAppConfigOnly bool)(string,error)  {
	var(
		sFullFileName string = tc.basepath + "/" + filename
	)
	newfile,err := tc.getRemoteFile(filename,bAppConfigOnly)
	if err != nil {
		// try to use local file instead
		if _,err := os.Stat(sFullFileName);err == nil{
			appzaplog.Warn("getRemoteFile failed,use the local config",zap.String("filename",filename),zap.Error(err))
			return "[fail] get remote config:" + filename + "fail,use the local config.",nil
		}
		// no local find, return err
		appzaplog.Error("getRemoteFile and local failed",zap.String("filename",filename),zap.Error(err))
		return err.Error(),err
	}


	same,err := samecontent(newfile,sFullFileName)

	if err != nil{
		appzaplog.Error("samecontent failed",zap.Error(err))
		return err.Error(),err
	}

	if !same {
		// move one by one
		for i := int64(tc.maxBakNum-1);i > 1;i--{
			if _,err := os.Stat(index2file(sFullFileName,i));err == nil{
				os.Rename(index2file(sFullFileName,i),index2file(sFullFileName,i+1))
			}else {
				appzaplog.Warn("stat file failed",zap.Error(err),zap.String("filename",index2file(sFullFileName,int64(i))))
			}
		}
		// move old to .1.bak
		if pathExist(sFullFileName){
			os.Rename(sFullFileName,index2file(sFullFileName,1))
		}
	}

	os.Rename(newfile,sFullFileName)

	if _,err := os.Stat(sFullFileName);err != nil{
		appzaplog.Error("Stat file failed",zap.String("sFullFileName",sFullFileName),zap.Error(err))
		return err.Error(),err
	}
	return "[succ] get remote config:"+filename,nil
}

func (tc *TarConfig)getRemoteFile(filename string,bAppConfigOnly bool)(string,error)  {
	if tc.configPrx == nil{
		appzaplog.Error("tarconfig proxy not ready",zap.String("filename",filename))
		return "",errors.New("tarconfig proxy not ready")
	}
	var (
		servername string
		bconfig string
		ret int32
		err error
	)
	if !bAppConfigOnly{
		servername = tc.servername
	}

	if tc.setdivision == ""{
		ret,err = tc.configPrx.LoadConfig(tc.appname,servername,filename,&bconfig)
	}else {
		//todo not support set yet, just impl it first
		cinfo := taf.ConfigInfo{
			Appname:tc.appname,
			Servername:servername,
			Filename:filename,
			BAppOnly:bAppConfigOnly,
			Setdivision:tc.setdivision,
		}
		ret,err = tc.configPrx.LoadConfigByInfo(cinfo,&bconfig)
	}
    appzaplog.Info("[+]getRemoteFile LoadConfig",zap.Int32("ret",ret),zap.Error(err))
	if err != nil {
		return "",err
	}
	if ret !=0{
		appzaplog.Warn("LoadConfig ret failed",zap.Int32("ret",ret))
		return "",errors.New("request failed")
	}

	newFile := tc.basepath + "/" + filename + "." + strconv.FormatInt(time.Now().Unix(),10)
	cfile,err := os.Create(newFile)
	if err != nil {
		appzaplog.Error("create newfile failed",zap.String("newFile",newFile),zap.Error(err))
		return "",err
	}
	defer cfile.Close()
	if _,err := cfile.WriteString(bconfig);err != nil{
		appzaplog.Error("WriteString failed",zap.String("newFile",newFile),zap.Error(err))
		return "",err
	}
	return newFile,nil
}