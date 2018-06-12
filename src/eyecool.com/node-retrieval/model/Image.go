package model

import "time"

type Image struct {
	Id               int64 `xorm:"pk autoincr 'id'"`
	RepositoryId     int
	ClusterId        int
	PubId            string
	PeopleId         int64
	ImageType        int32
	ImageRealPath    string
	ImageUri         string
	ImageUrl         string
	ImageName        string
	Status           int
	ImageContextPath string
	Param1           string
	Param2           string
	CreateTime       time.Time
	UpdateTime       time.Time
	Deleted          int32
}

func (self *Image) TableName() string {
	return "buz_image"
}
