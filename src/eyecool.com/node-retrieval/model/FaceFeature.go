package model

import "time"

type FaceFeature struct {
	PkId           int       `json:"pk_id" xorm:"pk autoincr"`
	FaceImageId    string
	RepositoryId   string
	RepositoryPkId int
	PeopleId       int64
	Feat           string
	Status         int
	UpdateTime     time.Time `json:"created_at" xorm:"<-"`
	CreateTime     time.Time `json:"created_at" xorm:"<-"`
	FaceRect       string
	FaceProp       string
	W              int
	X              int
	Y              int
	H              int
	ImageId        int64
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	T int `json:"t"`
	B int `json:"b"`
}

type Prop struct {
	Age         int `json:"age"`
	Gender      int `json:"gender"`
	Race        int `json:"race"`
	SmileLevel  int `json:"smileLevel"`
	BeautyLevel int `json:"beautyLevel"`
}

func (self *FaceFeature) TableName() string {
	return "buz_face_feature"
}
