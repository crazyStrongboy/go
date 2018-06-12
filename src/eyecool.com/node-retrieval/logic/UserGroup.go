package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"time"
	"github.com/polaris1119/logger"
)

type UserGroupLogic struct {
}

func (this *UserGroupLogic) UpdateStatus(groupId int, status int) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	group := new(model.UserGroup)
	group.Status = status
	_, err := MasterDB.Table(new(model.UserGroup).TableName()).ID(groupId).Cols("status").Update(group)
	session.Commit()
	return err
}

func (this *UserGroupLogic) SelectByParentId(parentId int) ([]*model.UserGroup, error) {
	groups := make([]*model.UserGroup, 0)
	err := MasterDB.Table(new(model.UserGroup).TableName()).Where("parent_id = ?", parentId).Find(&groups)
	return groups, err
}

func (this *UserGroupLogic) InsertGroup(group *model.UserGroup) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(group.TableName()).InsertOne(group)
	session.Commit()
	return err
}

func (this *UserGroupLogic) SelectGroupById(id int) (*model.UserGroup) {
	group := new(model.UserGroup)
	MasterDB.Table(group.TableName()).ID(id).Get(group)
	return group
}

func (this *UserGroupLogic) SelectPredecessorIds(parentId int) []int {
	groupPredecessorIds := make([]int, 0)
	logger := GetLogger(nil)
	err := MasterDB.Table(new(model.UserGroup).TableName()).Cols("id").Where("parent_id=? and status=0", parentId).Find(&groupPredecessorIds)
	if err != nil {
		logger.Fatal("SelectAllTopUserGroup error:", err)
	}
	return groupPredecessorIds
}

func (this *UserGroupLogic) SelectAllTopUserGroup(userGroup *model.UserGroup) []*model.UserGroup {
	userGroups := make([]*model.UserGroup, 0)
	logger := GetLogger(nil)
	err := MasterDB.Table(userGroup.TableName()).Find(&userGroups)
	if err != nil {
		logger.Fatal("SelectAllTopUserGroup error:", err)
	}
	return userGroups
}

func (this *UserGroupLogic) SelectGroupLevelById(id int) (bool,int) {
	groupLevel := 0
	has, _ := MasterDB.Table(new(model.UserGroup).TableName()).Cols("group_level").Where("id = ?", id).Get(&groupLevel)
	return has,groupLevel
}

func (this *UserGroupLogic) SelectUserGroupByLevel(level int) []*model.UserGroup {
	userGroups := make([]*model.UserGroup, 0)
	MasterDB.Table(new(model.UserGroup).TableName()).Where("group_level = ? and status=0", level).Find(&userGroups)
	return userGroups
}

func (this *UserGroupLogic) UpdateUserGroup(group *model.UserGroup) error {
	group.UpdateTime = time.Now()
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(group.TableName()).ID(group.Id).Cols("name").Cols("extra_meta").Update(group)
	if err != nil {
		session.Rollback()
		logger.Errorln("update userGroup error:", err)
		session.Commit()
		return err
	}
	session.Commit()
	return nil
}
