package model

import "time"

type AlarmInfo struct {
	Id                       int       `json:"id" xorm:"pk autoincr"`
	TaskId                   string
	CameraId                 string
	Surveillances            string
	RepositoryId             string
	AlarmPeopleId            int64
	AlarmImageContextPath    string
	AlarmCropImageProperties string
	AlarmCropImageUri        string
	AlarmOrigImageId         int64
	AlarmOrigImageUuid       string
	AlarmOrigImageRectIdx    string
	AlarmOrigImageUri        string
	AlarmScoreOthers         string
	AlarmScore               float32
	AlarmTmplFeatrueId1      string
	AlarmTmplFeatrueId2      string
	AlarmTmplFeatrueId3      string
	AlarmTmplScore1          float32
	AlarmTmplScore2          float32
	AlarmTmplScore3          float32
	Timestamp                int64
	ClusterId                int
	CreateTime               time.Time `json:"create_time" xorm:"<-"`
}

func (self *AlarmInfo) TableName() string {
	return "buz_alarm_info"
}
