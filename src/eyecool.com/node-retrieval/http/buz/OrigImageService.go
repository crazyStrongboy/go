package buz

import (
	"eyecool.com/node-retrieval/logic"
	"fmt"
	"strconv"
	"encoding/json"
	"eyecool.com/node-retrieval/utils"
)

type OrigImageService struct {
}

type StartOffset struct {
	Cluster_id int
	Offset     int //不填时为当前最新 (填0表示从头开始)
}
type NextOffset struct {
	ClusterId int `json:"cluster_id,omitempty"`
	Offset    int `json:"offset,omitempty"` //不填时为当前最新 (填0表示从头开始)
}

type OrigImageRequest struct {
	Start_offsets []*StartOffset
	Limit         int //表示每个集群最长返回多少个结果（不能超过5000)
	Face_image_id string
	Extra_fields  []string //导图时导入自定义字段，在检索结果中显示相关字段
}
type OrigResult struct {
	FaceImageId    string `json:"face_image_id,omitempty"`  //过人id($编号@$集群号)
	CameraId       string `json:"camera_id,omitempty"`      //摄像头编号 ($编号@$集群号)
	FaceImageUri   string `json:"face_image_uri,omitempty"` //人脸图URI
	PictureUri     string `json:"picture_uri,omitempty"`    //场景图URI
	FaceRect       *Rect  `json:"face_rect,omitempty"`      //人脸框
	Timestamp      int64  `json:"timestamp,omitempty"`      //过人的时间戳
	BornYear       string `json:"born_year,omitempty"`
	FaceImageIdStr string `json:"face_image_id_str,omitempty"`
	Gender         int    `json:"gender,omitempty"`
	IsWritable     bool   `json:"is_writable,omitempty"`
	Name           string `json:"name,omitempty"`
	Nation         string `json:"nation,omitempty"`
	PersonId       int64  `json:"person_id,omitempty"`
	RepositoryId   int    `json:"repository_id,omitempty"`
	CustomField    string `json:"custom_field,omitempty"`
}
type OrigImageResponse struct {
	Rtn         int           `json:"rtn"`                    //错误码
	Message     string        `json:"message,omitempty"`      //错误消息
	Results     []*OrigResult `json:"results,omitempty"`      //错误消息
	NextOffsets []*NextOffset `json:"next_offsets,omitempty"` //每个集群下一次请求的start_offset
	Total       int           `json:"total,omitempty"`
}

var origImgeLogic = new(logic.OrigImageLogic)

func (this *OrigImageService) GetCaptureImage(request *OrigImageRequest) *OrigImageResponse {
	results := make([]*OrigResult, 0)
	response := new(OrigImageResponse)
	startOffsets := request.Start_offsets
	if len(startOffsets) == 0 {
		response.Rtn = -1
		response.Message = "参数错误！"
		return response
	}
	limit := request.Limit
	if limit == 0 {
		limit = 5000
	} else if limit > 5000 {
		response.Rtn = -1
		response.Message = "limit不能超过5000！"
		return response
	}
	nextOffsets := make([]*NextOffset, 0)
	for _, v := range startOffsets {
		nextOffset := new(NextOffset)
		offSet := v.Offset
		clusterId := v.Cluster_id
		nextOffset.ClusterId = clusterId
		nextOffset.Offset = offSet + limit
		nextOffsets = append(nextOffsets, nextOffset)
		origImages, err := origImgeLogic.FindOrigImages(clusterId, offSet, limit)
		if err != nil {
			fmt.Println("FindOrigImages err", err)
			response.Rtn = -1
			response.Message = "FindOrigImages err！"
			return response
		}
		response.NextOffsets = nextOffsets
		if len(origImages) == 0 {
			response.Rtn = -1
			response.Message = "OrigImages is empty！"
			return response
		}
		for _, origImage := range origImages {
			faceNum := origImage.FaceNum
			faceRect := origImage.FaceRect
			rect := new(Rect)
			rectByte, _ := json.Marshal(faceRect)
			json.Unmarshal(rectByte, rect)

			for i := 0; i < faceNum; i++ {
				origResult := new(OrigResult)
				origResult.FaceImageId = strconv.Itoa(origImage.Id)
				origResult.CameraId = origImage.CameraId
				origResult.FaceImageUri = origImage.FaceImageUri
				origResult.PictureUri = origImage.PictureUri
				origResult.FaceRect = rect
				origResult.Timestamp = origImage.Timestamp
				results = append(results, origResult)
			}
		}
		response.Results = results
		response.Rtn = 0
		response.Message = "查询成功！"
	}
	return response
}

func (this *OrigImageService) GetSingleImage(request *OrigImageRequest) *OrigImageResponse {
	results := make([]*OrigResult, 0)
	response := new(OrigImageResponse)
	imageId, _, err := utils.GetIdAndClusterId(request.Face_image_id)
	if err != nil {
		response.Rtn = -1
		response.Message = "face_image_id参数不合格!"
		return response
	}
	hasF, feature := featureLogic.FindFaceFeatureByPkId(imageId)
	if hasF {
		hasP, people := peopleLogic.FindPeopleById(feature.PeopleId)
		hasI, image := imageLogic.FindImageById(feature.ImageId)
		if hasP && hasI {
			origResult := new(OrigResult)
			origResult.CustomField = people.CustomField
			origResult.BornYear = people.Birthday
			origResult.FaceImageId = request.Face_image_id
			origResult.FaceImageIdStr = request.Face_image_id

			rect := &Rect{
				X: feature.X,
				Y: feature.Y,
				H: feature.H,
				W: feature.W,
			}
			origResult.FaceRect = rect
			origResult.Gender = people.Gender
			origResult.IsWritable = false
			origResult.Name = people.Name
			origResult.Nation = people.Nation
			origResult.PersonId = people.Id
			origResult.PictureUri = image.ImageUri
			origResult.RepositoryId = image.RepositoryId
			origResult.Timestamp = image.CreateTime.Unix()
			results = append(results, origResult)
			response.Results = results
			response.Rtn = 0
			response.Total = 1
			response.Message = "获取成功!"
		}
	}
	return response
}
