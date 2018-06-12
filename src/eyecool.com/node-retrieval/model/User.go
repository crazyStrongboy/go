package model

import "time"

type User struct {
	Id            int    `xorm:"pk autoincr 'id'"`
	Name          string
	Password      string
	ExtraMeta     string
	Status        int
	GroupId       int
	ClusterId     int
	RepositoryId  int
	CreateTime    int64
	UserLevel     int
	ParentId      int
	PermissionMap string
	UpdateTime    time.Time
	Param1        string `json:"omitempty"`
	Param2        string `json:"omitempty"`
	Param3        string `json:"omitempty"`
	Param4        string `json:"omitempty"`

	//in
	Predecessor_id string `xorm:"-"`

	//out
	Extra_meta      string `xorm:"-"`
	Permission_map  string `xorm:"-"`
	Create_time     int    `xorm:"-"`
	Predecessor_ids []int  `xorm:"-"`
}

func (self *User) TableName() string {
	return "buz_user"
}
