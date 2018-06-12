package model

import "time"

type Retrieval struct {
	Id               int64     `json:"id" xorm:"pk autoincr"`
	RetrievalQueryId string
	RepositoryId     string
	AsyncQuery       bool
	CreatorId        int
	Status           int       `json:"status" xorm:"-"`
	Results          string
	ClusterId        int
	UpdateTime       time.Time `json:"created_at" xorm:"<-"`
	CreateTime       time.Time `json:"created_at" xorm:"<-"`

	ExtraFields   string
	PeopleId      int64
	RepositoryIds string
	CameraIds     string
	VideoIds      string
	CameraId      string
	VideoId       string
	PersonId      string
	Name          string
	Threshold     float64
	UsingAnn      bool
	Topk          int32
	ConditionJson string
	OrderJson     string
	StartIndex    int
	LimitResult   int
	Timestamp     string
	Total         int
}

func (self *Retrieval) TableName() string {
	return "buz_retrieval"
}
