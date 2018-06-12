package global

import (
	. "github.com/polaris1119/config"
)

type CameraInfo struct {
	CameraIp string `json:"cameraIp,omitempty"`
}
type ResponseResults struct {
	Result []*ResponseResult
	Code   string
}
type ResponseResult struct {
	CameraIp      string `json:"cameraIp,omitempty"`
	PeopleId      int64
	PeopleIdValue string
	RecordStatus  int // 1:新增 2:删除
}

var (
	PeopleInfo       map[string]int64
	HeartUrl         = "http://192.168.0.39:8080/admin/interface/heartbeat"
	GetPeopleInfoUrl = "http://192.168.0.39:8080/admin/interface/getPeopleInfosByCameraIP"
	CameraIp         = "192.163.0.211"
	Cron             = "*/5 * * * * ?"
)

func init() {
	HeartUrl, _ = ConfigFile.GetValue("camera", "heart_url")
	GetPeopleInfoUrl, _ = ConfigFile.GetValue("camera", "get_people_info_url")
	CameraIp, _ = ConfigFile.GetValue("camera", "camera_ip")
	Cron, _ = ConfigFile.GetValue("camera", "cron")

	PeopleInfo = make(map[string]int64, 0)
}

func AddOrRemove(result *ResponseResult) {
	peopleIdValue := result.PeopleIdValue
	status := result.RecordStatus
	peopleId := result.PeopleId
	if status == 1 {
		PeopleInfo[peopleIdValue] = peopleId
	}
	if status == 2 {
		delete(PeopleInfo, peopleIdValue)
	}
}