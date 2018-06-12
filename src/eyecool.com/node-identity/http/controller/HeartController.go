package controller

import (
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"encoding/json"
	"log"
	"github.com/dghubble/sling"
	"eyecool.com/node-identity/global"
)

type HeartController struct {
}

type HeartRequest struct {
	CameraIp    string `json:"cameraIp,omitempty"`    //相机Ip
	CaptureTime int64  `json:"captureTime,omitempty"` //当前时间毫秒值
	Heartbeat   int    `json:"heartbeat,omitempty"`
}

//-1 摄像机IP不能为空
//0 成功
//1 缓存中有此ip键记录但是没有此摄像机对象
//2 数据库中没有此ip记录
//3 异常报错
//4 参数错误
type HeartResponse struct {
	Status int `json:"status"` //心跳响应状态码
}

func (this *HeartController) HeartBeat(req *restful.Request, rsp *restful.Response) {
	request := new(HeartRequest)
	response := new(HeartResponse)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Println("HeartBeart Unmarshal err:", err)
		response.Status = 4
		responseBytes, _ := json.Marshal(response)
		rsp.ResponseWriter.Write(responseBytes)
		return
	}
	log.Println(request)
	response = sendHeart(request, response)
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

func sendHeart(req *HeartRequest, resp *HeartResponse) *HeartResponse {
	sling := sling.New()
	request, err := sling.Post(global.HeartUrl).Set("Content-Type", "application/json").BodyJSON(req).Request()
	if err != nil {
		log.Println("build request err:", err)
	}
	_, err = sling.Do(request, resp, nil)
	if err != nil {
		log.Println("send heart err :", err)
	}
	return resp
}
