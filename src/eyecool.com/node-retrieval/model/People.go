package model

import "time"

type People struct {
	Id                  int64 `xorm:"pk autoincr 'id'"`
	PubId               string
	PeopleNo            string
	ChannelId           int64
	Name                string
	Gender              int
	Region              int
	Birthday            string
	Nation              int
	Options             string
	PeopleIdType        int
	PersonId            string
	RepositoryPkId      int
	CreatorId           int
	ClusterId           int
	RepositoryId        string
	CustomField         string
	PeopleDetailAddress string
	GroupId             int
	PeopleAddress       string
	PeopleStatus        int
	PeopleType_id       int64
	PeopleComment       string
	PeopleParam1        string
	PeopleParam2        string
	PeopleParam3        string
	CreateTime          time.Time
	UpdateTime          time.Time
	Deleted             int
	ManHeight           float32
	ImageNumber         int
	Param1              string
	Param2              string
	Param3              string
}

func (self *People) TableName() string {
	return "buz_people"
}
