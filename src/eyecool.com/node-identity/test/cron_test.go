package test

import (
	"testing"
	"github.com/robfig/cron"
	"log"
	"fmt"
	"eyecool.com/node-identity/global"
	"github.com/dghubble/sling"
)

type IdCard struct {
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
	RecordStatus  int
}


func TestCron(t *testing.T) {
	cron := cron.New()
	defer cron.Stop()

	cron.AddFunc("*/5 * * * * ?", HttpPost)
	cron.Start()
	select {}
}

func TestHttpPost(t *testing.T) {
	HttpPost()
}

func HttpPost() {
	log.Println("send Http post")
	s := sling.New()
	idcard := &IdCard{
		CameraIp: "192.163.0.211",
	}
	request, err := s.Post(global.GetPeopleInfoUrl).Set("Content-Type", "application/json").BodyJSON(idcard).Request()

	if err != nil {
		fmt.Println("request err:", err)
	}
	responseResult := ResponseResults{}
	response, err := s.Do(request, &responseResult, nil)
	log.Println(err, responseResult)
	if err != nil {
		fmt.Println("response err :", err)
	}
	fmt.Println(response)
}
