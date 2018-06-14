package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"time"
	"log"
)

type UserLogic struct {
}

var userGroupLogic = new(UserGroupLogic)

func (this *UserLogic) FindUserByName(name string) (*model.User, bool, error) {
	existUser := new(model.User)
	has, err := MasterDB.Table(existUser.TableName()).Where("name = ? and status = 0", name).Get(existUser)
	if err != nil {
		log.Println("FindUserByName err :", err)
	}
	return existUser, has, err
}

func (this *UserLogic) FindUsersByGroupId(groupId int) ([]*model.User, error) {
	userList := make([]*model.User, 0)
	err := MasterDB.Table(new(model.User).TableName()).Where("group_id=? and status=0", groupId).Find(&userList)
	return userList, err
}

func (this *UserLogic) FindUserByLevel(levelId int) []*model.User {
	userList := make([]*model.User, 0)
	MasterDB.Table(new(model.User)).Where("user_level = ? and status=0", levelId).Find(&userList)
	return userList
}

func (this *UserLogic) UpdateUser(user *model.User) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	user.UpdateTime = time.Now()
	var err error
	if user.Name == "" {
		_, err = MasterDB.Table(user.TableName()).ID(user.Id).Cols("extra_meta").Update(user)
	} else {
		_, err = MasterDB.Table(user.TableName()).ID(user.Id).Cols("name").Cols("extra_meta").Update(user)
	}
	if err != nil {
		session.Rollback()
		log.Println("update user error:", err)
		session.Commit()
		return err
	}
	session.Commit()
	return nil
}

func (this *UserLogic) FindUserById(id int) bool {
	user := new(model.User)
	has, err := MasterDB.Table(user.TableName()).Where("id = ?", id).Get(user)
	if err != nil {
		log.Println("FindUserById err :", err)
	}
	return has
}

func (this *UserLogic) InsertUser(user *model.User) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	user.UpdateTime = time.Now()
	_, err := MasterDB.Table(user.TableName()).InsertOne(user)
	if err != nil {
		session.Rollback()
		logger.Errorln("InsertUser user error:", err)
		session.Commit()
		return err
	}
	session.Commit()
	return nil

}

func (this *UserLogic) DeleteUser(user *model.User) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	user.UpdateTime = time.Now()
	_, err := MasterDB.Table(new(model.User).TableName()).ID(user.Id).Cols("status").Update(user)
	if err != nil {
		session.Rollback()
		logger.Errorln("DeleteUser user error:", err)
		session.Commit()
		return err
	}
	session.Commit()
	return nil

}

func (this *UserLogic) UpdateStatusByGroupId(groupId int, status int) error {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	user := new(model.User)
	_, err := MasterDB.Table(new(model.User).TableName()).Cols("status").Where("group_id = ?", groupId).Update(user)
	if err != nil {
		session.Rollback()
		logger.Errorln("UpdateStatusByGroupId  error:", err)
		session.Commit()
		return err
	}
	session.Commit()
	return nil
}
