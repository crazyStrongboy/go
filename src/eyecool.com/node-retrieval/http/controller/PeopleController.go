package controller

import (
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"eyecool.com/node-retrieval/http/buz"
)

type PeopleController struct {
}

var peopleService = new(buz.PeopleService)
//导入图片
func (this *PeopleController) PictureSynchronized(req *restful.Request, rsp *restful.Response) {
	response := new(buz.PeopleResponse)
	request := new(buz.PeopleRequest)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, request)
	if err != nil {
		fmt.Println("PictureSynchronized Unmarshal  err : ", err, ":", request)
		response.Rtn = -1
		response.Message = err.Error()
		responseBytes, _ := json.Marshal(response)
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		rsp.ResponseWriter.Write(responseBytes)
		return
	}
	sessionId := req.HeaderParameter("session_id")
	if request.Picture_image_content_base64 != "" {
		user := cacheMap.GetUserSession(sessionId)
		if user != nil {
			response = peopleService.Insert(request, user.Id)
		} else {
			response.Rtn = -1
			response.Message = "用户未登录!"
		}
	} else {
		response.Rtn = -1
		response.Message = "图片base64 不能为空"
	}
	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min

	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

//修改/face/update POST
func (this *PeopleController) FaceUpdate(req *restful.Request, rsp *restful.Response) {
	response := new(buz.PeopleResponse)
	request := new(buz.PeopleRequest)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, request)
	if err != nil {
		fmt.Println("PictureSynchronized Unmarshal  err : ", err, ":", request)
		response.Rtn = -1
		response.Message = err.Error()
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		return
	}
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response = peopleService.Update(request)
	} else {
		response.Rtn = -1
		response.Message = "用户未登录!"
	}

	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

//删除图片/face/delete POST
func (this *PeopleController) FaceDelete(req *restful.Request, rsp *restful.Response) {
	response := new(buz.PeopleResponse)
	request := new(buz.PeopleRequest)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, request)
	if err != nil {
		fmt.Println("PictureSynchronized Unmarshal  err : ", err, ":", request)
		response.Rtn = -1
		response.Message = err.Error()
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		return
	}
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response = peopleService.Delete(request)
	} else {
		response.Rtn = -1
		response.Message = "用户未登录!"
	}
	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}
