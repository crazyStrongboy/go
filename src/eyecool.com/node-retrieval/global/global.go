package global

import (
	"eyecool.com/node-retrieval/model"
	"eyecool.com/node-retrieval/algorithm"
	"fmt"
	"github.com/labstack/gommon/log"
	"bytes"
	"encoding/base64"
	"github.com/polaris1119/set"
	"eyecool.com/node-retrieval/logic"
)

const (
	FEATURE_LENGTH              = 2600
	REPOS_MAX_FEATURE_SIZE  int = 20000
	CAMERA_MAX_FEATURE_SIZE int = 20000
	REPOS_PREFIX                = "repos_"
	CAMERA_PREFIX               = "camera_"
)
const (
	INSERT  = 0
	DELETE  = 1
	SUCCESS = "success"
	FAILED  = "failed"
)
const (
	Type_Image   = 0
	Type_Feature = 1
)

var G_Push_Pattern = false

var G_TaskTupleMap map[string][]*model.TaskTuple

var G_ChlFaceX *algorithm.ChlFaceX

var G_ReposHandleMap map[string]set.Interface = make(map[string]set.Interface)

var G_CameraHandleMap map[string]set.Interface = make(map[string]set.Interface)

func init() {
	fmt.Println("xxx global start do init ........")
	G_TaskTupleMap = logic.DefaultTask.FindTaskTupleMap()
	//初始化算法
	G_ChlFaceX = algorithm.NewChlFaceX()
	ret := G_ChlFaceX.ChlFaceSdkInit()
	if ret == -1 {
		log.Fatalf(" ChlFaceSdkInit[%d] faild ! ", ret)
	}
	//加载库特征
	loadRepositoryFeature()

}

func loadRepositoryFeature() {
	fmt.Println("xxx global start do loadRepositoryFeature ........")
	repos := logic.DefaultRepository.FindAll()
	for _, v := range repos {
		CreateReposFaceSdkCache(BuildReposKey(v.Id))
	}

	buff := new(bytes.Buffer)
	batch := 1;
	for _, r := range repos {
		faceFeatures := logic.DefaultFaceFeature.FindFaceFeaturesByRepositoryId(r.Id)
		incr, len := 0, 0
		imageIds := make([]string, 0, batch)
		peopleIds := make([]int64, 0, batch)
		for i, v := range faceFeatures {
			feats, err := base64.StdEncoding.DecodeString(faceFeatures[i].Feat)
			if err != nil {
				continue
			}
			l, err := buff.Write(feats)
			if err != nil {
				log.Fatal(" buff.Write err : ", err)
			}
			len += l;
			incr++
			imageIds = append(imageIds, v.FaceImageId)
			peopleIds = append(peopleIds, v.PeopleId)
			if (i+1)%batch == 0 {
				ListInsert(&len, &incr, peopleIds, imageIds, r, buff)
				imageIds = imageIds[:0]
				peopleIds = peopleIds[:0]
			}
		}
		if len > 0 {
			ListInsert(&len, &incr, peopleIds, imageIds, r, buff)
		}
	}
}
func CreateReposFaceSdkCache(reposKey string) bool {
	created := G_ChlFaceX.ChlFaceSdkListCreate(reposKey, REPOS_MAX_FEATURE_SIZE)
	if !created {
		log.Fatalf("ChlFaceSdkListCreate[%s]  err : size[%d] ", reposKey, REPOS_MAX_FEATURE_SIZE)
		return false
	}
	if _, ok := G_ReposHandleMap[reposKey]; !ok {
		G_ReposHandleMap[reposKey] = set.New(set.NonThreadSafe)
	}
	return true
}

func ListInsert(length *int, incr *int, peopleIds []int64, imageIds []string, r *model.Repository, buff *bytes.Buffer) (err error) {
	defer func() {
		buff.Reset()
		*incr, *length = 0, 0
	}()
	if *length != *incr*2600 {
		log.Fatalf(" length[%d] !=incr[%d]*2600  err : ", length, incr, err)
	}
	fmt.Printf("ListInsert length [%d] :   %d ", *incr, len(buff.Bytes()[0:*length]))
	total, pos := G_ChlFaceX.ChlFaceSdkListInsert(BuildReposKey(r.Id), -1, buff.Bytes()[0:*length], *incr)
	fmt.Println("repos id : ", REPOS_PREFIX+r.Id, "  total: ", total, "  pos: ", pos)

	setcoll := G_ReposHandleMap[BuildReposKey(r.Id)]

	if setcoll == nil || pos == -1 {
		fmt.Printf("repos id [%s] , setcoll [%v] err :  total [%d], pos[%d] \n", r.Id, setcoll, total, pos)
		return nil
	}
	for i := 0; i < (*incr); i++ {
		entry := &model.FeatureEntry{PeopleId: peopleIds[i], FaceImageId: imageIds[i], Pos: pos + i, Status: 0}
		setcoll.Add(entry)
		fmt.Printf(" peopleIds[%d]  Pos : [%d] \n", peopleIds[i], pos+i)
	}
	return nil
}
func loadCameraFeature() {
	cameras := logic.DefaultCamera.FindOnlineCameras();
	for _, v := range cameras {
		created := G_ChlFaceX.ChlFaceSdkListCreate(BuildCameraKey(v.Id), CAMERA_MAX_FEATURE_SIZE)
		if !created {
			log.Fatal("ChlFaceSdkListCreate err: ", BuildCameraKey(v.Id))
		}
		if _, ok := G_CameraHandleMap[BuildCameraKey(v.Id)]; !ok {
			G_CameraHandleMap[BuildCameraKey(v.Id)] = set.New(set.NonThreadSafe)
		}
	}
}

func BuildReposKey(id string) string {
	return REPOS_PREFIX + id
}
func BuildCameraKey(id string) string {
	return CAMERA_PREFIX + id
}
