package buz

import (
	"eyecool.com/node-retrieval/logic"
	"strconv"
)

type ClusterService struct {
}
type ClusterResponse struct {
	Rtn        int       `json:"rtn"`                   //接收状态。0表示接收正常，非0表示接收异常（<0表示错误，>0表示警告）
	Message    string    `json:"message,omitempty"`     //接收状态描述
	ClusterId  string    `json:"cluster_id,omitempty"`  //集群ID
	ClusterIds []*string `json:"cluster_ids,omitempty"` //集群IDs
}

var clusterLogic = new(logic.ClusterLogic)

func (this *ClusterService) GetClusterArray() *ClusterResponse {
	response := new(ClusterResponse)
	clusterIds := make([]*string, 0)
	clusterIntIds := clusterLogic.FindClusters()
	if len(clusterIntIds) > 0 {
		for _, v := range clusterIntIds {
			clusterId := strconv.Itoa(v)
			clusterIds = append(clusterIds, &clusterId)
		}
	}
	response.Rtn = 0
	response.Message = "获取成功"
	response.ClusterIds = clusterIds
	return response
}
