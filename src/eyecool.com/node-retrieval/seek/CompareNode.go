package seek

import (
	"eyecool.com/node-retrieval/model"
	"eyecool.com/node-retrieval/logic"
	"eyecool.com/node-retrieval/utils"
	"encoding/json"
	"strconv"
	"eyecool.com/node-retrieval/global"
	"fmt"
)

type CompareNodeLogic struct{}

var DefaultCompareNode = CompareNodeLogic{}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	T int `json:"t"`
	B int `json:"b"`
}

func (CompareNodeLogic) Compare(is *model.ImageSource) error {
	for i := 0; i < is.FaceNum; i++ {
		feats := is.FaceFeatureBufs[i*global.FEATURE_LENGTH:(i+1)*global.FEATURE_LENGTH]
		list, err := CompareTarget(is.RepositoryId, feats, int(is.Topk))
		if err != nil {
			return err
		}
		maxf := utils.Itof32(list.Max())
		if maxf > is.Threshold {
			alarm := &model.AlarmInfo{
				TaskId:       is.TaskId,
				CameraId:     is.CameraId,
				RepositoryId: is.RepositoryId,
				//AlarmPeopleId:         is.PeopleId,
				AlarmImageContextPath: is.ImageContextPath,
				//AlarmCropImageProperties string
				//AlarmCropImageUri string
				//AlarmOrigImageId :
				//AlarmOrigImageRectIdx:
				AlarmOrigImageUri: is.OrigImageUri,
				//AlarmScoreOthers string
				AlarmScore: maxf,
				//AlarmTmplFeatrueId1 int64
				//AlarmTmplFeatrueId2 int64
				//AlarmTmplFeatrueId3 int64
				//AlarmTmplScore1 float32
				//AlarmTmplScore2 float32
				//AlarmTmplScore3 float32
				Timestamp: is.CaptureTime,
				//ClusterId int
			}
			rectmap := make([]map[int]Rect, 0)
			err := json.Unmarshal([]byte(is.FaceRects), &rectmap)
			if err != nil {
				fmt.Println("Unmarshal LifecycleRequest err : ", err)
			} else {
				rect := rectmap[i]
				if rect != nil {
					rectBytes, _ := json.Marshal(rect)
					alarm.AlarmOrigImageRectIdx = string(rectBytes)
				}
			}

			if len := list.Len(); len >= 1 {
				pair := list[0]
				alarm.AlarmPeopleId = pair.Tmpl.PeopleId
				alarm.AlarmTmplFeatrueId1 = pair.Tmpl.FaceImageId
				alarm.AlarmTmplScore1 = utils.Itof32(pair.Similarity)
				if len >= 2 {
					pair = list[1]
					alarm.AlarmTmplFeatrueId2 = pair.Tmpl.FaceImageId
					alarm.AlarmTmplScore2 = utils.Itof32(pair.Similarity)
				}
				if len >= 3 {
					pair = list[2]
					alarm.AlarmTmplFeatrueId3 = pair.Tmpl.FaceImageId
					alarm.AlarmTmplScore3 = utils.Itof32(pair.Similarity)
				}
				if len >= 4 {
					others := list[3:]
					// [{"tmplFeatrueId":286394,"tmplScore":54},{"tmplFeatrueId":286568,"tmplScore":45}]
					otherms := make([]map[string]string, 0, others.Len())
					for i := 0; i < others.Len(); i++ {
						m := map[string]string{}
						m["tmplFeatrueId"] = others[i].Tmpl.FaceImageId
						m["tmplScore"] = strconv.Itoa(others[i].Similarity)
						otherms = append(otherms, m)
					}
					if bytes, err := json.Marshal(otherms); err == nil {
						alarm.AlarmScoreOthers = string(bytes)
					}
				}
			}
			logic.DefaultAlarmInfo.Insert(alarm)
		}
	}
	return nil
}
