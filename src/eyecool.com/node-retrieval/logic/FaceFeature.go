package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"fmt"
)

type FaceFeatureLogic struct{}

var DefaultFaceFeature = FaceFeatureLogic{}

func (FaceFeatureLogic) FindFaceFeaturesByRepositoryId(repositoryId string) []*model.FaceFeature {
	objLog := GetLogger(nil)
	features := make([]*model.FaceFeature, 0)
	err := MasterDB.Where("repository_id=? and status=0  ", repositoryId).Find(&features)
	if err != nil {
		objLog.Errorln("FaceFeatureLogic FindFaceFeaturesByRepositoryId error:", err)
		return nil
	}
	return features
}

func (this *FaceFeatureLogic) FindFaceFeatureByPkId(pkId int) (bool, *model.FaceFeature) {
	feature := new(model.FaceFeature)
	has, err := MasterDB.ID(int64(pkId)).Get(feature)
	if err!=nil {
		fmt.Println(err)
	}
	return has, feature
}

func (this *FaceFeatureLogic) Insert(feature *model.FaceFeature) (*model.FaceFeature, error) {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(feature.TableName()).InsertOne(feature)
	session.Commit()
	return feature, err
}

func (this *FaceFeatureLogic) UpdateFaceImageId(feature *model.FaceFeature) (*model.FaceFeature, error) {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(feature.TableName()).ID(feature.PkId).Cols("face_image_id").Update(feature)
	session.Commit()
	return feature, err
}

func (this *FaceFeatureLogic) DeleteByPeopleId(feature *model.FaceFeature) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(feature.TableName()).Where("people_id = ?", feature.PeopleId).Cols("status", "update_time").Update(feature)
	session.Commit()
	return err
}
func (featureLogic *FaceFeatureLogic) FindFaceFeatureByPeopleId(peopleId int64) (bool, *model.FaceFeature) {
	faceFeature := new(model.FaceFeature)
	has, _ := MasterDB.Table(faceFeature.TableName()).Where("people_id = ?", peopleId).Get(faceFeature)
	return has, faceFeature
}
