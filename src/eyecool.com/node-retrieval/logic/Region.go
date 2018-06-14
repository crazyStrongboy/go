package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"log"
)

type RegionLogic struct{}

var DefaultRegion = RegionLogic{}

type Sets struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	PermissionMap  string   `json:"permission_map"`
	PredecessorIds []string `json:"predecessor_ids"`
}

func (RegionLogic) QueryRegion() ([]*Sets, error) {
	sets := make([]*Sets, 0)
	cameras := make([]*model.Camera, 0)
	err := MasterDB.Where("status!=0").Find(&cameras)
	if err != nil {
		log.Println("AdLogic FindAll PageAd error:", err)
		return nil, err
	}
	for _, v := range cameras {
		set := &Sets{}
		set.Id = v.Id
		set.Name = v.Name
		set.PermissionMap = v.PermissionMap
		set.PredecessorIds = DefaultCamera.FindRegionList(v.RegionId)
		sets = append(sets, set)
	}
	return sets, nil
}

func (RegionLogic) FindByPrimaryKey(id int)(bool, *model.Region ){
	region := model.Region{}
	has,err:=MasterDB.Where("id=?", id).Get(&region)
	if err!=nil {
		log.Println("FindByPrimaryKey region err : ",err)
	}
	return has,&region
}

func (RegionLogic) InsertRegion(region *model.Region) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Insert(region)
	if err != nil {
		session.Rollback()
		log.Println("insert region error:", err)
		session.Commit()
		return err
	}
	session.Commit()
	return nil
}

func (RegionLogic) UpdateRegion(region *model.Region) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("id=?", region.Id).Update(region)
	if err != nil {
		session.Rollback()
		logger.Errorln("update region error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (RegionLogic) DeleteRegion(id int) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	regions := findByParentId(id)
	region := &model.Region{
		Id:     id,
		Status: 2,
	}
	_, err := MasterDB.Where("id=?", region.Id).Update(region)
	if err != nil {
		session.Rollback()
		log.Println("update region error:", err)
		return err
	}
	for _, v := range regions {
		DefaultRegion.DeleteRegion(v.Id)
	}

	session.Commit()
	return nil
}

func findByParentId(parentId int) []*model.Region {
	regions := make([]*model.Region, 0)
	MasterDB.Where("parent_id=?", parentId).Find(&regions)
	return regions
}
