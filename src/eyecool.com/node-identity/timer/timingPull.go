package timer

import (
	"github.com/robfig/cron"
	"github.com/dghubble/sling"
	. "eyecool.com/node-identity/global"
	"log"
)

func init() {
	go func() {
		cron := cron.New()
		defer cron.Stop()
		cron.AddFunc(Cron, HttpPost)
		cron.Start()
		log.Println("start pull people info from server..........")
		select {}
	}()

}

func HttpPost() {
	s := sling.New()
	cameraInfo := &CameraInfo{
		CameraIp: CameraIp,
	}
	//log.Println("HttpPost..........")
	request, err := s.Post(GetPeopleInfoUrl).Set("Content-Type", "application/json").BodyJSON(cameraInfo).Request()
	//log.Println(request)
	if err != nil {
		log.Println("HttpPost request err:", err)
	}
	responseResult := ResponseResults{}
	_, err = s.Do(request, &responseResult, nil)
	if err != nil {
		log.Println("HttpPost response err :", err)
	}
	if responseResult.Code == "0000" {
		if len(responseResult.Result) > 0 {
			for _, result := range responseResult.Result {
				AddOrRemove(result)
			}
		}
	}
	log.Println("PeopleInfo:", PeopleInfo)
}
