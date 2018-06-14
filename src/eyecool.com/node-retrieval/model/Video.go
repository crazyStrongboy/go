package model

import "time"

type Video struct {
	PkId          int    `json:"pk_id" xorm:"pk autoincr"`
	Id            string
	Name          string
	Url           string
	RepositoryId  string
	PermissionMap string
	ExtraMeta     string `xorm:"extra_meta"`
	Enabled       int
	CreatorId     int64
	ClusterId     int
	CameraId      int64
	CreateTime    int64
	RecParams     string
	UpdateTime    time.Time
	Param1        string
	Param2        string
	Param3        string
	Param4        string
}

func (self *Video) TableName() string {
	return "buz_video"
}
