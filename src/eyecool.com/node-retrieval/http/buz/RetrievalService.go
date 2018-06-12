package buz

import (
	"eyecool.com/node-retrieval/utils"
	"eyecool.com/node-retrieval/logic"
	"eyecool.com/node-retrieval/model"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
	"encoding/base64"
	"eyecool.com/node-retrieval/http/service"
	"strings"
)

type RetrievalService struct {
}

type Retrieval struct {
	Face_image_id  string   //查询此 face image id 对应的人脸
	Repository_ids []string //查询库的ID列表，repository_ids和camera_ids，video_ids三选一
	Camera_ids     []string //查询的摄像头id, repository_ids和camera_ids,video_ids 三选一
	Video_ids      []string //查询的离线视频id, repository_ids和camera_ids, video_ids三选一
	Threshold      float64  //相似度分数线
	Using_ann      bool     //检索是否开启ann来加速  默认 true
	Topk           int32    //检索结果最大值 (topk 增加会导致检索速度变慢)  默认100
}
type Condition struct {
	Repository_id string
	Camera_id     string
	Video_id      string
	Name          string
	Timestamp     string
	Person_id     string //身份证搜索
}

type RetrievalRequest struct {
	Retrieval_query_id string
	Extra_fields       []string
	Async_query        bool //是否是异步检索 默认false
	Retrieval          *Retrieval
	Start              int //从第几个结果开始返回
	Limit              int //返回至多多少个结果
	Condition          *Condition
	Order              interface{}
}
type RetrievalResponse struct {
	Rtn              int                `json:"rtn"`               //错误码
	Message          string             `json:"message,omitempty"` //错误消息
	RetrievalQueryId string             `json:"retrieval_query_id,omitempty"`
	Total            int                `json:"total"`
	Results          []*RetrievalResult `json:"results,omitempty"`
}
type RetrievalResult struct {
	Annotation     int    `json:"annotation"`
	FaceImageId    string `json:"face_image_id,omitempty"` //过人id($编号@$集群号)
	FaceImageIdStr string `json:"face_image_id_str,omitempty"`
	FaceImageUri   string `json:"face_image_uri,omitempty"` //人脸图URI
	PictureUri     string `json:"picture_uri,omitempty"`    //场景图URI
	FaceRect       *Rect  `json:"face_rect,omitempty"`      //人脸框
	Timestamp      int64  `json:"timestamp,omitempty"`      //过人的时间戳
	BornYear       string `json:"born_year,omitempty"`
	Gender         int    `json:"gender,omitempty"`
	IsWritable     bool   `json:"is_writable,omitempty"`
	Name           string `json:"name,omitempty"`
	Nation         int `json:"nation"`
	PersonId       string `json:"person_id,omitempty"`
	CustomField    string `json:"custom_field,omitempty"`
	RepositoryId   string `json:"repository_id,omitempty"`
	Similarity     int    `json:"similarity,omitempty"`
	CameraId       string `json:"camera_id,omitempty"` //摄像头编号 ($编号@$集群号)
}
type Order struct {
	Similarity int
	Timestamp  int
}
type ResultInfo struct {
	Similarity  int    `json:"similarity,omitempty"`
	FaceImageId string `json:"faceImageId,omitempty"` // 地库人id
	Id          int64  `json:"id,omitempty"`          // 抓拍图片id
	Tmpl        Tmpl   `json:"tmpl"`
}
type Tmpl struct {
	FaceImageId string `json:"faceImageId,omitempty"`
	Pos         int    `json:"pos,omitempty"`
	Status      int    `json:"status,omitempty"`
	Id          int64  `json:"id,omitempty"`
}
type RetrievalParam struct {
	RepositoryParam string
	CameraParam     string
	VideoParam      string
	FaceImageId     string
	Feat            string
}

type RetrievalResults []*RetrievalResult

func (a RetrievalResults) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a RetrievalResults) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

func (infos RetrievalResults) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return infos[j].Timestamp < infos[i].Timestamp
}

type ResultInfos []*ResultInfo

func (a ResultInfos) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a ResultInfos) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

func (infos ResultInfos) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return infos[j].Similarity < infos[i].Similarity
}

var retrievalLogic = new(logic.RetrievalLogic)
var retrievalUitLogic = new(logic.RetrievalUnitLogic)

//人脸检索
func (this *RetrievalService) InsertAndRetrieval(request *RetrievalRequest) *RetrievalResponse {
	response := new(RetrievalResponse)
	if request.Retrieval_query_id == "" { //新查询
		response = insertAndRequestRetrieval(request)
	} else { //查询缓存
		response = queryCacheById(request)
	}
	return response
}

//条件查询
func (this *RetrievalService) ConditionQuery(request *RetrievalRequest) *RetrievalResponse {
	response := new(RetrievalResponse)
	if request.Condition == nil {
		response.Rtn = -1
		response.Message = "condition参数不正确!"
		return response
	}
	retrieval := new(model.Retrieval)
	condition := request.Condition
	if condition.Repository_id != "" {
		retrieval.RepositoryId = condition.Repository_id
	}
	if condition.Camera_id != "" {
		retrieval.CameraId = condition.Camera_id
	}
	if condition.Video_id != "" {
		retrieval.VideoId = condition.Video_id
	}
	if condition.Name != "" {
		retrieval.Name = condition.Name
	}
	if condition.Timestamp != "" {
		retrieval.Timestamp = condition.Timestamp
	}
	if condition.Person_id != "" {
		retrieval.PersonId = condition.Person_id
	}
	retrieval.StartIndex = request.Start
	retrieval.LimitResult = request.Limit
	peoples := peopleLogic.FindByRetrieval(retrieval)
	//填充条件检索结果
	retrievalResults := fillConditionRetrievalResult(peoples)
	if len(retrievalResults) > 0 {
		response.Total = len(retrievalResults)
		sort.Sort(retrievalResults)
		if request.Order != "" && request.Order != nil {
			orderByte, _ := json.Marshal(request.Order)
			order := new(Order)
			json.Unmarshal(orderByte, order)
			timestamp := order.Timestamp
			if timestamp == 1 {
				sort.Reverse(RetrievalResults(retrievalResults))
			}
		}
		retrievalResults = getConditionStartLimit(&retrievalResults, retrieval)
	}
	response.Message = "ok"
	response.Rtn = 0
	response.Results = retrievalResults
	return response
}

func getConditionStartLimit(results *RetrievalResults, retrieval *model.Retrieval) RetrievalResults {
	newResults := RetrievalResults{}
	size := len(*results)
	startIndex := retrieval.StartIndex
	limit := retrieval.LimitResult
	if size > 0 && startIndex >= 0 && limit > 0 && size > startIndex {
		num := startIndex + limit
		if size < num {
			num = size
		}
		for i := startIndex; i < num; i++ {
			newResults = append(newResults, (*results)[i])
		}
	}
	return newResults
}

//填充条件检索结果
func fillConditionRetrievalResult(peoples []*model.People) RetrievalResults {
	retrievalResults := make([]*RetrievalResult, 0)
	if len(peoples) > 0 {
		for _, people := range peoples {
			fmt.Println(people)
			result := &RetrievalResult{
				CustomField:    people.CustomField,
				FaceImageId:    strconv.Itoa(int(people.Id)) + "@" + strconv.Itoa(people.ClusterId),
				FaceImageIdStr: strconv.Itoa(int(people.Id)) + "@" + strconv.Itoa(people.ClusterId),
				BornYear:       people.Birthday,
				Gender:         people.Gender,
				Name:           people.Name,
				Nation:         people.Nation,
				PersonId:       people.PersonId,
				RepositoryId:   people.RepositoryId,
				Timestamp:      time.Now().Unix(),
			}
			hasF, feature := featureLogic.FindFaceFeatureByPeopleId(people.Id) //仅针对一个图片一张人脸的情况
			if hasF {
				result.FaceRect = &Rect{
					X: feature.X,
					Y: feature.Y,
					H: feature.H,
					W: feature.W,
				}
				hasI, image := imageLogic.FindImageById(feature.ImageId)
				if hasI {
					result.PictureUri = image.ImageUrl + "@" + strconv.Itoa(image.ClusterId)
					result.FaceImageUri = image.ImageUrl + "@" + strconv.Itoa(image.ClusterId)
				}
			}
			retrievalResults = append(retrievalResults, result)
		}
	}
	return retrievalResults
}

func insertAndRequestRetrieval(request *RetrievalRequest) *RetrievalResponse {
	response := new(RetrievalResponse)
	retrievalParam := request.Retrieval
	if retrievalParam == nil {
		response.Rtn = -1
		response.Message = "参数不正确!"
		return response
	}
	if retrievalParam.Face_image_id == "" {
		response.Rtn = -1
		response.Message = "Face_image_id不正确!"
		return response
	}
	faceId, _, err := utils.GetIdAndClusterId(retrievalParam.Face_image_id)
	if err != nil || faceId == -2 {
		response.Rtn = -1
		response.Message = "Face_image_id--->faceId不正确!"
		return response
	}
	hasF, feature := featureLogic.FindFaceFeatureByPkId(faceId)
	if !hasF {
		response.Rtn = -1
		response.Message = "can not find faceFeature by id :" + strconv.Itoa(faceId)
		return response
	}
	response = insertRetrieval(feature, request)
	return response
}
func insertRetrieval(feature *model.FaceFeature, request *RetrievalRequest) *RetrievalResponse {
	response := new(RetrievalResponse)
	retrieval := new(model.Retrieval)
	retrievalParam := new(RetrievalParam)
	retrieval.PeopleId = feature.PeopleId
	repositoryIds := request.Retrieval.Repository_ids
	cameraIds := request.Retrieval.Camera_ids
	videoIds := request.Retrieval.Video_ids
	if len(repositoryIds) > 0 {
		repositoryParam := utils.SplitArrayByComma(repositoryIds)
		retrievalParam.RepositoryParam = repositoryParam
		repByte, _ := json.Marshal(repositoryIds)
		retrieval.RepositoryIds = string(repByte)
	}
	if len(cameraIds) > 0 {
		cameraParam := utils.SplitArrayByComma(cameraIds)
		retrievalParam.CameraParam = cameraParam
		camByte, _ := json.Marshal(cameraIds)
		retrieval.CameraIds = string(camByte)
	}
	if len(videoIds) > 0 {
		videoParam := utils.SplitArrayByComma(videoIds)
		retrievalParam.VideoParam = videoParam
		vidByte, _ := json.Marshal(videoIds)
		retrieval.VideoIds = string(vidByte)
	}
	retrieval.UsingAnn = request.Retrieval.Using_ann
	if request.Retrieval.Topk == 0 {
		request.Retrieval.Topk = 100
	}
	retrieval.Topk = request.Retrieval.Topk
	retrieval.Threshold = request.Retrieval.Threshold
	similarity := -1
	//排序
	if request.Order != "" && request.Order != nil {
		orderByte, _ := json.Marshal(request.Order)
		order := new(Order)
		json.Unmarshal(orderByte, order)
		similarity = order.Similarity
		retrieval.OrderJson = string(orderByte)
	}
	retrieval.StartIndex = request.Start
	retrieval.LimitResult = request.Limit
	retrieval.AsyncQuery = request.Async_query
	if len(request.Extra_fields) > 0 {
		extByte, _ := json.Marshal(request.Extra_fields)
		retrieval.ExtraFields = string(extByte)
	}
	retrieval.ClusterId = 1
	retrieval.CreatorId = 1
	retrieval.CreateTime = time.Now()
	retrieval.UpdateTime = time.Now()
	retrievalLogic.Insert(retrieval)
	retrieval.RetrievalQueryId = strconv.Itoa(int(retrieval.Id)) + "@" + strconv.Itoa(retrieval.ClusterId)
	retrievalLogic.UpdateRetrievalById(retrieval)

	retrievalParam.Feat = feature.Feat
	retrievalParam.FaceImageId = request.Retrieval.Face_image_id
	//判断同步检索还是异步检索
	if !request.Async_query {
		//调用Go接口进行检索
		msg := sendToGoRetrieval(retrievalParam, retrieval)
		if "success" == msg.Msg {
			time.Sleep(1*time.Second)
			_, retrievalExsit := retrievalLogic.SelectRetrievalById(retrieval.Id)
			fmt.Println(retrievalExsit)
			response.Rtn = 0
			response.Message = "检索成功"
			response.RetrievalQueryId = retrievalExsit.RetrievalQueryId
			if retrievalParam.RepositoryParam != "" {
				results := pushResults(retrievalExsit, similarity)
				response.Total = retrievalExsit.Total
				response.Results = results
			}
			if retrievalParam.CameraParam != "" {
				results := pushCameraSnapshotResults(retrievalExsit, similarity)
				_, retrievalExsit = retrievalLogic.SelectRetrievalById(retrieval.Id)  //再次获取total记录数
				response.Total = retrievalExsit.Total
				response.Results = results
			}
			return response
		}
		response.Rtn = -1
		response.Message = "GOGOGO failed!"
		return response
	}

	sendToGoRetrieval(retrievalParam, retrieval)
	response.Rtn = 0
	response.Message = "操作成功"
	response.RetrievalQueryId = request.Retrieval_query_id
	return response
}

//调用接口进行检索,分repository和camera
func sendToGoRetrieval(retrievalParam *RetrievalParam, retrieval *model.Retrieval) *model.RetrievalResponse {
	response := &model.RetrievalResponse{}
	featByte, err := base64.StdEncoding.DecodeString(retrievalParam.Feat)
	fmt.Println("sendToGoRetrieval feat:", retrievalParam.Feat)
	retrievalRequest := model.RetrievalRequest{
		Uuid:          "",
		RetrievalId:   int32(retrieval.Id),
		FaceImageId:   retrievalParam.FaceImageId,
		CameraIds:     retrievalParam.CameraParam,
		RepositoryIds: retrievalParam.RepositoryParam,
		VideoIds:      retrievalParam.VideoParam,
		Threshold:     retrieval.Threshold,
		Topk:          retrieval.Topk,
		Async:         retrieval.AsyncQuery,
		Params:        "",
		Feats:         featByte,
	}

	if retrievalParam.RepositoryParam != "" {
		err = service.RetrievalRepositoryTarget(nil, &retrievalRequest, response)
		if err != nil {
			response.Msg = "error"
		}
	}
	if retrievalParam.CameraParam != "" {
		cameraIds := strings.Split(retrievalParam.CameraParam, ",")
		for i := range cameraIds {
			resp := &model.RetrievalResponse{}
			retrievalRequest.CameraIds = cameraIds[i]
			err = service.RetrievalCameraTarget(nil, &retrievalRequest, resp)
			if err != nil {
				continue
			}
			response.Total += resp.Total
			response.RetrievalId = resp.RetrievalId
			response.Msg = resp.Msg
		}
	}
	return response
}
func queryCacheById(request *RetrievalRequest) *RetrievalResponse {
	results := RetrievalResults{}
	response := new(RetrievalResponse)
	retrievalId, _, err := utils.GetIdAndClusterId(request.Retrieval_query_id)
	if err != nil {
		response.Rtn = -1
		response.Message = "Retrieval_query_id不正确!"
		return response
	}
	has, retrieval := retrievalLogic.SelectRetrievalById(int64(retrievalId))
	if !has {
		response.Rtn = -1
		response.Message = "未查询到此记录：" + request.Retrieval_query_id
		return response
	}
	if has {
		if retrieval.Results == "" {
			response.Rtn = -1
			response.Message = "查询到此记录：" + request.Retrieval_query_id + "未检索出结果"
			return response
		}
		response.Rtn = 0
		response.Message = "检索成功"
		response.RetrievalQueryId = request.Retrieval_query_id
		similarity := -1
		//排序
		if request.Order != "" && request.Order != nil {
			orderByte, _ := json.Marshal(request.Order)
			order := new(Order)
			json.Unmarshal(orderByte, order)
			similarity = order.Similarity
			retrieval.OrderJson = string(orderByte)
		}
		//地库搜索
		if retrieval.RepositoryIds != "" {
			retrieval.StartIndex = request.Start
			retrieval.LimitResult = request.Limit
			results = pushResults(retrieval, similarity)
			response.Results = results
			response.Total = retrieval.Total
		}
		//摄像机抓拍搜索
		if retrieval.CameraIds != "" {
			results = pushCameraSnapshotResults(retrieval, similarity)
			response.Results = results
			response.Total = retrieval.Total
		}
	}
	return response
}
func pushCameraSnapshotResults(retrieval *model.Retrieval, similarity int) []*RetrievalResult {
	retrievalResults := RetrievalResults{}
	results := retrieval.Results
	resultInfos := ResultInfos{}
	if retrieval.Total == 0 { //没有组装过
		resultInfos = getCameraRetrievalResultAndSort(retrieval.Id);
		fmt.Println(" pushCameraSnapshotResults  -----resultInfos: ", resultInfos)
		if len(resultInfos) > 0 {
			retrieval.Total = len(resultInfos)
			resultInfoByte, _ := json.Marshal(resultInfos)
			retrieval.Results = string(resultInfoByte)
			retrievalLogic.UpdateRetrievalById(retrieval)
		}
	} else { //组装过了
		err := json.Unmarshal([]byte(results), &resultInfos)
		if err != nil {
			fmt.Println(err)
		}
	}
	// 排序规则
	if similarity == 1 { // 排序规则
		sort.Reverse(resultInfos)
	}
	resultInfos = getStartLimit(&resultInfos, retrieval)
	//填充摄像机搜索结果
	retrievalResults = fillCameraRetrievalResult(&resultInfos, similarity)
	return retrievalResults
}

func getCameraRetrievalResultAndSort(retrievalId int64) ResultInfos {
	newInfos := ResultInfos{}
	retrievalUnits := retrievalUitLogic.FindByRetrievalId(retrievalId)
	if len(retrievalUnits) > 0 {
		for _, unit := range retrievalUnits {
			results := unit.Results
			//获取resultMap
			resultInfos := &ResultInfos{}
			json.Unmarshal([]byte(results), resultInfos)
			//fmt.Println("retrievalUitLogic results",results,"resultInfos", resultInfos)
			if len(*resultInfos) > 0 {
				for _, resultInfo := range *resultInfos {
					resultInfo.Id = resultInfo.Tmpl.Id
					fmt.Println("resultInfo", resultInfo)
					newInfos = append(newInfos, resultInfo)
				}
			}
		}
	}
	//fmt.Println("getCameraRetrievalResultAndSort newInfos newInfos newInfos", newInfos)
	sort.Sort(newInfos)
	return newInfos
}

func fillCameraRetrievalResult(resultInfos *ResultInfos, similarity int) RetrievalResults {
	retrievalResults := RetrievalResults{}
	if len(*resultInfos) > 0 {
		for _, info := range *resultInfos {
			hasO, origImage := origImgeLogic.FindOrigImageById(info.Id)
			if hasO {
				_, clusterId, _ := utils.GetIdAndClusterId(origImage.CameraId)
				rectStr := origImage.FaceRect
				rect := &Rect{}
				json.Unmarshal([]byte(rectStr), rect)
				retrievalResult := &RetrievalResult{
					CameraId:       origImage.CameraId,
					FaceImageId:    strconv.Itoa(origImage.Id) + "@" + strconv.Itoa(clusterId),
					FaceImageIdStr: strconv.Itoa(origImage.Id) + "@" + strconv.Itoa(clusterId),
					FaceImageUri:   origImage.ImageContextPath + "/" + origImage.FaceImageUri + "@" + strconv.Itoa(clusterId),
					FaceRect:       rect,
					PictureUri:     origImage.ImageContextPath + "/" + origImage.FaceImageUri + "@" + strconv.Itoa(clusterId),
					Timestamp:      time.Now().Unix(),
					Similarity:     similarity,
				}
				retrievalResults = append(retrievalResults, retrievalResult)
			}
		}
	}
	return retrievalResults
}

//根据repository检索
func pushResults(retrieval *model.Retrieval, similarity int) []*RetrievalResult {
	retrievalResults := RetrievalResults{}
	results := retrieval.Results
	resultInfos := &ResultInfos{}
	if retrieval.Total == 0 { //没有组装过
		repositoryIds := retrieval.RepositoryIds
		cameraIds := retrieval.CameraIds
		keys := make([]string, 0)
		if repositoryIds != "" {
			json.Unmarshal([]byte(repositoryIds), &keys)
		}
		if cameraIds != "" {
			json.Unmarshal([]byte(cameraIds), &keys)
		}
		fmt.Println("pushResults ------------------- keys:", keys, "repositoryIds:", repositoryIds,"results:",results)
		resultInfos = getResultAndSort(results, keys);
		if len(*resultInfos) > 0 {
			retrieval.Total = len(*resultInfos)
			resultInfoByte, _ := json.Marshal(resultInfos)
			retrieval.Results = string(resultInfoByte)
			retrievalLogic.UpdateRetrievalById(retrieval)
		}
	} else { //组装过了
		json.Unmarshal([]byte(results), &resultInfos)
	}
	//排序规则
	if similarity == 1 {
		//反转下resultInfos
		sort.Reverse(resultInfos)
	}

	//分页返回值
	*resultInfos = getStartLimit(resultInfos, retrieval)
	fmt.Println("pushResults resultInfos:",*resultInfos,"长度是:",len(*resultInfos))

	//填充结果
	retrievalResults = fillRetrievalResults(resultInfos, similarity)

	fmt.Println("pushResults retrievalResults:",retrievalResults)

	return retrievalResults
}

//填充retrievalResults结果
func fillRetrievalResults(resultInfos *ResultInfos, similarity int) RetrievalResults {
	retrievalResults := RetrievalResults{}
	if len(*resultInfos) > 0 {
		for _, result := range *resultInfos {
			fmt.Println("fillRetrievalResults result :",result)
			faceImageIdStr := result.Tmpl.FaceImageId
			faceImageId, _, err := utils.GetIdAndClusterId(faceImageIdStr)
			if err != nil || faceImageId == -2 {
				fmt.Println(err)
				continue
			}
			hasF, feature := featureLogic.FindFaceFeatureByPkId(faceImageId)
			if hasF {
				hasI, image := imageLogic.FindImageById(feature.ImageId)
				hasP, people := peopleLogic.FindPeopleById(feature.PeopleId)
				if hasI && hasP {
					retrievalResult := &RetrievalResult{
						Annotation:     0,
						FaceImageId:    feature.FaceImageId,
						FaceImageIdStr: feature.FaceImageId,
						BornYear:       people.Birthday,
						FaceImageUri:   image.ImageUrl + "@" + strconv.Itoa(image.ClusterId),
						FaceRect: &Rect{
							X: feature.X,
							Y: feature.Y,
							H: feature.H,
							W: feature.W,
						},
						Gender:       people.Gender,
						IsWritable:   false,
						Name:         people.Name,
						Nation:       people.Nation,
						PictureUri:   image.ImageUrl + "@" + strconv.Itoa(image.ClusterId),
						RepositoryId: people.RepositoryId,
						Timestamp:    time.Now().Unix(),
						PersonId:     people.PersonId,
						Similarity:   similarity,
					}
					retrievalResults = append(retrievalResults, retrievalResult)
				}
			}
		}
	}
	return retrievalResults;
}

//分页返回值  startIndex:开始位置  limitResult:获取长度
func getStartLimit(resultInfos *ResultInfos, retrieval *model.Retrieval) ResultInfos {
	newInfos := ResultInfos{}
	size := len(*resultInfos)
	startIndex := retrieval.StartIndex
	limit := retrieval.LimitResult
	if size > 0 && startIndex >= 0 && limit > 0 && size > startIndex {
		num := startIndex + limit
		if size < num {
			num = size
		}
		for i := startIndex; i < num; i++ {
			newInfos = append(newInfos, (*resultInfos)[i])
		}
	}
	return newInfos
}

//获取结果并降序排序
func getResultAndSort(results string, keys []string) *ResultInfos {
	resultInfos := ResultInfos{}
	//获取resultMap
	resultMap := map[string]ResultInfos{}
	json.Unmarshal([]byte(results), &resultMap)
	fmt.Println("getResultAndSort results :", results)
	if len(keys) > 0 {
		for _, key := range keys {
			if v, f := resultMap[key]; f {
				for _, info := range v {
					resultInfo := &ResultInfo{
						Similarity:  info.Similarity,
						Tmpl:        info.Tmpl,
						FaceImageId: info.FaceImageId,
					}
					resultInfos = append(resultInfos, resultInfo)
				}
			}
		}
	}
	fmt.Println("getResultAndSort resultInfos : ", resultInfos)
	if len(resultInfos) > 0 {
		sort.Sort(&resultInfos)
	}
	return &resultInfos
}
