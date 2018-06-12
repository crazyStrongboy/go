package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
)

type AlarmInfoLogic struct{}

var DefaultAlarmInfo =AlarmInfoLogic{}

func (AlarmInfoLogic) Insert(alarmInfo *model.AlarmInfo) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	_, err := MasterDB.Insert(alarmInfo)
	if err != nil {
		session.Rollback()
		logger.Errorln("insert alarmInfo error:", err)
		return  err
	}
	session.Commit()
	return nil
}