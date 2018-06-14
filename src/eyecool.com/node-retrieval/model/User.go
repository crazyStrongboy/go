package model

import "time"

type User struct {
	Id            int       `json:"id" xorm:"pk autoincr 'id'"`
	Name          string    `json:"name"`
	Password      string    `json:"password"`
	ExtraMeta     string    `json:"extra_meta"`
	Status        int       `json:"omitempty"`
	GroupId       int       `json:"omitempty"`
	ClusterId     int       `json:"omitempty"`
	RepositoryId  int       `json:"omitempty"`
	CreateTime    int64     `json:"create_time"`
	UserLevel     int       `json:"omitempty"`
	ParentId      int       `json:"omitempty"`
	PermissionMap string    `json:"permission_map"`
	UpdateTime    time.Time `json:"omitempty"`
	Param1        string    `json:"omitempty"`
	Param2        string    `json:"omitempty"`
	Param3        string    `json:"omitempty"`
	Param4        string    `json:"omitempty"`

	Predecessor_id string `json:"predecessor_ids" xorm:"-"`
}

func (self *User) TableName() string {
	return "buz_user"
}
