package buz

import (
	"eyecool.com/node-retrieval/logic"
	"eyecool.com/node-retrieval/model"
	"time"
	"strconv"
	"strings"
	"eyecool.com/node-retrieval/utils"
)

type VideoResponse struct {
	Rtn     int            `json:"rtn"`
	Message string         `json:"message"`
	Videos  []*logic.Video `json:"videos"`
}

type InsertVideoResponse struct {
	Id      string `json:"id"`
	Rtn     int    `json:"rtn"`
	Message string `json:"message"`
}

type VideoRequest struct {
	Id        string
	Name      string
	Url       string
	Enabled   int
	RecParams string `json:"rec_params"`
	ExtraMeta string `json:"extra_meta"`
}

func QueryVideo() *VideoResponse {
	result := &VideoResponse{}
	videos, err := logic.DefaultVideo.QueryVideo()
	if err != nil {
		result.Rtn = -1
		result.Message = "查询失败！"
		return result
	}
	result.Videos = videos
	result.Rtn = 0
	result.Message = "查询成功!"
	return result
}

func InsertVideo(video *VideoRequest) *InsertVideoResponse {
	result := &InsertVideoResponse{}
	if video.Name == "" {
		result.Rtn = -1
		result.Message = "库名不能为空!"
		return result
	}
	if length := strings.Count(video.Name, "") - 1; length > 128 {
		result.Message = "库名不能大于128个字符"
		result.Rtn = -1
		return result
	}
	if video.Url == "" {
		result.Rtn = -1
		result.Message = "url不能为空!"
		return result
	}
	v := &model.Video{
		Enabled:    1,
		ExtraMeta:  video.ExtraMeta,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now(),
		Name:       video.Name,
		Url:        video.Url,
	}
	err := logic.DefaultVideo.InsertVideo(v)
	if err != nil {
		result.Rtn = -1
		result.Message = "添加失败"
		return result
	}
	id := strconv.Itoa(v.PkId)
	result.Rtn = 0
	result.Message = "添加成功"
	result.Id = id
	return result
}

func UpdateVideo(v *VideoRequest) *model.RespMsg {
	result := &model.RespMsg{}
	if length := strings.Count(v.Name, "") - 1; length > 128 {
		result.Message = "库名不能大于128个字符"
		result.Rtn = -1
		return result
	}
	pkId, clusterId, err := utils.GetIdAndClusterId(v.Id)
	if err != nil {
		result.Rtn = -1
		result.Message = "参数错误"
		return result
	}
	has, video := logic.DefaultVideo.FindVideoByPkId(pkId)
	if !has {
		result.Rtn = -1
		result.Message = "该库不存在!"
		return result
	}
	video.UpdateTime = time.Now()
	video.ClusterId = clusterId
	video.Enabled = v.Enabled
	video.RecParams = v.RecParams
	video.ExtraMeta = v.ExtraMeta
	if v.Url != "" {
		video.Url = v.Url
	}
	if v.Name != "" {
		video.Name = v.Name
	}

	err = logic.DefaultVideo.UpdateVideo(video)
	if err != nil {
		result.Message = "更新失败"
		result.Rtn = -1
		return result
	}
	result.Message = "更新成功"
	result.Rtn = 0
	return result
}

func DeleteVideo(id string) *model.RespMsg {
	result := &model.RespMsg{}
	pkId, _, err := utils.GetIdAndClusterId(id)
	if err != nil || pkId == -2 {
		result.Rtn = -1
		result.Message = "id不规范!"
		return result
	}
	has, _ := logic.DefaultVideo.FindVideoByPkId(pkId)
	if !has {
		result.Rtn = -1
		result.Message = "您要删除的视频不存在!"
		return result
	}
	video := &model.Video{
		PkId: pkId,
	}
	err = logic.DefaultVideo.DeleteVideo(video)
	if err != nil {
		result.Rtn = -1
		result.Message = "删除失败"
		return result
	}
	result.Message = "删除成功"
	result.Rtn = 0
	return result

}
