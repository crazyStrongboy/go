package controller

import (
	"github.com/emicklei/go-restful"
	"log"
	"encoding/json"
	"eyecool.com/node-retrieval/http/buz"
	"io/ioutil"
	"fmt"
	"eyecool.com/node-retrieval/utils"
	"eyecool.com/node-retrieval/model"
)



type CameraController struct {

}



func (this *CameraController)InsertCamera(req *restful.Request,res *restful.Response){
	log.Print("Received BusinessController.InsertCamera API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	fmt.Println(req.Request.Header)
	cacheMap:=utils.CacheMap{}
	flag:=cacheMap.CheckSession(sessionId)
	result:=&buz.InsertCameraResponse{}
	if flag{
		camera :=buz.CameraRequest{}
		body, _ := ioutil.ReadAll(req.Request.Body)
		err:=json.Unmarshal(body,&camera)
		if err != nil {
			log.Println("InsertCamera err:",err)
			result.Rtn=-1
			result.Message="参数错误！"
		}else{
			//数据入库
			result=buz.InsertCamera(&camera)
		}
	}else{
		result.Rtn=-1
		result.Message="用户未登录！"
	}
	fmt.Println(req.Request.Method)
	SetResponse(res)
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)

}

func (this *CameraController)CameraQuery(req *restful.Request,res *restful.Response){
	log.Print("Received CameraController.CameraQuery API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	flag:=cacheMap.CheckSession(sessionId)
	//flag=true
	result:=&buz.CameraResponse{}
	if flag{
		//查询数据库
		result=buz.CameraQuery()
	}else{
		result.Rtn=-1
		result.Message="用户未登录"
	}
	fmt.Println(req.Request.Method)
	SetResponse(res)
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}

func (this *CameraController)DeleteCamera(req *restful.Request,res *restful.Response){
	log.Print("Received CameraController.DeleteCamera API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	flag:=cacheMap.CheckSession(sessionId)
	//flag=true
	result:=&model.RespMsg{}
	if flag{
		m:=req.Request.URL.Query()
		id:=m.Get("id")
		fmt.Println(id)
		if id == "" {
			result.Rtn=-1
			result.Message="id不能为空!"
			SetResponse(res)
			responseBytes, _ := json.Marshal(result)
			res.ResponseWriter.Write(responseBytes)
			return
		}
		result=buz.DeleteCamera(id)
	}else{
		result.Rtn=-1
		result.Message="用户未登录"
	}
	fmt.Println(req.Request.Method)
	SetResponse(res)
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}

func (this *CameraController)UpdateCamera(req *restful.Request,res *restful.Response){
	log.Print("Received CameraController.UpdateCamera API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	flag:=cacheMap.CheckSession(sessionId)
	//flag=true
	result:=&model.RespMsg{}
	if flag{
		camera :=buz.CameraRequest{}
		body, _ := ioutil.ReadAll(req.Request.Body)
		err:=json.Unmarshal(body,&camera)
		if err != nil {
			fmt.Println("Unmarshal err : ", err)
			result.Rtn=-1
			result.Message = "参数错误！"
		}else{
			//更新
			result=buz.UpdateCamera(&camera)
		}
	}else{
		result.Rtn=-1
		result.Message="用户未登录"
	}
	fmt.Println(req.Request.Method)
	SetResponse(res)
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}





