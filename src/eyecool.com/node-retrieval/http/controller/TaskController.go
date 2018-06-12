package controller

import (
	"encoding/json"
	"eyecool.com/node-retrieval/http/buz"
	"github.com/emicklei/go-restful"
	"log"
	"eyecool.com/node-retrieval/utils"
	"io/ioutil"
	"eyecool.com/node-retrieval/model"
	"fmt"
)

type TaskController struct {}


func (this *TaskController)InsertTask(req *restful.Request,res *restful.Response){
	log.Println("Received TaskController.InsertTask API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	//检查用户是否登陆
	flag:=cacheMap.CheckSession(sessionId)
	flag=true
	result:=&buz.InsertTaskResponse{}
	if flag{
		r:=buz.TaskRequest{}
		body,_:=ioutil.ReadAll(req.Request.Body)
		err:=json.Unmarshal(body,&r)
		if err!=nil{
			log.Println("InsertRepository err:",err)
			result.Rtn=-1
			result.Message="参数错误！"
			responseBytes, _ := json.Marshal(result)
			res.ResponseWriter.Write(responseBytes)
		}else{
			//入库
			result=buz.InsertTask(&r)
		}


	}else{

		result.Rtn=-1
		result.Message="用户未登录"
	}

	fmt.Println(req.Request.Method)
	res.Header().Set("Access-Control-Allow-Origin","*")
	res.Header().Set("Access-Control-Allow-Methods","POST,GET,DELETE,PUT")
	res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	res.Header().Set("Access-Control-Max-Age", "1800");//30 min
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}


func (this *TaskController) DeleteTask(req *restful.Request,res *restful.Response){
	log.Print("Received TaskController.DeleteTask API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	flag:=cacheMap.CheckSession(sessionId)
	flag=true
	result:=&model.RespMsg{}
	if flag{
		m:=req.Request.URL.Query()
		id:=m.Get("id")
		//删除
		result=buz.DeleteTask(id)
	}else{

		result.Rtn=-1
		result.Message="用户未登录"
	}

	fmt.Println(req.Request.Method)
	res.Header().Set("Access-Control-Allow-Origin","*")
	res.Header().Set("Access-Control-Allow-Methods","POST,GET,DELETE,PUT")
	res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	res.Header().Set("Access-Control-Max-Age", "1800");//30 min
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}


func (this *TaskController) QueryTask(req *restful.Request,res *restful.Response){
	log.Print("Received TaskController.QueryTask API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	//判断用户是否登陆
	flag:=cacheMap.CheckSession(sessionId)
	flag=true
	result:=&buz.TaskResponse{}
	if flag{
		//查询数据库
		result=buz.QueryTask()
	}else{

		result.Rtn=-1
		result.Message="用户未登录"
	}

	fmt.Println(req.Request.Method)
	res.Header().Set("Access-Control-Allow-Origin","*")
	res.Header().Set("Access-Control-Allow-Methods","POST,GET,DELETE,PUT")
	res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	res.Header().Set("Access-Control-Max-Age", "1800");//30 min
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}

func (this *TaskController)UpdateTask(req *restful.Request,res *restful.Response){
	log.Println("Received TaskController.UpdateTask API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	//检查用户是否登陆
	flag:=cacheMap.CheckSession(sessionId)
	flag=true
	result:=&model.RespMsg{}
	if flag{
		r:=buz.TaskUpdateRequest{}
		body,_:=ioutil.ReadAll(req.Request.Body)
		err:=json.Unmarshal(body,&r)
		if err!=nil{
			log.Println("UpdateTask err:",err)
			result.Rtn=-1
			result.Message="参数错误！"
		}else{
			result=buz.UpdateTask(&r)
		}

	}else{

		result.Rtn=-1
		result.Message="用户未登录"
	}

	fmt.Println(req.Request.Method)
	res.Header().Set("Access-Control-Allow-Origin","*")
	res.Header().Set("Access-Control-Allow-Methods","POST,GET,DELETE,PUT")
	res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	res.Header().Set("Access-Control-Max-Age", "1800");//30 min
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}