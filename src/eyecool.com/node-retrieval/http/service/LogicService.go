package service

import (
	"eyecool.com/node-retrieval/model"
	"context"
	"log"
	"strings"
	"eyecool.com/node-retrieval/utils"
	"eyecool.com/node-retrieval/seek"
	"strconv"
	"encoding/json"
	"eyecool.com/node-retrieval/global"
	"github.com/polaris1119/set"
	"fmt"
	"eyecool.com/node-retrieval/algorithm"
	"encoding/base64"
	"time"
	"eyecool.com/node-retrieval/logic"
)

func RetrievalRepositoryTarget(ctx context.Context, req *model.RetrievalRequest, resp *model.RetrievalResponse) error {
	log.Println("Received Dispatch.RetrievalRepositoryTarget request 2222222222222222 ", req.Uuid)
	repositoryIds := strings.Split(req.RepositoryIds, ",")
	resultList := make(map[string]utils.FeaturePairList, len(repositoryIds))
	for i := 0; i < len(repositoryIds) && len(repositoryIds[i]) > 0; i++ {
		list, err := seek.CompareTarget(repositoryIds[i], req.Feats, int(req.Topk))
		if err != nil {
			continue
		}
		resultList[repositoryIds[i]] = list
	}
	resultBytes, _ := json.Marshal(resultList)
	logic.DefaultRetrieval.UpdateRetrievalResults(int(req.RetrievalId), string(resultBytes))
	resp.RetrievalId = strconv.Itoa(int(req.RetrievalId))
	resp.Msg = "success"
	return nil
}

func RetrievalCameraTarget(ctx context.Context, req *model.RetrievalRequest, resp *model.RetrievalResponse) error {
	log.Println("Received Dispatch.RetrievalCameraTarget request 2222222222222222 ", req.Uuid)
	if req.Async {
		go compareOrigImage(req)
	} else {
		err := compareOrigImage(req)
		if err != nil {
			resp.Msg = global.FAILED
		} else {
			resp.Msg = global.SUCCESS
		}
	}
	resp.RetrievalId = strconv.Itoa(int(req.RetrievalId))
	return nil
}

func RetrievalRepositoryFeatureInsert(ctx context.Context, req *model.RetrievalFeatureRequest, resp *model.RetrievalFeatureResponse) error {
	log.Println("Received Dispatch.RetrievalRepositoryFeatrueInsert request 2222222222222222 ", req.Uuid)
	if req.Type == global.DELETE {
		if setcoll, ok := global.G_ReposHandleMap[global.BuildReposKey(req.RepositoryId)]; ok {
			setcoll.Each(func(entry interface{}) bool {
				if e := entry.(*model.FeatureEntry); e != nil && e.FaceImageId == req.Id {
					e.Status = 1
					return false
				}
				return true
			})
		}
		return nil
	} else {
		total, pos := global.G_ChlFaceX.ChlFaceSdkListInsert(global.BuildReposKey(req.RepositoryId), -1, req.Feats, int(req.FeatNum))
		fmt.Println("total: ", total, "  pos: ", pos)
		if pos >= 0 {
			faceImageIds := strings.Split(req.Id, ",")
			peopleIds := strings.Split(req.PeopleId, ",")
			if _, ok := global.G_ReposHandleMap[global.BuildReposKey(req.RepositoryId)]; !ok {
				global.G_ReposHandleMap[global.BuildReposKey(req.RepositoryId)] = set.New(set.NonThreadSafe)
			}
			setcoll := global.G_ReposHandleMap[global.BuildReposKey(req.RepositoryId)]
			for i := 0; i < int(req.FeatNum); i++ {
				pid, _ := strconv.Atoi(peopleIds[i])
				entry := &model.FeatureEntry{PeopleId: int64(pid), FaceImageId: faceImageIds[i], Pos: pos + i, Status: 0}
				setcoll.Add(entry)
			}
		} else {
		}
		if resp != nil {
			resp.Pos = int32(pos)
			resp.Total = int32(total)
		}
	}
	return nil
}

func InsertOrigImage(ctx context.Context, req *model.OrigImageRequest) error {
	log.Println("Received Dispatch.InsertOrigImage request 2222222222222222 ", req)
	image := &model.OrigImageFull{
		Uuid:             req.Uuid,
		CameraId:         req.CameraId,
		ClusterId:        1,
		ImageName:        req.ImageName,
		ImageRealPath:    req.ImageRealPath,
		FaceNum:          int(req.FaceNum),
		FeatList:         req.FeatList,
		FaceRect:         req.FaceRect,
		FaceProp:         req.FaceProp,
		Timestamp:        req.Timestamp,
		ImageContextPath: req.ImageContextPath,
		FaceImageUri:     req.FaceImageUri,
		UpdateTime:       time.Now(),
	}
	err := logic.DefaultOrigImageFull.Insert(image)
	if err != nil {
		//rsp.Msg = global.FAILED
	} else {
		//	rsp.Msg = global.SUCCESS
	}
	return nil
}

func RepositoryLifecycle(ctx context.Context, req *model.LifecycleRequest, rsp *model.LifecycleResponse) error {
	setcoll, ok := global.G_ReposHandleMap[global.BuildReposKey(req.RepositoryId)]
	switch req.Type {
	case global.INSERT:
		if !ok {
			ret := global.CreateReposFaceSdkCache(global.BuildReposKey(req.RepositoryId))
			if !ret {
				rsp.Msg = global.FAILED
				return nil
			}
			fmt.Println("INSERT ret: ", ret, "  map: ", global.G_ReposHandleMap)
		}

	case global.DELETE:
		if ok {
			//清除内存列表
			setcoll.Clear()
			delete(global.G_ReposHandleMap, global.BuildReposKey(req.RepositoryId))
			//清除算法缓存
			global.G_ChlFaceX.ChlFaceSdkListDestroy(global.BuildReposKey(req.RepositoryId))
			fmt.Println("DELETE : ", "  map: ", global.G_ReposHandleMap)
		}
	}
	rsp.Msg = global.SUCCESS
	return nil
}

func compareOrigImage(req *model.RetrievalRequest) error {
	origImageTopNList := utils.NewOrigImageTopNList(int(req.Topk))
	origImages := logic.DefaultOrigImage.FindAllByCameraId(req.CameraIds)
	//插入一条记录保存子查询任务
	unit := &model.RetrievalUnit{
		RetrievalId: int(req.RetrievalId),
		CameraId:    req.CameraIds,
		Type:        1,
		Total:       len(origImages),
		CreateTime:  time.Now(),
	}
	logic.DefaultRetrievalUnit.Insert(unit)
	batch, incr := 3000, 0
	for i, v := range origImages {
		if v.FaceNum == 1 {
			compFeats, err := base64.StdEncoding.DecodeString(v.FeatList)
			if err != nil {
				continue
			}
			score := global.G_ChlFaceX.ChlFaceSdkFeatureCompare(algorithm.DefaultChannelNo, req.Feats, compFeats)
			v.FeatList = ""
			origImageTopNList.Put(&utils.OrigImagePair{
				Similarity: score,
				Tmpl:       v,
			})
		} else if v.FaceNum > 1 {
			featLists := strings.Split(v.FeatList, ",")
			for i, _ := range featLists {
				compFeats, err := base64.StdEncoding.DecodeString(featLists[i])
				if err != nil {
					continue
				}
				score := global.G_ChlFaceX.ChlFaceSdkFeatureCompare(algorithm.DefaultChannelNo, req.Feats, compFeats)
				v.FeatList = ""
				origImageTopNList.Put(&utils.OrigImagePair{
					Similarity: score,
					Tmpl:       v,
				})
			}
		}
		incr = i
		if incr%batch == 0 {
			resultBytes, _ := json.Marshal(origImageTopNList.TopNList)
			logic.DefaultRetrievalUnit.UpdateRetrievalUnitResults(unit.Id, incr, string(resultBytes))
		}
	}
	if incr%batch != 0 {
		resultBytes, _ := json.Marshal(origImageTopNList.TopNList)
		logic.DefaultRetrievalUnit.UpdateRetrievalUnitResults(unit.Id, incr, string(resultBytes))
	}
	return nil
}
