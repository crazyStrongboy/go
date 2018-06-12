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
	"time"
)

const (
	FEATURE_LENGTH              = 2600
	REPOS_MAX_FEATURE_SIZE  int = 20000
	CAMERA_MAX_FEATURE_SIZE int = 20000
	REPOS_PREFIX                = "repos_"
	CAMERA_PREFIX               = "camera_"
	INSERT_BATCH_CNT            = 1
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
	fmt.Println(" xxxx start chlface init ....")
	ret := G_ChlFaceX.ChlFaceSdkInit()
	if ret == -1 {
		log.Fatalf(" ChlFaceSdkInit[%d] faild ! ", ret)
	} else {
		log.Print("algorithm chlface init success !!! ")
	}
	//首次加载库特征数据到目标库
	loadRepositoryFeatures()
	//启动刷新缓存库的定时器
	startRefreshCachedReposTicker()
}

//定时刷新缓存目标库,更新数据到算法缓存中
func startRefreshCachedReposTicker() {
	t := time.NewTicker(30 * time.Second)
	go func() {
		i := 1
		for {
			<-t.C
			tNow := time.Now()
			timeNow := tNow.Format("2006-01-02 15:04:05")
			fmt.Printf("%s start %d refresh repos\n", timeNow, i)
			loadAllVaildRepos()
			i++
		}
	}()
	return
}

func loadRepositoryFeatures() {
	fmt.Println("xxx global start do loadRepositoryFeatures ........")
	repos := loadAllVaildRepos()
	if repos == nil || len(repos) == 0 {
		fmt.Println("load global  Repository is empty !!!")
		return
	}
	buff := new(bytes.Buffer)
	for _, r := range repos {
		RepositoryFeatureInsertCache(r.Id, buff)
	}
}
func loadAllVaildRepos() []*model.Repository {
	//加载有效的目标库数据
	repos := logic.DefaultRepository.FindAll()
	for _, v := range repos {
		ret := CreateReposFaceSdkCache(BuildReposKey(v.Id))
		fmt.Printf("exist repos id [%s%s] , cached [%t] \n", REPOS_PREFIX, v.Id, ret)
	}
	//做一些清除多余的数据操作
	for key := range G_ChlFaceX.HandlesCachedSize {
		exist := false
		//查找是不是库里应该加入的
		for _, v := range repos {
			if key == BuildReposKey(v.Id) {
				exist = true
			}
		}
		//如果不是则要移除库的缓存数据
		if !exist {
			G_ChlFaceX.ChlFaceSdkListDestroy(key)
			if set, ok := G_ReposHandleMap[key]; ok {
				set.Clear()
				delete(G_ReposHandleMap, key)
			}
		}
	}

	return repos
}

func RepositoryFeatureInsertCache(reposId string, buff *bytes.Buffer) {
	faceFeatures := logic.DefaultFaceFeature.FindFaceFeaturesByRepositoryId(reposId)
	incr, len := 0, 0
	imageIds := make([]string, 0, INSERT_BATCH_CNT)
	peopleIds := make([]int64, 0, INSERT_BATCH_CNT)
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
		if (i+1)%INSERT_BATCH_CNT == 0 {
			ListInsert(&len, &incr, peopleIds, imageIds, reposId, buff)
			imageIds = imageIds[:0]
			peopleIds = peopleIds[:0]
		}
	}
	if len > 0 {
		ListInsert(&len, &incr, peopleIds, imageIds, reposId, buff)
	}
}

func ListInsert(length *int, incr *int, peopleIds []int64, imageIds []string, reposId string, buff *bytes.Buffer) (err error) {
	defer func() {
		buff.Reset()
		*incr, *length = 0, 0
	}()
	if *length != *incr*2600 {
		log.Fatalf(" length[%d] !=incr[%d]*2600  err : ", length, incr, err)
	}
	fmt.Printf("ListInsert length [%d] :   %d ", *incr, len(buff.Bytes()[0:*length]))
	total, pos := G_ChlFaceX.ChlFaceSdkListInsert(BuildReposKey(reposId), -1, buff.Bytes()[0:*length], *incr)
	fmt.Println("repos id : ", REPOS_PREFIX+reposId, "  total: ", total, "  pos: ", pos)

	setcoll := G_ReposHandleMap[BuildReposKey(reposId)]

	if setcoll == nil || pos == -1 {
		fmt.Printf("repos id [%s] , setcoll [%v] err :  total [%d], pos[%d] \n", reposId, setcoll, total, pos)
		return nil
	}
	for i := 0; i < (*incr); i++ {
		entry := &model.FeatureEntry{PeopleId: peopleIds[i], FaceImageId: imageIds[i], Pos: pos + i, Status: 0}
		setcoll.Add(entry)
		fmt.Printf(" peopleIds[%d]  Pos : [%d] \n", peopleIds[i], pos+i)
	}
	return nil
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
