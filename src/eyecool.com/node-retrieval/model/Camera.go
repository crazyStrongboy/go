package model

import "time"

type Camera struct {
	PkId         int  `json:"pk_id" xorm:"pk autoincr"`
	Id   string
	Name   string
	Url  string
	Ip string
	PredecessorId string
	RecParams string
	PermissionMap string
	ExtraMeta string
	Status int
	RegionId int
	CreatorId int
	ClusterId int
	UpdateTime  time.Time `json:"created_at" `
	CreateTime  time.Time `json:"created_at" `
}

type CameraResponse struct{
	Id string
	Name string
	Url string
	Enabled int
	RecParams string
	PermisssionMap string
	predecessorIds []string
}

func (self *Camera) TableName() string {
	return "buz_camera"
}

