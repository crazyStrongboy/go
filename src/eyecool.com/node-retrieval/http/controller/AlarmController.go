package controller

import (
	"github.com/emicklei/go-restful"
	"eyecool.com/node-retrieval/http/buz"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type AlarmController struct {
}

var alarmService = new(buz.AlarmService)
// 命中接口:告警查询接口   /hit/alert
func (this *AlarmController) HitAlert(req *restful.Request, rsp *restful.Response) {
	request := new(buz.AlarmRequest)
	response := new(buz.AlarmResponse)

	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, request)
	if err != nil {
		fmt.Println("HitAlert Unmarshal  err : ", err, ":", request)
		response.Rtn = -1
		response.Message = err.Error()
		SetResponse(rsp)
		responseBytes, _ := json.Marshal(response)
		rsp.ResponseWriter.Write(responseBytes)
		return
	}

	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response = alarmService.HitAlert(request)
	} else {
		response.Rtn = -1
		response.Message = "用户未登录!"
	}
	SetResponse(rsp)
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}
