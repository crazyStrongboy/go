package model

import "time"

type Region struct {
	Id         int  `json:"id" xorm:"pk autoincr"`
	ClusterId    int
	Name   string
	PermissionMap string
	ExtraMeta  string
	Status  int
	Mlevel int
	ParentId  int
	CreateTime    int64
	UpdateTime time.Time
}

func (self *Region) TableName() string {
	return "buz_region"
}
