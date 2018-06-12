package controller

import (
	"github.com/emicklei/go-restful"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"github.com/satori/go.uuid"
	"eyecool.com/node-retrieval/model"
	"strings"
	"eyecool.com/node-retrieval/http/service"
	"eyecool.com/node-retrieval/global"
	"encoding/base64"
	"time"
	"eyecool.com/node-retrieval/utils"
)

var SurveillanceUrl = "http://192.168.0.192:8091/surveillance/api/feature/match"

type FaceController struct {
}

func (s *FaceController) Verify(req *restful.Request, rsp *restful.Response) {
	log.Print("Received AccessController.Verify API request")
	rsp.WriteEntity(map[string]string{
		"message": "Hi, this is the Verify API",
	})
	rsp.ResponseWriter.Write([]byte("xxxxxxxxResponseWriter Verify xxxxxxxx"))
}

func (this *FaceController) InsertOrigImage(req *restful.Request, rsp *restful.Response) {
	log.Print("Received FaceController.InsertOrigImage API request : ", req.Request.RemoteAddr)
	oi := model.OrigImageRequest{}
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &oi)
	if err != nil {
		fmt.Println("Unmarshal OrigImageRequest err : ", err)
	}
	u1, _ := uuid.NewV4()
	if oi.Uuid == "" {
		oi.Uuid = strings.Replace(u1.String(), "-", "", -1)
	}
	//推送到布控系统
	this.pushToSurveillance(&oi)

	//数据入库
	err = service.InsertOrigImage(nil, &oi)
	if err != nil {
		fmt.Println("insert error : ", err)
	}
	resp := make(map[string]string)
	resp["msg"] = "successs"
	responseBytes, _ := json.Marshal(resp)
	rsp.ResponseWriter.Write(responseBytes)
}

func (s *FaceController) pushToSurveillance(oi *model.OrigImageRequest) {
	//开启 推送到布控系统
	if global.G_Push_Pattern {
		featArr := strings.Split(oi.FeatList, ",")
		featureBufs := make([]byte, 0, global.FEATURE_LENGTH*oi.FaceNum)

		for i, _ := range featArr {
			featBytes, err := base64.StdEncoding.DecodeString(featArr[i])
			if err != nil {
				fmt.Println(" DecodeString error : ", err)
				continue
			}
			featureBufs = append(featureBufs[:], featBytes[:]...)
		}
		imageSource := &model.ImageSource{
			Type:              global.Type_Feature, // 0:image ,1:feature
			CameraId:          oi.CameraId,
			CameraIp:          oi.CameraIp,
			CreateTime:        time.Now().UnixNano(),
			CaptureTime:       oi.Timestamp,
			OrigPath:          oi.ImageRealPath,
			ImageOriginalName: oi.ImageName,
			ImageContextPath:  oi.ImageContextPath,
			OrigImageUri:      oi.FaceImageUri,
			OrigImageUuid:     oi.Uuid,
			FaceNum:           int(oi.FaceNum),
			FaceRects:         oi.FaceRect,
			FaceFeatureBufs:   featureBufs,
		}
		imageSourceBytes, _ := json.Marshal(imageSource)
		fmt.Println("DoBytesPost imageSourceBytes len : ", len(imageSourceBytes), "  featureBufs len :", len(featureBufs))
		//kafka
		_, err := utils.DoBytesPost(SurveillanceUrl, imageSourceBytes)
		if err != nil {
			fmt.Println("byte err : ", err)
		}
	}
}
