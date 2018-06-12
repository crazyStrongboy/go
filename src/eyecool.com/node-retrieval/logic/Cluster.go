package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
)

type ClusterLogic struct {
}

func (this *ClusterLogic) GetClusterArray() []int {
	clusterIntIds := make([]int, 0)
	MasterDB.Table(new(model.Cluster)).Cols("id").Find(&clusterIntIds)
	return clusterIntIds
}
