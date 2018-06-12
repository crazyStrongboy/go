package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"log"
)

type OrigImageLogic struct{}

var DefaultOrigImage = OrigImageLogic{}

func (OrigImageLogic) FindAllByCameraId(camera_id string) []*model.OrigImage {
	objLog := GetLogger(nil)
	origImages := make([]*model.OrigImage, 0)
	err := MasterDB.Where("camera_id=?", camera_id).Find(&origImages)
	if err != nil {
		objLog.Errorln("OrigImageLogic FindAllByCameraId error:", err)
		return nil
	}
	return origImages
}

type OrigImageFullLogic struct{}

var DefaultOrigImageFull = OrigImageFullLogic{}

func (OrigImageFullLogic) Insert(image *model.OrigImageFull) error {
	//logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	_, err := MasterDB.Insert(image)
	if err != nil {
		session.Rollback()
		log.Println("insert image error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (this *OrigImageLogic) FindOrigImages(clusterId int, offset int, limit int) ([]*model.OrigImageFull, error) {
	origImages := make([]*model.OrigImageFull, 0)
	err := MasterDB.Table(new(model.OrigImageFull).TableName()).Where("cluster_id = ?", clusterId).Desc("update_time").Limit(limit, offset).Find(&origImages)
	return origImages, err
}
func (imageLogic *OrigImageLogic) FindOrigImageById(id int64) (bool, *model.OrigImageFull) {
	origImage := new(model.OrigImageFull)
	has, _ := MasterDB.Table(new(model.OrigImageFull).TableName()).ID(id).Get(origImage)
	return has, origImage
}
