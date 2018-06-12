package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"time"
	"log"
)

type RepositoryLogic struct{}

var DefaultRepository = RepositoryLogic{}

type Result struct {
	Id               string      `json:"id"`
	Name             string      `json:"name"`
	TotalPictureNum  int         `json:"total_picture_num"`
	FaceImageNum     int         `json:"face_image_num"`
	FailedPictureNum int         `json:"failed_picture_num"`
	CreatorId        int         `json:"creator_id"`
	CreateTime       int64       `json:"create_time"`
	PermissionMap    string      `json:"permission_map"`
	ExtraMeta        interface{} `json:"extra_meta"`
}

func (RepositoryLogic) FindAll() []*model.Repository {
	pageRepos := make([]*model.Repository, 0)
	err := MasterDB.Where("status=0").Find(&pageRepos)
	if err != nil {
		log.Fatal("RepositoryLogic FindAll  error:", err)
		return nil
	}
	return pageRepos
}

func (RepositoryLogic) QueryRepository() ([]*Result, error) {
	objLog := GetLogger(nil)
	results := make([]*Result, 0)
	repositorys := make([]*model.Repository, 0)
	err := MasterDB.Where("status=0").Find(&repositorys)
	if err != nil {
		objLog.Errorln("AdLogic FindAll PageAd error:", err)
		return nil, err
	}
	for _, v := range repositorys {
		result := &Result{}
		result.Id = v.Id
		result.Name = v.Name
		result.TotalPictureNum = v.TotalPictureNum
		result.FaceImageNum = v.FaceImageNum
		result.FailedPictureNum = v.FailedPictureNum
		result.CreatorId = v.CreatorId
		result.CreateTime = v.CreateTime
		result.PermissionMap = v.PermissionMap
		result.ExtraMeta = v.ExtraMeta
		results = append(results, result)
	}
	return results, nil
}

//查询该库是否存在
func (RepositoryLogic) SelectByName(name string) bool {
	objLog := GetLogger(nil)
	repository := model.Repository{}
	flag, err := MasterDB.Where("name=? and status=0", name).Get(&repository)
	if err != nil {
		objLog.Errorln("AdLogic FindAll PageAd error:", err)
		return false
	}
	return flag
}

//插入
func (RepositoryLogic) InsertRepository(repository *model.Repository) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Insert(repository)
	if err != nil {
		session.Rollback()
		logger.Errorln("insert alarmInfo error:", err)
		return err
	}
	session.Commit()
	return nil

}

func (RepositoryLogic) UpdateRepository(repository *model.Repository) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("pk_id=?", repository.PkId).Update(repository)
	if err != nil {
		session.Rollback()
		logger.Errorln("Update Repository error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (RepositoryLogic) DeleteRepository(id int) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	//删除人像库
	r := &model.Repository{
		PkId:       id,
		Status:     1,
		UpdateTime: time.Now(),
	}
	_, err := MasterDB.Where("pk_id=?", r.PkId).Update(r)
	if err != nil {
		session.Rollback()
		logger.Errorln("delete Repository error:", err)
		return err
	}
	//删除库下面的所有people
	people := &model.People{
		PeopleStatus:   1,
		RepositoryPkId: id,
		UpdateTime:     time.Now(),
	}
	_, err = MasterDB.Where("repository_pk_id=?", people.RepositoryPkId).Update(people)
	if err != nil {
		session.Rollback()
		logger.Errorln("delete people error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (RepositoryLogic) SelectRepositoryById(id int) bool {
	objLog := GetLogger(nil)
	repository := model.Repository{}
	flag, err := MasterDB.Where("id=? and status=0", id).Get(&repository)
	if err != nil {
		objLog.Errorln("RepositoryLogic selectRepositoryById PageAd error:", err)
		return false
	}
	return flag
}

func (this *RepositoryLogic) SelectByPrimaryKey(id int) (bool, *model.Repository) {
	repository := new(model.Repository)
	has, _ := MasterDB.ID(id).Get(repository)
	//fmt.Println(repository)
	return has, repository
}

func (this *RepositoryLogic) InsertFailRepository(repository *model.Repository) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.ID(repository.PkId).Cols("failed_picture_num", "update_time").Update(repository)
	if err != nil {
		session.Rollback()
		logger.Errorln("InsertRepository  error:", err)
		session.Commit()
		return err
	}
	session.Commit()
	return nil
}
