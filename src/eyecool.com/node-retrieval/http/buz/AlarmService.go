package buz

import (
	"eyecool.com/node-retrieval/logic"
	"eyecool.com/node-retrieval/model"
	"time"
	"sort"
	"log"
	"encoding/json"
)

type AlarmService struct {
}

type AlarmRequest struct {
	Task_ids      []string
	Camera_ids    []string
	Surveillances []logic.Surveillance
	Extra_fields  []string //导图时导入自定义字段，在检索结果中显示相关字段
	Hit_condition *logic.HitCondition
	Order         string
	Start         int
	Limit         int
}

//type Surveillance struct {
//	Camera_id     string  //某一个特定布控的camera_id
//	Repository_id string  //某一个特定布控的repository_id
//	Threshold     float64 //某一个特定布控的阈值
//}

type AlarmResponse struct {
	Total   int           `json:"total"`
	Rtn     int           `json:"rtn"`
	Message string        `json:"message,omitempty"`
	Results []AlarmResult `json:"pair_results"`
}

type AlarmResult struct {
	RecogResult  RecogResult  `json:"recog_result,omitempty"`  //识别结果
	PeopleResult PeopleResult `json:"people_result,omitempty"` //抓拍结果
	HitDetails   HitDetails   `json:"hit_detail,omitempty"`    //命中详情
}

//目标人员信息
type PeopleResult struct {
	BronYear       string      `json:"born_year,omitempty"`     //出生年月
	FaceImageId    string      `json:"face_image_id,omitempty"` //过人id($编号@$集群号)
	FaceImageIdStr string      `json:"face_image_id_str,omitempty"`
	FaceImageUri   string      `json:"face_image_uri,omitempty"` //人脸图URI
	FaceRect       *model.Rect `json:"face_rect,omitempty"`      //人脸框
	PictureUri     string      `json:"picture_uri,omitempty"`    //场景图URI
	Timestamp      int64       `json:"timestamp,omitempty"`      //过人的时间戳
	Name           string      `json:"name,omitempty"`
	Nation         int         `json:"nation,omitempty"`
	PersonId       string      `json:"person_id,omitempty"`
	Gender         int         `json:"gender,omitempty"` //性别: 1-male, 0-female
	Age            int         `json:"age,omitempty"`    //年龄
	Race           int         `json:"race,omitempty"`   //种族: 0-未知, 1-白种人, 2-黄种人, 3-黑种人
	SmileLevel     int         `json:"smile_level"`      //微笑程度: 0-100
	BeautyLevel    int         `json:"beauty_level"`     //颜值: 0-100
	RepositoryId   string      `json:"repository_id,omitempty"`
	CustomField    string      `json:"custom_field,omitempty"`
}

//认证结果
type CompareResult struct {
	Score     float32 `json:"score"`
	FeatureId string  `json:"feature_id"`
}

type RecogResult struct {
	CameraId       string           `json:"camera_id,omitempty"`     //摄像头编号 ($编号@$集群号)
	FaceImageId    string           `json:"face_image_id,omitempty"` //过人id($编号@$集群号)
	FaceImageIdStr string           `json:"face_image_id_str,omitempty"`
	FaceImageUri   string           `json:"face_image_uri,omitempty"` //人脸图URI
	FaceRect       *model.Rect      `json:"face_rect,omitempty"`      //人脸框
	PictureUri     string           `json:"picture_uri,omitempty"`    //场景图URI
	Timestamp      int64            `json:"timestamp,omitempty"`      //过人的时间戳
	Gender         int              `json:"rec_gender,omitempty"`     //性别: 1-male, 0-female
	Age            int              `json:"rec_age,omitempty"`        //年龄
	Race           int              `json:"rec_race,omitempty"`       //种族: 0-未知, 1-白种人, 2-黄种人, 3-黑种人
	SmileLevel     int              `json:"rec_smile_level"`          //微笑程度: 0-100
	BeautyLevel    int              `json:"rec_beauty_level"`         //颜值: 0-100
	RepositoryId   string           `json:"repository_id,omitempty"`
	CustomField    string           `json:"custom_field,omitempty"`
	IsHit          bool             `json:"is_hit,omitempty"` //是否命中
	TopN           []*CompareResult `json:"topN,omitempty"`   //是否命中
}

type HitDetails struct {
	FaceImageId     string  `json:"face_image_id,omitempty"`
	HitFaceImageId  string  `json:"hit_face_image_id,omitempty"`
	HitRepositoryId string  `json:"hit_repository_id,omitempty"`
	HitSimilarity   float32 `json:"hit_similarity,omitempty"`
	HitType         int     `json:"hit_type,omitempty"`
	Id              int     `json:"id,omitempty"`
	Timestamp       int64   `json:"timestamp,omitempty"`
}

type CompareResults []*CompareResult

func (a CompareResults) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a CompareResults) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

func (infos CompareResults) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return infos[j].Score < infos[i].Score
}

type AlarmInfos []*model.AlarmInfo

func (a AlarmInfos) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a AlarmInfos) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

func (infos AlarmInfos) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return infos[j].AlarmScore < infos[i].AlarmScore
}

var alarmLogic = new(logic.AlarmInfoLogic)

func (service *AlarmService) HitAlert(request *AlarmRequest) *AlarmResponse {
	//验证参数
	response := judgeRequestParams(request)
	if response.Rtn != 0 {
		return response
	}
	//根据条件查询告警信息
	alarmInfos := findAlarmsBaseOnCondition(request)
	//排序
	if request.Order == "-1" {
		sort.Reverse(AlarmInfos(alarmInfos))
	} else {
		sort.Sort(AlarmInfos(alarmInfos))
	}
	//封装返回结果
	if alarmInfos != nil && len(alarmInfos) > 0 {
		err, results := fillAlarmResult(alarmInfos)
		if err != nil {
			response.Rtn = -1
			response.Message = err.Error()
			return response
		}
		response.Results = results
		response.Total = len(results)
	} else {
		response.Results = []AlarmResult{}
		response.Total = 0
	}
	response.Rtn = 0
	response.Message = "OK!"
	return response
}

//填充响应结果
func fillAlarmResult(alarmInfos []*model.AlarmInfo) (error, []AlarmResult) {
	alarmResults := make([]AlarmResult, 0)
	for _, alarmResult := range alarmInfos {
		faceImageId := alarmResult.AlarmTmplFeatrueId1
		hasF, feature := featureLogic.FindFaceFeatureByFaceImageId(faceImageId)
		//origUuid := alarmResult.AlarmOrigImageUuid
		//hasO, origImage := origImgeLogic.FindOrigImageByUUID(origUuid)
		if hasF {
			peopleId := feature.PeopleId
			imageId := feature.ImageId
			hasP, people := peopleLogic.FindPeopleById(peopleId)
			hasI, image := imageLogic.FindImageById(imageId)
			if hasP && hasI {
				compareResults := CompareResults{}
				rect := new(model.Rect)
				prop := new(model.Prop)
				json.Unmarshal([]byte(feature.FaceRect), rect)
				json.Unmarshal([]byte(feature.FaceProp), prop)
				peopleResult := PeopleResult{
					BronYear:       people.Birthday,
					FaceImageId:    people.PubId,
					FaceImageIdStr: people.PubId,
					FaceRect:       rect,
					FaceImageUri:   image.ImageUrl,
					PictureUri:     image.ImageUrl,
					Timestamp:      people.CreateTime.Unix(),
					Name:           people.Name,
					Nation:         people.Nation,
					PersonId:       people.PersonId,
					Gender:         people.Gender,
					Age:            prop.Age,
					Race:           prop.Race,
					SmileLevel:     prop.SmileLevel,
					BeautyLevel:    prop.BeautyLevel,
					RepositoryId:   people.RepositoryId,
					CustomField:    people.CustomField,
				}
				compareResults = packAndSortCompareResult(alarmResult, compareResults)
				regRect := new(model.Rect)
				regProp := new(model.Prop)
				json.Unmarshal([]byte(alarmResult.AlarmFaceRect), regRect)
				json.Unmarshal([]byte(alarmResult.AlarmFaceProp), regProp)
				recogResult := RecogResult{
					FaceImageId:    alarmResult.AlarmOrigImageUuid,
					FaceImageIdStr: alarmResult.AlarmOrigImageUuid,
					FaceRect:       regRect,
					FaceImageUri:   alarmResult.AlarmOrigImageUri,
					PictureUri:     alarmResult.AlarmOrigImageUri,
					Timestamp:      alarmResult.CreateTime.Unix(),
					Gender:         regProp.Gender,
					Age:            regProp.Age,
					Race:           regProp.Race,
					SmileLevel:     regProp.SmileLevel,
					BeautyLevel:    regProp.BeautyLevel,
					RepositoryId:   alarmResult.RepositoryId,
					CustomField:    "",
					TopN:           compareResults,
				}
				hitDetails := HitDetails{
					FaceImageId:     alarmResult.AlarmOrigImageUuid,
					HitFaceImageId:  people.PubId,
					HitRepositoryId: people.RepositoryId,
					HitSimilarity:   alarmResult.AlarmScore,
					HitType:         1,
					Id:              alarmResult.Id,
					Timestamp:       alarmResult.CreateTime.Unix(),
				}
				alarmResult := AlarmResult{
					RecogResult:  recogResult,
					PeopleResult: peopleResult,
					HitDetails:   hitDetails,
				}
				alarmResults = append(alarmResults, alarmResult)
			} else {
				log.Println("fillAlarmResult can not find people or image :", faceImageId)
			}
		} else {
			log.Println("fillAlarmResult can not find feature by faceImageId :", faceImageId)
			continue
		}
	}
	return nil, alarmResults
}

//打包并排序结果
func packAndSortCompareResult(alarmResult *model.AlarmInfo, compareResults CompareResults) CompareResults {
	if alarmResult.AlarmTmplFeatrueId1 != "" {
		compareResult := CompareResult{
			Score:     alarmResult.AlarmTmplScore1,
			FeatureId: alarmResult.AlarmTmplFeatrueId1,
		}
		compareResults = append(compareResults, &compareResult)
	}
	if alarmResult.AlarmTmplFeatrueId2 != "" {
		compareResult := CompareResult{
			Score:     alarmResult.AlarmTmplScore2,
			FeatureId: alarmResult.AlarmTmplFeatrueId2,
		}
		compareResults = append(compareResults, &compareResult)
	}
	if alarmResult.AlarmTmplFeatrueId3 != "" {
		compareResult := CompareResult{
			Score:     alarmResult.AlarmTmplScore3,
			FeatureId: alarmResult.AlarmTmplFeatrueId3,
		}
		compareResults = append(compareResults, &compareResult)
	}
	if alarmResult.AlarmScoreOthers != "" {
		results := CompareResults{}
		err := json.Unmarshal([]byte(alarmResult.AlarmScoreOthers), &results)
		if err != nil {
			log.Println("Unmarshal AlarmScoreOthers err :", err)
		} else {
			for _, v := range results {
				compareResults = append(compareResults, v)
			}
		}
	}
	sort.Sort(compareResults)
	return compareResults
}

//根据条件查询数据
func findAlarmsBaseOnCondition(request *AlarmRequest) []*model.AlarmInfo {
	condition := request.Hit_condition
	if len(request.Task_ids) > 0 {
		taskIds := request.Task_ids
		return alarmLogic.FindAlarmInfosByTaskIdsAndCondition(taskIds, condition)
	}
	if len(request.Camera_ids) > 0 {
		cameraIds := request.Camera_ids
		return alarmLogic.FindAlarmInfosByCameraIdssAndCondition(cameraIds, condition)
	}
	if len(request.Surveillances) > 0 {
		surveillances := request.Surveillances
		return alarmLogic.FindAlarmInfosBySurveillancessAndCondition(surveillances, condition)
	}
	return nil
}

//验证参数是否合格
func judgeRequestParams(request *AlarmRequest) *AlarmResponse {
	response := new(AlarmResponse)
	if len(request.Task_ids) == 0 && len(request.Camera_ids) == 0 && len(request.Surveillances) == 0 {
		log.Println("task_ids，camera_ids，surveillances不能全为空!")
		response.Rtn = -1
		response.Message = "task_ids，camera_ids，surveillances三个中需要填一个!"
		return response
	}
	if request.Hit_condition == nil {
		response.Rtn = -1
		response.Message = "需要填写过滤条件!"
		return response
	}
	if request.Order == "" {
		response.Rtn = -1
		response.Message = "需要填写结果排序方式!"
		return response
	}
	if request.Start == 0 {
		response.Rtn = -1
		response.Message = "start需要大于0!"
		return response
	}
	if request.Limit == 0 {
		response.Rtn = -1
		response.Message = "limit需要大于0!"
		return response
	}
	condition := request.Hit_condition
	timestamp := condition.Timestamp
	var startTime, endTime time.Time
	var err error
	if timestamp.Gte == "" {
		startTime, _ = time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00")
	} else {
		startTime, err = time.Parse("2006-01-02 15:04:05", timestamp.Gte)
		if err != nil {
			response.Rtn = -1
			response.Message = "timestamp gte不正确!"
			return response
		}
	}
	if timestamp.Lte == "" {
		endTime = time.Now()
	} else {
		endTime, err = time.Parse("2006-01-02 15:04:05", timestamp.Lte)
		if err != nil {
			response.Rtn = -1
			response.Message = "timestamp Lte不正确!"
			return response
		}
	}
	condition.StartTime = startTime
	condition.EndTime = endTime
	condition.Start = request.Start
	condition.Limit = request.Limit
	return response
}
