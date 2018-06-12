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

type RepositoryController struct{}

func (this *RepositoryController) QueryRepository(req *restful.Request, res *restful.Response) {
	log.Print("Received RepositoryController.QueryRepository API request : ", req.Request.RemoteAddr)
	sessionId := req.HeaderParameter("session_id")
	cacheMap := utils.CacheMap{}
	//判断用户是否登陆
	flag := cacheMap.CheckSession(sessionId)
	flag = true
	if flag {
		//查询数据库
		result := buz.QueryRepository()
		fmt.Println(req.Request.Method)
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		res.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)
	} else {
		result := &buz.RepositoryResponse{}
		result.Rtn = -1
		result.Message = "用户未登录"
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)
	}
}

func (this *RepositoryController) InsertRepository(req *restful.Request, res *restful.Response) {
	log.Println("Received RepositoryController.QueryRepository API request : ", req.Request.RemoteAddr)
	sessionId := req.HeaderParameter("session_id")
	cacheMap := utils.CacheMap{}
	//检查用户是否登陆
	flag := cacheMap.CheckSession(sessionId)
	flag = true
	if flag {
		user := cacheMap.GetUserSession(sessionId)
		r := buz.RepositoryRequest{}
		body, _ := ioutil.ReadAll(req.Request.Body)
		err := json.Unmarshal(body, &r)
		if err != nil {
			log.Println("InsertRepository err:", err)
			result := &buz.InsertRepositoryResponse{}
			result.Rtn = -1
			result.Message = "参数错误！"
			responseBytes, _ := json.Marshal(result)

			res.ResponseWriter.Write(responseBytes)
		}
		//入库
		result := buz.InsertRepository(&r, user)
		fmt.Println(req.Request.Method)
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		res.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)

	} else {
		result := &buz.InsertRepositoryResponse{}
		result.Rtn = -1
		result.Message = "用户未登录"
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)
	}

}

func (this *RepositoryController) UpdateRepository(req *restful.Request, res *restful.Response) {
	log.Println("Received RepositoryController.UpdateRepository API request : ", req.Request.RemoteAddr)
	sessionId := req.HeaderParameter("session_id")
	cacheMap := utils.CacheMap{}
	//检查用户是否登陆
	flag := cacheMap.CheckSession(sessionId)
	flag = true
	if flag {
		r := buz.RepositoryRequest{}
		body, _ := ioutil.ReadAll(req.Request.Body)
		err := json.Unmarshal(body, &r)
		if err != nil {
			log.Println("InsertRepository err:", err)
			result := &model.RespMsg{}
			result.Rtn = -1
			result.Message = "参数错误！"
			responseBytes, _ := json.Marshal(result)
			res.ResponseWriter.Write(responseBytes)
		}
		result := buz.UpdateRepository(&r)
		fmt.Println(req.Request.Method)
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		res.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)
	} else {
		result := &model.RespMsg{}
		result.Rtn = -1
		result.Message = "用户未登录"
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)
	}
}

func (this *RepositoryController) DeleteRepository(req *restful.Request, res *restful.Response) {
	log.Print("Received RepositoryController.DeleteRepository API request : ", req.Request.RemoteAddr)
	sessionId := req.HeaderParameter("session_id")
	cacheMap := utils.CacheMap{}
	flag := cacheMap.CheckSession(sessionId)
	flag = true
	if flag {
		m := req.Request.URL.Query()
		id := m.Get("id")
		//删除
		result := buz.DeleteRepository(id)
		fmt.Println(req.Request.Method)
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		res.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)
	} else {
		result := &buz.InsertRegionResponse{}
		result.Rtn = -1
		result.Message = "用户未登录"
		responseBytes, _ := json.Marshal(result)
		res.ResponseWriter.Write(responseBytes)
	}
}
