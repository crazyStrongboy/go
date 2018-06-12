package model

import "time"

type ImageFail struct {
	Id               int64
	ImageRealPath    string
	ImageUri         string
	ImageName        string
	ImageContextPath string
	ImageDesc        string
	FailCode         int
	Name             string
	Gender           int
	Region           int
	Birthday         string
	Nation           string
	Options          string
	PersonId         string
	CreatorId        int
	GroupId          int
	CreateTime       time.Time
	PeopleLibId      int
	ClusterId        int
	RepositoryId     string
	Param1           string
	Param2           string
	ImageUrl         string
}

func (this *ImageFail) TableName() string {
	return "buz_image_fail"
}
