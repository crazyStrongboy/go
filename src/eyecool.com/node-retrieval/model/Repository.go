package model

import "time"

type Repository struct {
	PkId             int       `json:"pk_id" xorm:"pk autoincr"`
	Id               string
	Name             string
	RepositoryId     string    `xorm:"-"`
	RepositoryType   int
	TotalPictureNum  int
	FaceImageNum     int
	FailedPictureNum int
	CreatorId        int
	PermissionMap    string
	ExtraMeta        string    `xorm:"extra_meta"`
	Status           int
	Options          string
	ClusterId        int
	UpdateTime       time.Time `json:"created_at" `
	CreateTime       int64     `json:"created_at" `
}

func (self *Repository) TableName() string {
	return "buz_repository"
}
