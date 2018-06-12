package model

type VideoCamera struct {
	Id int `xorm: "pk autoincr 'id'"`
	CameraId string
	Conntype int
	Status int
	RtspUrl string `xorm:"char 'rtspUrl'"`
	Sipid string `xorm:"char"`
	Describes string
	RtspUrl2 string `xorm:'rtspUrl2'`
	Param1 string
	Param2 string
	Param3 string

}

func (self *VideoCamera)TableName()string  {
	return "video_camera"
}

