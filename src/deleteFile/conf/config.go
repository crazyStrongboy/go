package conf

import (
	"github.com/Unknwon/goconfig"
	"os"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

const mainIniPath = "/config/env.ini"
const LogName = "deletePic.log"

var (
	ConfigFile *goconfig.ConfigFile
	logDir = "D:\\opt\\eyecool\\modules"
	ROOT string
)
func InitConf(){
	curFilename := os.Args[0]
	binaryPath, err := exec.LookPath(curFilename)
	if err != nil {
		panic(err)
	}
	binaryPath, err = filepath.Abs(binaryPath)
	if err != nil {
		panic(err)
	}
	ROOT = filepath.Dir(binaryPath)//程序运行的当前位置

	configPath := ROOT + mainIniPath
	ConfigFile, err = goconfig.LoadConfigFile(configPath)
	clean, err := ConfigFile.GetSection("clean")
	if v, ok := clean["logDir"]; ok {
		//如果配置了logDir,则选择配置的文件夹作为当前日志的存放点
		logDir = v
	}
	if err != nil{
		SetLogger(err)
	}
}

//设置日志的方法
func SetLogger(myerr interface{}) {
	os.MkdirAll(logDir, os.ModePerm)
	fileName := logDir + "\\" + LogName
	if _, err := os.Stat(fileName); err != nil {
		//文件不存在
		os.Create(fileName)
	} else {
		//文件存在
	}
	logfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(999)
	}
	logger := log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	defer logfile.Close()
	logger.Println(myerr)
}
