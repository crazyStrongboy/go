package controller

import (
	"github.com/emicklei/go-restful"
	"log"
	"encoding/json"
	"eyecool.com/node-retrieval/http/buz"
	"io/ioutil"
	"fmt"
	"eyecool.com/node-retrieval/utils"
)

type CameraController struct {
}

func (this *CameraController) InsertCamera(req *restful.Request, res *restful.Response) {
	log.Print("Received BusinessController.InsertCamera API request : ", req.Request.RemoteAddr)
	sessionId := req.HeaderParameter("session_id")
	cacheMap := utils.CacheMap{}
	flag := cacheMap.CheckSession(sessionId)
	flag = true
	if flag {
		user := cacheMap.GetUserSession(sessionId)
		camera := buz.CameraRequest{}
		body, _ := ioutil.ReadAll(req.Request.Body)
		err := json.Unmarshal(body, &camera)
		if err != nil {
			fmt.Println("Unmarshal err : ", err)
		}
		fmt.Println("获取json中的retrievalQuery:", &camera)

		//数据入库
		result := buz.InsertCamera(&camera, user)
		fmt.Println(req.Request.Method)
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		res.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)

	} else {
		result := make(map[string]interface{})
		result["rtn"] = -1
		result["message"] = "用户未登录"
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)
	}

}

func (this *CameraController) CameraQuery(req *restful.Request, res *restful.Response) {
	log.Print("Received BusinessController.CameraQuery API request : ", req.Request.RemoteAddr)
	result := buz.CameraQuery()
	responseBytes, _ := json.Marshal(result)
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "GET")
	res.ResponseWriter.Write(responseBytes)
}

func (this *CameraController) DeleteCamera(req *restful.Request, res *restful.Response) {
	log.Print("Received BusinessController.DeleteCamera API request : ", req.Request.RemoteAddr)
	m := req.Request.URL.Query()
	id := m.Get("id")
	fmt.Println(id)
	result := buz.DeleteCamera(id)

	fmt.Println(req.Request.Method)
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	res.Header().Set("Access-Control-Max-Age", "1800"); //30 min

	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}

func (this *CameraController) UpdateCamera(req *restful.Request, res *restful.Response) {
	log.Print("Received BusinessController.UpdateCamera API request : ", req.Request.RemoteAddr)
	response := new(buz.CameraResponse)
	camera := buz.CameraRequest{}
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &camera)
	if err != nil {
		fmt.Println("Unmarshal err : ", err)
		response.Rtn = -1
		response.Message = err.Error()
		responseBytes, _ := json.Marshal(response)
		res.ResponseWriter.Write(responseBytes)
		return
	}
	result := buz.UpdateCamera(&camera)
	fmt.Println(req.Request.Method)
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	res.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}
