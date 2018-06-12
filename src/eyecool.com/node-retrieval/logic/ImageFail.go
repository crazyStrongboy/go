package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"fmt"
)

type ImageFailLogic struct {
}

func (this *ImageFailLogic) Insert(fail *model.ImageFail) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.InsertOne(fail)
	if err != nil {
		session.Rollback()
		logger.Errorln("InsertUser ImageFail error:", err)
		session.Commit()
		return err
	}
	session.Commit()
	return nil
}
func (failLogic *ImageFailLogic) FindByRepositoryId(start int, limit int, repositoryId string) ([]*model.ImageFail) {
	imageFails := make([]*model.ImageFail, 0)
	err := MasterDB.Table(new(model.ImageFail)).Where("repository_id = ?", repositoryId).Limit(limit, start).Find(&imageFails)
	if err != nil {
		fmt.Println("FindByRepositoryId err :", err)
	}
	return imageFails
}
