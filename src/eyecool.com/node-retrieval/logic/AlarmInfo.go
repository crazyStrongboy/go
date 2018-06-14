package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"time"
	"log"
)

type HitCondition struct {
	Hit_similarity HitSimilarity
	Timestamp      Timestamp
	StartTime      time.Time
	EndTime        time.Time
	Start          int
	Limit          int
}

type HitSimilarity struct {
	Gte int //报警的阈值
}
type Timestamp struct {
	Gte string //start time  格式是:yyyy-MM-dd  HH:mm:ss---->2018-06-13 14:50:33
	Lte string //end time  格式是:yyyy-MM-dd HH:mm:ss ---->2018-06-13 14:50:33
}
type AlarmInfoLogic struct{}

var DefaultAlarmInfo = AlarmInfoLogic{}

func (AlarmInfoLogic) Insert(alarmInfo *model.AlarmInfo) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	_, err := MasterDB.Insert(alarmInfo)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("Insert alarm info error :", err)
		return err
	}
	session.Commit()
	return nil
}
func (infoLogic *AlarmInfoLogic) FindAlarmInfosByTaskIdsAndCondition(taskIds []string, condition *HitCondition) []*model.AlarmInfo {
	alarmInfos := make([]*model.AlarmInfo, 0)
	session := MasterDB.Table(model.AlarmInfo{})
	session.In("task_id", taskIds)
	session.Where("alarm_score > ?", condition.Hit_similarity.Gte)
	session.Where("? < create_time < ?", condition.StartTime, condition.EndTime)
	session.Limit(condition.Limit, condition.Start)
	err := session.Find(&alarmInfos)
	if err != nil {
		log.Println("FindAlarmInfosByTaskIdsAndCondition err :", err)
	}
	return alarmInfos
}
func (infoLogic *AlarmInfoLogic) FindAlarmInfosByCameraIdssAndCondition(cameraIds []string, condition *HitCondition) []*model.AlarmInfo {
	alarmInfos := make([]*model.AlarmInfo, 0)
	session := MasterDB.Table(new(model.AlarmInfo).TableName())
	session.In("camera_id", cameraIds)
	session.Where("alarm_score > ?", condition.Hit_similarity.Gte)
	session.Where("? < create_time < ?", condition.StartTime, condition.EndTime)
	session.Limit(condition.Limit, condition.Start)
	err := session.Find(&alarmInfos)
	if err != nil {
		log.Println("FindAlarmInfosByCameraIdssAndCondition err :", err)
	}
	return alarmInfos
}
func (infoLogic *AlarmInfoLogic) FindAlarmInfosBySurveillancessAndCondition(surveillances []Surveillance, condition *HitCondition) []*model.AlarmInfo {
	alarmInfos := make([]*model.AlarmInfo, 0)
	session := MasterDB.Table(new(model.AlarmInfo).TableName())
	for _, surveillance := range surveillances {
		session.Or("alarm_score > ? and repository_id = ? and camera_id = ? and alarm_score > ? and ? < create_time < ?", condition.Hit_similarity.Gte, surveillance.RepositoryId, surveillance.CameraId, surveillance.Threshold, condition.StartTime, condition.EndTime)
	}
	session.Limit(condition.Limit, condition.Start)
	err := session.Find(&alarmInfos)
	if err != nil {
		log.Println("FindAlarmInfosBySurveillancessAndCondition err :", err)
	}
	return alarmInfos
}
