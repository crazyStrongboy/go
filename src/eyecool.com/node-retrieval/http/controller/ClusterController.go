package controller

import (
	"github.com/emicklei/go-restful"
	"eyecool.com/node-retrieval/http/buz"
	"encoding/json"
	"strconv"
)

type ClusterController struct {
}
var clusterService = new(buz.ClusterService)
func (this *ClusterController) GetSelfClusterId(req *restful.Request, rsp *restful.Response) {
	response := new(buz.ClusterResponse)
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		clusterId := user.ClusterId
		if clusterId != 0 {
			response.Rtn = 0
			response.Message = "获取成功!"
			response.ClusterId = strconv.Itoa(clusterId)
		}else{
			response.Rtn = -1
			response.Message = "获取集群号失败!"
		}
	}else{
		response.Rtn = -1
		response.Message = "用户未登录!"
	}
	SetResponse(rsp)
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}


func (this *ClusterController) GetClusterIds(req *restful.Request, rsp *restful.Response) {
	response := new(buz.ClusterResponse)
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response = clusterService.GetClusterArray()
	}else{
		response.Rtn = -1
		response.Message = "用户未登录!"
	}
	SetResponse(rsp)
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}
