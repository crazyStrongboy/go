package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"fmt"
)

type ImageLogic struct {
}

func (this *ImageLogic) Insert(image *model.Image) (*model.Image, error) {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(image.TableName()).InsertOne(image)
	session.Commit()
	return image, err
}

func (this *ImageLogic) DeleteByPeopleId(image *model.Image) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(image.TableName()).Where("people_id = ?", image.PeopleId).Cols("status", "update_time").Update(image)
	session.Commit()
	return err
}


func (this *ImageLogic) FindImageById(id int64) (bool, *model.Image) {
	image := new(model.Image)
	has, err := MasterDB.Table(image.TableName()).ID(id).Get(image)
	if err!=nil {
		fmt.Println(err)
	}
	return has, image
}