package controller

import (
	"encoding/base64"
	"encoding/json"
	"eyecool.com/node-retrieval/model"
	"github.com/emicklei/go-restful"
	"eyecool.com/node-retrieval/http/service"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"eyecool.com/node-retrieval/global"
	"eyecool.com/node-retrieval/http/buz"
)

type RetrievalController struct {
}

type RetrievalQuery struct {
	RetrievalQueryId string `json:"retrieval_query_id"`
	ExtraFields      string `json:"extra_fields"`
	AsyncQuery       bool   `json:"async_query"`
	Retrieval        RetrievalUnit
	Order            string
	Start            int
	Limit            int
}

type RetrievalUnit struct {
	FaceImageId   string   `json:"face_image_id"`
	RepositoryIds []string `json:"repository_ids"`
	CameraIds     []string `json:"camera_ids"`
	VideoIds      []string `json:"video_ids"`
	Threshold     float64  `json:"threshold"`
	UsingAnn      bool     `json:"using_ann"`
	Topk          int      `json:"topk"`
}

func (s *RetrievalController) Retrieval(req *restful.Request, rsp *restful.Response) {
	log.Print("Received AccessController.Retrieval API request")
	rsp.WriteEntity(map[string]string{
		"message": "Hi, this is the Verify API",
	})
	retrievalQuery := RetrievalQuery{}
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &retrievalQuery)
	if err != nil {
		fmt.Println("Unmarshal err : ", err)
	}
	fmt.Println("获取json中的retrievalQuery:", retrievalQuery)
	rsp.ResponseWriter.Write([]byte("xxxxxxxxResponseWriter Verify xxxxxxxx"))
}

func (s *RetrievalController) ViewRepos(req *restful.Request, rsp *restful.Response) {
	log.Print("Received ViewRepos API request")
	outstr := ""
	for key, set := range global.G_ReposHandleMap {
		outstr += fmt.Sprintf(" Repos key : %s len : %d  \r\n  ", key, set.Size())
	}
	for key, size := range global.G_ChlFaceX.HandlesCachedSize {
		outstr += fmt.Sprintf(" ChlFace cached key : %s  cachesize : %d  \r\n  ", key, size)
	}
	if outstr == "" {
		outstr = " cache repos is empty !!!"
	}
	rsp.ResponseWriter.Write([]byte(outstr))
}

type FeatureRequest struct {
	Uuid         string
	PeopleId     string
	CameraId     string
	RepositoryId string
	Id           string
	Type         int32
	FeatNum      int32
	FeatsBase64  string
}

type RetrievalRequest struct {
	Uuid          string
	RetrievalId   int32
	FaceImageId   string
	CameraIds     string
	RepositoryIds string
	VideoIds      string
	Threshold     float64
	Topk          int32
	Async         bool
	Params        string
	FeatsBase64   string
}

/**
	通过上传一张临时图片，检索若干个相机抓拍库，返回topn 的目标人员结果（查询结果先保存在RetrievalUnit表）
 */
func (c *RetrievalController) RetrievalCameraTarget(req *restful.Request, rsp *restful.Response) {
	rr := RetrievalRequest{}
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &rr)
	if err != nil {
		fmt.Println("Unmarshal FeatureRequest err : ", err)
	}
	featBytes, err := base64.StdEncoding.DecodeString(rr.FeatsBase64)
	retrievalRequest := model.RetrievalRequest{
		Uuid:          rr.Uuid,
		RetrievalId:   rr.RetrievalId,
		FaceImageId:   rr.FaceImageId,
		CameraIds:     rr.CameraIds,
		RepositoryIds: rr.RepositoryIds,
		VideoIds:      rr.VideoIds,
		Threshold:     rr.Threshold,
		Topk:          rr.Topk,
		Async:         rr.Async,
		Params:        rr.Params,
		Feats:         featBytes,
	}
	response := &model.RetrievalResponse{}
	cameraIds := strings.Split(rr.CameraIds, ",")
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
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

/**
	通过上传一张临时图片，检索若干个目标底库，返回topn 的目标人员结果
 */
func (c *RetrievalController) RetrievalRepositoryTarget(req *restful.Request, rsp *restful.Response) {
	rr := RetrievalRequest{}
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &rr)
	if err != nil {
		fmt.Println("Unmarshal FeatureRequest err : ", err)
	}
	featBytes, err := base64.StdEncoding.DecodeString(rr.FeatsBase64)

	retrievalRequest := model.RetrievalRequest{
		Uuid:          rr.Uuid,
		RetrievalId:   rr.RetrievalId,
		FaceImageId:   rr.FaceImageId,
		CameraIds:     rr.CameraIds,
		RepositoryIds: rr.RepositoryIds,
		VideoIds:      rr.VideoIds,
		Threshold:     rr.Threshold,
		Topk:          rr.Topk,
		Async:         rr.Async,
		Params:        rr.Params,
		Feats:         featBytes,
	}
	response := &model.RetrievalResponse{}
	err = service.RetrievalRepositoryTarget(nil, &retrievalRequest, response)
	if err != nil {
		fmt.Println(" RetrievalRepositoryTarget err : ", err)
	}
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

/**
	目标底库里添加一条特征记录
 */
func (s *RetrievalController) RepositoryFeatureInsert(req *restful.Request, rsp *restful.Response) {
	log.Print("Received AccessController.RepositoryFeatureInsert API request")
	fr := FeatureRequest{}
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &fr)
	if err != nil {
		fmt.Println("Unmarshal FeatureRequest err : ", err)
	}
	featbytes, err := base64.StdEncoding.DecodeString(fr.FeatsBase64)
	retrievalFeatureRequest := &model.RetrievalFeatureRequest{
		RepositoryId: fr.RepositoryId,
		Id:           fr.Id,
		PeopleId:     fr.PeopleId,
		Type:         fr.Type,
		FeatNum:      fr.FeatNum,
		Feats:        featbytes,
	}
	resp := &model.RetrievalFeatureResponse{}
	err = service.RetrievalRepositoryFeatureInsert(nil, retrievalFeatureRequest, resp)
	responseBytes, _ := json.Marshal(resp)
	rsp.ResponseWriter.Write(responseBytes)
}

/**
	操作目标底库的生命周期，如增加，或删除 底库
 */
func (s *RetrievalController) RepositoryLifecycle(req *restful.Request, rsp *restful.Response) {
	log.Print("Received AccessController.RepositoryLifecycle API request")
	lifecycleRequest := model.LifecycleRequest{}
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &lifecycleRequest)
	if err != nil {
		fmt.Println("Unmarshal LifecycleRequest err : ", err)
	}
	response := &model.LifecycleResponse{}
	err = service.RepositoryLifecycle(nil, &lifecycleRequest, response)
	fmt.Println("Req LifecycleRequest : ", lifecycleRequest, "  err : ", err)
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

var retrievalService = new(buz.RetrievalService)
//人脸检索接口  /retrieval  POST
func (this *RetrievalController) PictureSynchronized(req *restful.Request, rsp *restful.Response) {
	response := new(buz.RetrievalResponse)
	request := new(buz.RetrievalRequest)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, request)
	fmt.Println("PictureSynchronized controller request :", request)
	if err != nil {
		fmt.Println("PictureSynchronized Unmarshal  err : ", err, ":", request)
		response.Rtn = -1
		response.Message = err.Error()

		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		return
	}
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response = retrievalService.InsertAndRetrieval(request)
	} else {
		response.Rtn = -1
		response.Message = "用户未登录!"
	}

	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

// 条件查询接口  /condition/query  POST
func (this *RetrievalController) ConditionQuery(req *restful.Request, rsp *restful.Response) {
	response := new(buz.RetrievalResponse)
	request := new(buz.RetrievalRequest)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, request)
	if err != nil {
		fmt.Println("ConditionQuery Unmarshal  err : ", err, ":", request)
		response.Rtn = -1
		response.Message = err.Error()
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		return
	}
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response = retrievalService.ConditionQuery(request)
	} else {
		response.Rtn = -1
		response.Message = "用户未登录!"
	}
	responseBytes, _ := json.Marshal(response)
	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min

	rsp.ResponseWriter.Write(responseBytes)
}
