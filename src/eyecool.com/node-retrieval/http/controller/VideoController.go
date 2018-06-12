package controller

import (
	"github.com/emicklei/go-restful"
	"eyecool.com/node-retrieval/utils"
	"eyecool.com/node-retrieval/http/buz"
	"encoding/json"
	"log"
	"io/ioutil"
	"eyecool.com/node-retrieval/model"
	"fmt"
)

type VideoController struct {}

func (this *VideoController)QueryVideo(req *restful.Request,res *restful.Response){
	log.Print("Received VideoController.QueryVideo API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	result:=&buz.VideoResponse{}
	//判断用户是否登陆
	flag:=cacheMap.CheckSession(sessionId)
	flag=true
	if flag{
		//查询数据库
		result=buz.QueryVideo()
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

func (this *VideoController)InsertVideo(req *restful.Request,res *restful.Response){
	log.Println("Received RepositoryController.QueryRepository API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	//检查用户是否登陆
	flag:=cacheMap.CheckSession(sessionId)
	flag=true
	result:=&buz.InsertVideoResponse{}
	if flag{
		v:=buz.VideoRequest{}
		body,_:=ioutil.ReadAll(req.Request.Body)
		err:=json.Unmarshal(body,&v)
		if err!=nil{
			log.Println("InsertRepository err:",err)
			result:=&buz.InsertVideoResponse{}
			result.Rtn=-1
			result.Message="参数错误！"
		}else{
			//入库
			result=buz.InsertVideo(&v)
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


func (this *VideoController)UpdateVideo(req *restful.Request,res *restful.Response){
	log.Println("Received VideoController.UpdateVideo API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	//检查用户是否登陆
	flag:=cacheMap.CheckSession(sessionId)
	flag=true
	result:=&model.RespMsg{}
	if flag{
		r:=buz.VideoRequest{}
		body,_:=ioutil.ReadAll(req.Request.Body)
		err:=json.Unmarshal(body,&r)
		if err!=nil{
			log.Println("UpdateVideo err:",err)
			result.Rtn=-1
			result.Message="参数错误！"
		}else{
			result=buz.UpdateVideo(&r)
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

func (this *VideoController) DeleteVideo(req *restful.Request,res *restful.Response){
	log.Print("Received VideoController.DeleteVideo API request : ", req.Request.RemoteAddr)
	sessionId:=req.HeaderParameter("session_id")
	cacheMap:=utils.CacheMap{}
	flag:=cacheMap.CheckSession(sessionId)
	flag=true
	result:=&model.RespMsg{}
	if flag{
		m:=req.Request.URL.Query()
		id:=m.Get("id")
		//删除
		result=buz.DeleteVideo(id)
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
