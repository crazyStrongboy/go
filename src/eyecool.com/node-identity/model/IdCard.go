package model

import "time"

type IdCard struct {
	Id               int64  `json:"id" xorm:"pk autoincr"`
	PubId            string
	IdNum            string
	Name             string
	Sex              int //性别 0 男 1女
	Nation           string
	Address          string
	IdcOpenUnit      string
	Finger           string `json:"finger" xorm:"text"`
	BirthDate        string `json:"birthdate" xorm:"string 'birthdate'"`
	Age              int
	Status           int
	ExpireDate       string
	EffectedDate     string
	Issue            string
	ImageContextPath string
	ImageUri         string
	RecogTime        time.Time //识别时间
	Path             string    //身份证图片地址
	CreateTime       time.Time
	UpdateTime       time.Time
	Deleted          int //是否已删除,0:正常 2:已删除
}

func (self *IdCard) TableName() string {
	return "buz_idcard"
}
