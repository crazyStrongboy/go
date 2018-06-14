package model

import "time"

type TaskTuple struct {
	TaskId           string  `json:"TaskId"  xorm:"TaskId"`
	TaskName         string  `json:"TaskName"  xorm:"TaskName"`
	TaskChildrenPkId int     `json:"TaskChildrenPkId"  xorm:"TaskChildrenPkId"`
	CameraId         string  `json:"CameraId"  xorm:"CameraId"`
	RepositoryId     string  `json:"RepositoryId"  xorm:"RepositoryId"`
	Threshold        float32 `json:"Threshold"  xorm:"Threshold"`
	Topk             int
}

type Task struct {
	PkId       int       `json:"pk_id" xorm:"pk autoincr"`
	Id         string
	Name       string
	Status     int
	ClusterId  int
	UpdateTime time.Time `json:"created_at" `
	CreateTime int64     `json:"created_at" `
	Param1     string
	Param2     string
	Param3     string
	Param4     string
}

func (self *Task) TableName() string {
	return "buz_task"
}

type TaskChildren struct {
	PkId         int       `json:"pk_id" xorm:"pk autoincr"`
	Id           string
	TaskId       string
	CameraId     string
	CameraIp     string
	Name         string    `xorm:"-"`
	TaskPkId     int       `xorm:"-"`
	CameraPkId   int       `xorm:"-"`
	RepositoryId string
	Threshold    float64
	ExtraMeta    string
	Status       int
	ClusterId    int
	UpdateTime   time.Time `json:"created_at"`
	CreateTime   int64     `json:"created_at" `
	Param1       string
	Param2       string
	Param3       string
	Param4       string
}

func (self *TaskChildren) TableName() string {
	return "buz_task_children"
}
