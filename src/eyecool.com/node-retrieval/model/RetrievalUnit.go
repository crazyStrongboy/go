package model

import "time"

type RetrievalUnit struct {
	Id         int  `json:"id" xorm:"pk autoincr"`
	RetrievalId   int
	CameraId  string
	Results string
	Type int
	Total int
	DealNum int
	CreateTime  time.Time `json:"create_time" xorm:"<-"`
}

func (self *RetrievalUnit) TableName() string {
	return "buz_retrieval_unit"
}

