package controller

import (
	"github.com/emicklei/go-restful"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"eyecool.com/node-retrieval/http/buz"
)

type ImageFailController struct {
}

var imageFailService = new(buz.ImageFailService)

// 获取导图失败的请求 /repository/picture/failed  POST
func (this *ImageFailController) GetFailImage(req *restful.Request, rsp *restful.Response) {
	response := new(buz.ImageFailResponse)
	request := new(buz.ImageFailRequest)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, request)
	if err != nil {
		fmt.Println("GetFailImage Unmarshal  err : ", err, ":", request)
		response.Rtn = -1
		response.Message = err.Error()
		SetResponse(rsp)
		responseBytes, _ := json.Marshal(response)
		rsp.ResponseWriter.Write(responseBytes)
		return
	}
	if request.Repository_id == "" {
		response.Rtn = -1
		response.Message = "repository_id不能为空!"
		SetResponse(rsp)
		responseBytes, _ := json.Marshal(response)
		rsp.ResponseWriter.Write(responseBytes)
		return
	}

	if request.Start == 0 {
		response.Rtn = -1
		response.Message = "start不能为空!"
		SetResponse(rsp)
		responseBytes, _ := json.Marshal(response)
		rsp.ResponseWriter.Write(responseBytes)
		return
	}

	if request.Limit == 0 {
		response.Rtn = -1
		response.Message = "limit不能为空!"
		SetResponse(rsp)
		responseBytes, _ := json.Marshal(response)
		rsp.ResponseWriter.Write(responseBytes)
		return
	}
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response = imageFailService.GetFailImage(request)
	} else {
		response.Rtn = -1
		response.Message = "用户未登录!"
	}
	SetResponse(rsp)
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}
