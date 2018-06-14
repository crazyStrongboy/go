package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"encoding/json"
	"strconv"
	"time"
	"fmt"
	"log"
	"github.com/polaris1119/logger"
)

type CameraLogic struct{}

var DefaultCamera = CameraLogic{}

type Cameras struct {
	Id             string      `json:"id"`
	Name           string      `json:"name"`
	Url            string      `json:"url"`
	Enabled        int         `json:"enabled"`
	RecParams      string      `json:"rec_params"`
	PermissionMap  string      `json:"permission_map"`
	PredecessorIds []string    `json:"predecessor_ids"`
	ExtraMeta      interface{} `json:"extra_meta"`
	Status         int         `json:"status"`
}

func (CameraLogic) FindOnlineCameras() []*model.Camera {
	cameras := make([]*model.Camera, 0)
	err := MasterDB.Where("status=1 or status=2 ").Find(&cameras)
	if err != nil {
		log.Println("AdLogic FindAll PageAd error:", err)
		return nil
	}
	return cameras
}

func (CameraLogic) Insert(camera *model.Camera) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	_, err := MasterDB.Insert(camera)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("insert camera error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (CameraLogic) CameraQuery() []*Cameras {
	cameras := make([]*model.Camera, 0)
	err := MasterDB.Where("status!=6").Find(&cameras)
	if err != nil {
		log.Println("AdLogic FindAll PageAd error:", err)
		return nil
	}
	result := make([]*Cameras, 0)

	for _, v := range cameras {
		camera := &Cameras{}
		camera.Id = v.Id
		camera.Name = v.Name
		camera.Url = v.Url
		var enabled int
		if (v.Status == 5 || v.Status == 0) {
			enabled = 0;
		} else {
			enabled = 1;
		}
		camera.Enabled = enabled
		recJson, _ := json.Marshal(v.RecParams)
		camera.RecParams = string(recJson)
		perJson, _ := json.Marshal(v.PermissionMap)
		camera.PermissionMap = string(perJson)
		camera.ExtraMeta = v.ExtraMeta
		camera.Status = v.Status
		if l := len(v.PredecessorId); l > 0 {
			id, _ := strconv.Atoi(v.PredecessorId)
			camera.PredecessorIds = DefaultCamera.FindRegionList(id)
		}
		result = append(result, camera)

	}
	return result

}

func (CameraLogic) FindCameraById(id int) (bool, *model.Camera) {
	camera := model.Camera{}
	flag, err := MasterDB.Where("id=? and status!=6", id).Get(&camera)
	if err != nil {
		log.Println("CameraLogic selectCameraById PageAd error:", err)
		return false, nil
	}
	return flag, &camera
}

func (CameraLogic) FindRegionList(id int) []string {
	result := make([]string, 0)
	region := model.Region{}
	flag := true
	for flag {
		flg, _ := MasterDB.Where("id=?", id).Get(&region)
		if flg {
			regionId := strconv.Itoa(region.Id)
			cluster := strconv.Itoa(region.ClusterId)
			result = append(result, regionId+"@"+cluster)
			if region.ParentId == 0 {
				flag = false
			} else {
				id = region.ParentId
			}
		} else {
			flag = false
		}
	}
	return result

}

func (CameraLogic) FindIP(ip string) bool {
	camera := model.Camera{}
	flg, _ := MasterDB.Where("ip=? and status !=6 ", ip).Get(&camera)
	return flg
}

func (CameraLogic) DeleteCamera(id int) error {
	camera := &model.Camera{
		PkId:       id,
		Status:     6,
		UpdateTime: time.Now(),
	}
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	//更新camera状态
	_, err := MasterDB.Where("pk_id=?", id).Update(camera)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("delete camera error:", err)
		return err
	}
	//同步更新video_camera状态
	vc := &model.VideoCamera{
		Status: 6,
	}
	_, err = MasterDB.Where("camera_id=?", id).Update(vc)
	if err != nil {
		session.Rollback()
		logger.Errorln("delete camera error:", err)
		return err
	}
	session.Commit()
	return nil

}

func (CameraLogic) Update(camera *model.Camera) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("pk_id=?", camera.PkId).Update(camera)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("update camera error:", err)
		return err
	}
	session.Commit()
	return nil
}
func (CameraLogic) InsertVideoCamera(vc *model.VideoCamera) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	_, err := MasterDB.Insert(vc)
	if err != nil {
		fmt.Println(err)
		session.Rollback()
		session.Commit()
		log.Println("insert videoCamera error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (CameraLogic) UpdateVideoCamera(vc *model.VideoCamera) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("camera_id=?", vc.Id).Update(vc)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("delete camera error:", err)
		return err
	}
	session.Commit()
	return nil
}
