package main

import (
	"io/ioutil"
	"os"
	"time"
	. "deleteFile/conf"
	"strconv"
	"deleteFile/db"
	"strings"
)

var period float64 = 1

func init() {
	InitConf()
	db.InitMysql()
}
func main() {
	clean, err := ConfigFile.GetSection("clean")
	if err != nil {
		SetLogger(err)
	}
	if v, ok := clean["period"]; ok {
		//如果配置了logDir,则选择配置的文件夹作为当前日志的存放点
		period1, err := strconv.ParseFloat(v, 32/64)
		period = period1
		if err != nil {
			SetLogger(err)
		}
	}

	if deleteDir, ok := clean["deleteDir"]; ok {
		dirs := strings.Split(deleteDir, ",")
		for i := 0; i < len(dirs); i++ {
			//存在
			dir := strings.TrimSpace(dirs[i])
			if dir == "" {
				continue
			}
			DeleteDir(dirs[i])
		}
	} else {
		SetLogger("请先配置您想要清除的文件夹")
	}

	deleteDataBase()
}

func deleteDataBase() {
	session := db.MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	sql := "delete from buz_orig_image where create_time < NOW()- INTERVAL ? day"
	_, err := db.MasterDB.Exec(sql, period)
	if err != nil {
		SetLogger("--------------deleteDataBase  :" + err.Error())
	}
	session.Commit()
}

func DeleteDir(dirPth string) {
	SetLogger("要清除的文件夹名是:" + dirPth + "----周期是:" + strconv.FormatFloat(period, 'f', -1, 64))
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		SetLogger(err)
	}
	pthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			DeleteDir(dirPth + pthSep + fi.Name())
		} else {
			lastModified := fi.ModTime()
			sumD := time.Now().Sub(lastModified)
			day := sumD.Hours() / 24
			if day > period {
				err := os.Remove(dirPth + pthSep + fi.Name())
				if err != nil {
					SetLogger(err)
					continue
				}
			}
		}
	}
}
