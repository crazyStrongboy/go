package buz

import (
	"eyecool.com/node-retrieval/utils"
	"eyecool.com/node-retrieval/model"
	"eyecool.com/node-retrieval/logic"
	"time"
	"fmt"
	"strconv"
	"encoding/base64"
	"eyecool.com/node-retrieval/http/service"
	"os"
	"io/ioutil"
	"eyecool.com/node-retrieval/algorithm"
	. "github.com/polaris1119/config"
)

var (
	people_target_path = "/home/eyecool/imageSavePath/tempImage/"
	nginx_path         = "http://192.168.0.192/"
)

type PeopleService struct {
}

// {"x": 左上角x坐标, "y": 左上角y坐标, "w": 宽度, "h": 高度}
type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"h"`
	W int `json:"w"`
}

//表示图片入库的结果
type Result struct {
	FaceImageId  string `json:"face_image_id,omitempty"`  //此人入库后的 id ($编号@$集群号)
	FaceImageUri string `json:"face_image_uri,omitempty"` //此人人脸图片的 uri，在单集群接口的基础上加上 "@$集群号
	FaceRecg     *Rect  `json:"face_recg,omitempty"`      //表示人脸在原图中的位置, {"x": 左上角x坐标, "y": 左上角y坐标, "w": 宽度, "h": 高度}
}

type PeopleResponse struct {
	Rtn        int       `json:"rtn"`               //错误码
	Message    string    `json:"message,omitempty"` //错误消息
	Results    []*Result `json:"results,omitempty"`
	PictureUri string    `json:"picture_uri,omitempty"` //表示入库图片的 uri ，在单集群接口的基础上加上 "@$集群号"
}

type UpdateInfo struct {
	Person_id string //证件号可以是18位身份证号或者其他格式的证件号码
	Name      string //查询人姓名
	Born_year string //出生日期。格式为"YYYY-mm-dd"， 比如"1990-10-10
	Gender    string //性别。0未知, 1男, 2女
}

type PeopleRequest struct {
	Repository_id                string      //图片所在库的id ($编号@$集群号)
	Picture_image_content_base64 string      //base64编码后的照片内容
	Name                         string      //查询人姓名
	Region                       int         //区域编号
	Birthday                     string      //出生日期。格式为"YYYY-mm-dd"， 比如"1990-10-10
	Gender                       int         //性别。0未知, 1男, 2女
	Nation                       string      //民族, 见 , 0表示未知
	Person_id                    string      //证件号可以是18位身份证号或者其他格式的证件号码
	Options                      string      //建库选项
	Custom_field                 string      //导图时添加任意多的额外自定义字段。
	Face_image_id                string      //要修改的 人像库中face 的face_image_id
	Set                          *UpdateInfo //里面填写每个域表示需要更新的域, 可以修改的阈值包括
}

func init() {
	people_target_path, _ = ConfigFile.GetValue("path", "people_target_path")
	nginx_path, _ = ConfigFile.GetValue("path", "nginx_path")
}

var peopleLogic = new(logic.PeopleLogic)
var featureLogic = new(logic.FaceFeatureLogic)
//导入图片
func (this *PeopleService) Insert(request *PeopleRequest, userId int) *PeopleResponse {
	response := new(PeopleResponse)
	people := new(model.People)
	repositoryPkId := 1
	clusterId := 1
	repository := new(model.Repository)
	var err error
	if request.Repository_id != "" {
		repositoryPkId, clusterId, err = utils.GetIdAndClusterId(request.Repository_id)
		if err != nil {
			response.Rtn = -1
			response.Message = err.Error()
			return response
		}
		people.RepositoryId = request.Repository_id
		var has bool
		has, repository = repositoryLogic.SelectByPrimaryKey(repositoryPkId)
		if !has || (has && repository.Status == 1) {
			response.Rtn = -1
			response.Message = "库不存在!"
			return response
		}
	}
	people.ClusterId = clusterId
	people.RepositoryPkId = repositoryPkId
	//存图片到本地,并检测是否有人脸
	response, err = writeFileToDiskAndDetect(request, repository, people, clusterId, userId)

	return response
}

func (this *PeopleService) Update(request *PeopleRequest) *PeopleResponse {
	response := new(PeopleResponse)
	facePkIdStr := request.Face_image_id
	if facePkIdStr == "" {
		response.Rtn = -1
		response.Message = "Face_image_id不能为空!"
		return response
	}
	facePkId, _, err := utils.GetIdAndClusterId(facePkIdStr)
	if err != nil || facePkId == -2 {
		response.Rtn = -1
		response.Message = "Face_image_id参数错误!"
		return response
	}
	hasF, feature := featureLogic.FindFaceFeatureByPkId(facePkId)
	if hasF {
		hasP, people := peopleLogic.FindPeopleById(feature.PeopleId)
		if hasP {
			if request.Set == nil {
				response.Rtn = -1
				response.Message = "Set参数不能为空!"
				return response
			}
			personId := request.Set.Person_id
			name := request.Set.Name
			birthday := request.Set.Born_year
			genderStr := request.Set.Gender
			if personId != "" {
				people.PersonId = personId
			}
			if name != "" {
				people.Name = name
			}
			if birthday != "" {
				people.Birthday = birthday
			}
			if genderStr != "" {
				gender, err := strconv.Atoi(genderStr)
				if err == nil {
					people.Gender = gender
				} else {
					fmt.Println("strconv.Atoi(genderStr)", err)
				}
			}
			people.UpdateTime = time.Now()
			err = peopleLogic.UpdateById(people)
			if err != nil {
				response.Rtn = -1
				response.Message = "更新失败!"
				return response
			}
			response.Rtn = 0
			response.Message = "更新成功!"
		}
	}
	return response
}

func (this *PeopleService) Delete(request *PeopleRequest) *PeopleResponse {
	response := new(PeopleResponse)
	faceImageId := request.Face_image_id
	faceId, _, err := utils.GetIdAndClusterId(faceImageId)
	if err != nil || faceId == -2 {
		response.Rtn = -1
		response.Message = "参数错误!"
		return response
	}
	_, feature := featureLogic.FindFaceFeatureByPkId(faceId)
	peopleId := feature.PeopleId
	deletePeople(peopleId)
	deleteImage(peopleId)
	deleteFeature(peopleId)
	response.Rtn = 0
	response.Message = "删除成功!"

	// 通知go 删除缓存
	notifyDeleteCache("", "", 1, faceImageId)
	return response
}

//retrieval/repository_feature_insert
//删除缓存
func notifyDeleteCache(enCodeB64 string, repositoryId string, kind int32, faceId string) {
	featbytes, _ := base64.StdEncoding.DecodeString(enCodeB64)
	retrievalFeatureRequest := &model.RetrievalFeatureRequest{
		RepositoryId: repositoryId,
		Id:           faceId,
		PeopleId:     "",
		Type:         kind,
		FeatNum:      1,
		Feats:        featbytes,
	}
	service.RetrievalRepositoryFeatureInsert(nil, retrievalFeatureRequest, nil)
}
func deleteFeature(peopleId int64) {
	feature := new(model.FaceFeature)
	feature.PeopleId = peopleId
	feature.Status = 2
	feature.UpdateTime = time.Now()
	featureLogic.DeleteByPeopleId(feature)
}
func deleteImage(peopleId int64) {
	image := new(model.Image)
	image.PeopleId = peopleId
	image.Status = 2
	image.UpdateTime = time.Now()
	imageLogic.DeleteByPeopleId(image)
}
func deletePeople(peopleId int64) {
	people := new(model.People)
	people.Id = peopleId
	people.PeopleStatus = 2
	people.UpdateTime = time.Now()
	peopleLogic.DeleteById(people)
}

//存图片到本地
func writeFileToDiskAndDetect(request *PeopleRequest, repository *model.Repository, people *model.People, clusterId int, userId int) (*PeopleResponse, error) {
	response := new(PeopleResponse)
	datePath := time.Now().Format("20060102")
	targetPath := people_target_path + "/" + datePath
	if _, err := os.Stat(targetPath); err != nil {
		os.MkdirAll(targetPath, os.ModePerm)
	}
	imageName := utils.UUID() + ".jpg"
	realPath := targetPath + "/" + imageName
	imageUri := "tempImage/" + datePath + "/" + imageName
	data, _ := base64.StdEncoding.DecodeString(request.Picture_image_content_base64)
	err := ioutil.WriteFile(realPath, data, os.ModePerm)
	if err != nil {
		fmt.Println("people insert write error", err, "path:", realPath)
		response.Rtn = -1
		response.Message = "文件写入失败"
		return response, err
	}
	//检测是否有人脸
	_, width, height, rgb24Data := algorithm.NewChlFaceX().ReadImageFile(realPath, 0, 0)
	hasFace, faceResult := algorithm.NewChlFaceX().ChlFaceSdkDetectFace(0, rgb24Data, width, height, true)
	fmt.Println("width:", width, "height:", height, "hasFace:", hasFace)
	if hasFace == 0 || faceResult == nil {
		//检测不到人脸
		response.Rtn = -1
		response.Message = "检测不到人脸"
		err = insertImageFail(request, clusterId, userId, imageName, realPath, imageUri)
		if err != nil {
			fmt.Println("insertImageFail error", err)
			response.Message = "insertImageFail 失败"
		} else {
			//检测不到人脸的图片入库成功,则该repository失败图片+1
			repository.UpdateTime = time.Now()
			repository.FailedPictureNum = repository.FailedPictureNum + 1
			err = repositoryLogic.InsertFailRepository(repository)
			if err != nil {
				fmt.Println("InsertRepository error", err)
				response.Message = "insertImageFail 失败"
			}
		}
		return response, err
	}
	//入库上传的people
	people, err = insertPeople(request, people)
	if err != nil {
		response.Rtn = -1
		response.Message = "insertPeople failed!!"
		return response, err
	}
	//入库上传的图片
	image, err := insertImage(clusterId, imageUri, people, realPath, imageName)
	if err != nil {
		response.Rtn = -1
		response.Message = "insertImage failed!!"
		return response, err
	}
	//提取特征
	response = extractFeatureAndInsert(data, faceResult, people, image, width, height)
	return response, nil
}
func insertImage(clusterId int, imageUri string, people *model.People, realPath string, imageName string) (*model.Image, error) {
	image := new(model.Image)
	image.CreateTime = time.Now()
	image.UpdateTime = time.Now()
	image.PubId = utils.UUID()
	image.ClusterId = clusterId
	image.ImageContextPath = nginx_path
	image.ImageUrl = nginx_path + "/" + imageUri
	image.RepositoryId = people.RepositoryPkId
	image.ImageRealPath = realPath
	image.ImageUri = imageUri
	image.Status = 0
	image.ImageType = 2
	image.ImageName = imageName
	image.PeopleId = people.Id
	var err error
	image, err = imageLogic.Insert(image)
	return image, err
}

func insertPeople(request *PeopleRequest, people *model.People) (*model.People, error) {
	if request.Person_id != "" {
		people.PersonId = request.Person_id
	}
	if request.Custom_field != "" {
		people.CustomField = request.Custom_field
	}
	people.CreateTime = time.Now()
	people.UpdateTime = time.Now()
	people.PubId = utils.UUID()
	people.PeopleStatus = 0
	people.CreatorId = 1
	people.Name = request.Name
	people.Gender = request.Gender
	people.Region = request.Region
	people.Birthday = request.Birthday
	people.Nation = request.Nation
	people.Options = request.Options
	people.CustomField = request.Custom_field
	_, err := peopleLogic.Insert(people)
	return people, err
}

//提取特征并入库
func extractFeatureAndInsert(data []byte, faceResult *algorithm.FACE_DETECT_RESULTX, people *model.People, image *model.Image, width int, height int) (*PeopleResponse) {
	_, feature := algorithm.NewChlFaceX().ChlFaceSdkFeatureGet(0, data, width, height, faceResult)
	response := new(PeopleResponse)
	rect := faceResult.GetRECT()
	responseRect := new(Rect)
	responseRect.X = rect.Left
	responseRect.Y = rect.Right
	responseRect.W = rect.Top
	responseRect.H = rect.Bottom
	featBase64 := base64.StdEncoding.EncodeToString(feature)
	insertAndUpdateFeature(featBase64, responseRect, image, people, response)
	return response
}
func insertAndUpdateFeature(featBase64 string, rect *Rect, image *model.Image, people *model.People, response *PeopleResponse) {
	faceFeature := new(model.FaceFeature)
	faceFeature.Feat = featBase64
	faceFeature.CreateTime = time.Now()
	faceFeature.UpdateTime = time.Now()
	faceFeature.W = rect.W
	faceFeature.X = rect.X
	faceFeature.Y = rect.Y
	faceFeature.H = rect.H
	faceFeature.ImageId = image.Id
	faceFeature.PeopleId = people.Id
	faceFeature.RepositoryPkId = people.RepositoryPkId
	faceFeature.RepositoryId = people.RepositoryId
	faceFeature.Status = 0
	_, err := featureLogic.Insert(faceFeature)
	if err != nil {
		response.Rtn = -1
		response.Message = "Insert feature failed!!"
		return
	}
	faceFeature.FaceImageId = strconv.Itoa(faceFeature.PkId) + "@" + strconv.Itoa(image.ClusterId)
	_, err = featureLogic.UpdateFaceImageId(faceFeature)
	if err != nil {
		response.Rtn = -1
		response.Message = "UpdateFaceImageId failed!!"
		return
	}
	//导入成功后的返回结果
	response.Rtn = 0
	response.Message = "导图成功!"
	response.PictureUri = image.ImageUrl + "@" + strconv.Itoa(image.ClusterId)
	results := make([]*Result, 1)
	result := new(Result)
	result.FaceImageId = faceFeature.FaceImageId
	result.FaceImageUri = response.PictureUri
	result.FaceRecg = rect
	results[0] = result
	response.Results = results
}

func insertImageFail(request *PeopleRequest, clusterId int, userId int, imageName string, realPath string, imageUri string) error {
	imageFail := new(model.ImageFail)
	imageFail.CreateTime = time.Now()
	imageFail.Birthday = request.Birthday
	imageFail.ClusterId = clusterId
	imageFail.CreatorId = userId
	imageFail.Gender = request.Gender
	imageFail.ImageContextPath = nginx_path
	imageFail.ImageDesc = "图片检测不到人脸"
	imageFail.ImageName = imageName
	imageFail.ImageRealPath = realPath
	imageFail.ImageUri = imageUri
	imageFail.ImageUrl = nginx_path + "/" + imageUri
	imageFail.Name = request.Name
	imageFail.Nation = request.Nation
	imageFail.PersonId = request.Person_id
	imageFail.Region = request.Region
	imageFail.RepositoryId = request.Repository_id
	err := imageFailLogic.Insert(imageFail)
	return err
}
