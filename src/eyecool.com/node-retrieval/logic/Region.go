package logic

import( "eyecool.com/node-retrieval/model"
. "eyecool.com/node-retrieval/db"
)
type RegionLogic struct{}

var DefaultRegion =RegionLogic{}

type Sets struct {
	Id string `json:"id"`
	Name string `json:"name"`
	PermissionMap string `json:"permission_map"`
	PredecessorIds []string `json:"predecessor_ids"`
}

func (RegionLogic)QueryRegion()([]*Sets,error){
	objLog := GetLogger(nil)
	sets:=make([]*Sets, 0)
	cameras := make([]*model.Camera, 0)
	err := MasterDB.Where("status!=0").Find(&cameras)
	if err != nil {
		objLog.Errorln("AdLogic FindAll PageAd error:", err)
		return nil,err
	}
	for _,v:=range cameras{
		set:=&Sets{}
		set.Id=v.Id
		set.Name=v.Name
		set.PermissionMap=v.PermissionMap
		set.PredecessorIds=DefaultCamera.GetRegionList(v.RegionId)
		sets=append(sets, set)
	}
	return sets,nil
}

func (RegionLogic)SelectByPrimaryKey(id int) model.Region{
	region:=model.Region{}
	MasterDB.Where("id=?",id).Get(&region)
	return region
}

func (RegionLogic)InsertRegion(region *model.Region)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()

	_, err := MasterDB.Insert(region)
	if err != nil {
		session.Rollback()
		logger.Errorln("insert region error:", err)
		return  err
	}
	session.Commit()
	return nil
}

func (RegionLogic) UpdateRegion(region *model.Region)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_,err:=MasterDB.Where("id=?",region.Id).Update(region)
	if err != nil {
		session.Rollback()
		logger.Errorln("update region error:", err)
		return  err
	}
	session.Commit()
	return nil
}

func (RegionLogic) DeleteRegion(id int) error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	regions:=selectByParentId(id)
	region:=&model.Region{
		Id:     id,
		Status: 2,
	}
	_,err:=MasterDB.Where("id=?",region.Id).Update(region)
	if err != nil {
		session.Rollback()
		logger.Errorln("update region error:", err)
		return  err
	}
	for _,v:=range regions{
		DefaultRegion.DeleteRegion(v.Id)
	}

	session.Commit()
	return nil
}

func selectByParentId(parentId int) []*model.Region{
	regions:=make([]*model.Region,0)
	MasterDB.Where("parent_id=?",parentId).Find(&regions)
	return regions
}
