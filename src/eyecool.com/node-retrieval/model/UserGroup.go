package model

import "time"

type UserGroup struct {
	Id             int    `xorm:"pk autoincr 'id'"`
	Name           string
	ExtraMeta      string
	Status         int
	ClusterId      int
	RepositoryId   string
	CreateTime     int64
	GroupLevel     int
	ParentId       int
	UpdateTime     time.Time
	Param1         string
	Param2         string
	Param3         string
	Param4         string
	PredecessorIds []int  `xorm:"-"`
	Predecessor_id string `xorm:"-"`
	Extra_meta     string `xorm:"-"`
	Create_time    int64  `xorm:"-"`
}

func (self *UserGroup) TableName() string {
	return "buz_user_group"
}
