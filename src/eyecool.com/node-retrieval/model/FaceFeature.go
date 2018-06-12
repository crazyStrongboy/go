package model

import "time"

type FaceFeature struct {
	PkId         int  `json:"pk_id" xorm:"pk autoincr"`
	FaceImageId   string
	RepositoryId   string
	RepositoryPkId int
	PeopleId  int64
	Feat string
	Status int
	UpdateTime  time.Time `json:"created_at" xorm:"<-"`
	CreateTime  time.Time `json:"created_at" xorm:"<-"`

	W int
	X int
	Y int
	H int
	ImageId int64
}

func (self *FaceFeature) TableName() string {
	return "buz_face_feature"
}

