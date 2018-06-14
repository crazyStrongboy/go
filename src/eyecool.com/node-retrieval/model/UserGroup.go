package model

import "time"

type UserGroup struct {
	Id             int       `json:"id" xorm:"pk autoincr 'id'"`
	Name           string    `json:"name"`
	ExtraMeta      string    `json:"extra_meta"`
	Status         int       `json:"omitempty"`
	ClusterId      int       `json:"omitempty"`
	RepositoryId   string    `json:"omitempty"`
	CreateTime     int64     `json:"create_time"`
	GroupLevel     int       `json:"omitempty"`
	ParentId       int       `json:"omitempty"`
	UpdateTime     time.Time `json:"omitempty"`
	Param1         string    `json:"omitempty"`
	Param2         string    `json:"omitempty"`
	Param3         string    `json:"omitempty"`
	Param4         string    `json:"omitempty"`
	PredecessorIds []int     `json:"predecessor_ids" xorm:"-"`
	Predecessor_id string    `json:"-" xorm:"-"`
	Extra_meta     string    `json:"-" xorm:"-"`
	Create_time    int64     `json:"-" xorm:"-"`
}

func (self *UserGroup) TableName() string {
	return "buz_user_group"
}
