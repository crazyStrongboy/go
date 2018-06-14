package buz

import (
	"eyecool.com/node-retrieval/logic"
	"fmt"
	"eyecool.com/node-retrieval/model"
	"eyecool.com/node-retrieval/utils"
	"time"
	"strconv"
)

type RegionRequest struct {
	Id            string
	PredecessorId string `json:"predecessor_id"`
	Name          string
}

type RegionResponse struct {
	Sets    []*logic.Sets `json:"sets"`
	Rtn     int           `json:"rtn"`
	Message string        `json:"message"`
}

type InsertRegionResponse struct {
	Id      string `json:"id"`
	Rtn     int    `json:"rtn"`
	Message string `json:"message"`
}

//查询摄像机区域
func QueryRegion() *RegionResponse {
	result := &RegionResponse{}
	sets, err := logic.DefaultRegion.QueryRegion()
	if err != nil {
		result.Rtn = -1
		result.Message = "查询失败！"
		return result
	}
	fmt.Println(sets)
	result.Sets = sets
	result.Rtn = 0
	result.Message = "查询成功!"
	return result
}

//插入摄像头区域
func InsertRegion(region *RegionRequest) *InsertRegionResponse {
	result := &InsertRegionResponse{}
	parentId, clusterId, err := utils.GetIdAndClusterId(region.PredecessorId)
	if err != nil || parentId == -2 || clusterId == -2 {
		result.Rtn = -1
		result.Message = "参数错误!"
		return result
	}
	_, reg := logic.DefaultRegion.FindByPrimaryKey(parentId)
	r := &model.Region{
		Name:       region.Name,
		Mlevel:     reg.Mlevel + 1,
		Status:     0,
		ClusterId:  clusterId,
		ParentId:   parentId,
		UpdateTime: time.Now(),
		CreateTime: time.Now().Unix(),
	}
	err = logic.DefaultRegion.InsertRegion(r)
	if err != nil {
		result.Rtn = -1
		result.Message = "插入失败!"
		return result
	}
	id := strconv.Itoa(r.Id)
	result.Id = id
	result.Rtn = 0
	result.Message = "插入成功!"
	return result

}

//更新摄像头区域
func UpdateRegion(region *RegionRequest) *model.RespMsg {
	result := &model.RespMsg{}
	id, clusterId, err := utils.GetIdAndClusterId(region.Id)
	if err != nil || id == -2 || clusterId == -2 {
		result.Rtn = -1
		result.Message = "参数错误!"
		return result
	}
	r := &model.Region{
		ClusterId:  clusterId,
		Id:         id,
		UpdateTime: time.Now(),
		Name:       region.Name,
	}
	err = logic.DefaultRegion.UpdateRegion(r)
	if err != nil {
		result.Rtn = -1
		result.Message = "更新错误!"
		return result
	}
	result.Rtn = 0
	result.Message = "更新成功!"
	return result
}

//删除摄像头区域
func DeleteRegion(regionId string) *model.RespMsg {
	result := &model.RespMsg{}
	id, _, err := utils.GetIdAndClusterId(regionId)
	if err != nil || id == -2 {
		result.Rtn = -1
		result.Message = "参数错误!"
		return result
	}
	has, _ := logic.DefaultRegion.FindByPrimaryKey(id)
	if !has {
		result.Rtn = -1
		result.Message = "该数据不存在!"
		return result
	}
	err = logic.DefaultRegion.DeleteRegion(id)
	if err != nil {
		result.Rtn = -1
		result.Message = "删除失败"
		return result
	}
	result.Rtn = 0
	result.Message = "删除成功"
	return result
}
