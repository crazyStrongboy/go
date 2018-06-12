package model

import "time"

type Cluster struct {
	Id         int `xorm:"pk autoincr 'id'"`
	Name       string
	Options    string
	CreatorId  int
	UpdateTime time.Time
	Param1     string
	Param2     string
	Param3     string
	Param4     string
}

func (this *Cluster)TableName()string  {
	return "buz_cluster"
}