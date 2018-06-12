package buz

import (
	"eyecool.com/node-retrieval/model"
	"time"
	"github.com/polaris1119/logger"
	"eyecool.com/node-retrieval/utils"
	"strconv"
	"fmt"
	"eyecool.com/node-retrieval/logic"
)

type UserGroupService struct {
}
type UserGroupResponse struct {
	Rtn     int    `json:"rtn"`               //错误码
	Message string `json:"message,omitempty"` //错误消息
	Id      string `json:"id,omitempty"`      //创建组的Id
}

var userGroupLogic = new(logic.UserGroupLogic)

func (this *UserGroupService) UpdateUserGroup(group *model.UserGroup, response *UserGroupResponse) {
	group.UpdateTime = time.Now()
	err := userGroupLogic.UpdateUserGroup(group)
	if err != nil {
		logger.Errorln("update userGroup error:", err)
		response.Rtn = -1
		response.Message = "update userGroup error!"
		return
	}
	response.Rtn = 0
	response.Message = "修改成功!"
}

func (this *UserGroupService) InsertUserGroup(group *model.UserGroup, response *UserGroupResponse) {
	predecessor_id := group.Predecessor_id
	parentId, clusterId, err := utils.GetIdAndClusterId(predecessor_id)
	if err != nil {
		response.Rtn = -1
		response.Message = "Predecessor_id不合格!"
		return
	}
	if parentId == -2 {
		group.GroupLevel = 0
	} else {
		existGroup := userGroupLogic.SelectGroupById(parentId)
		group.GroupLevel = existGroup.GroupLevel + 1
		group.ParentId = parentId
	}
	group.ClusterId = clusterId
	group.CreateTime = time.Now().Unix()
	group.UpdateTime = time.Now()
	group.RepositoryId = group.Predecessor_id
	group.Status = 0
	err = userGroupLogic.InsertGroup(group)
	if err != nil {
		response.Rtn = -1
		response.Message = "插入失败!"
		return
	}
	response.Rtn = 0
	response.Message = "插入成功!"
	response.Id = strconv.Itoa(group.Id)
}

func (this *UserGroupService) DeleteUserGroup(param string) *UserGroupResponse {
	response := new(UserGroupResponse)
	parentId, _, err := utils.GetIdAndClusterId(param)
	if err != nil {
		response.Rtn = -1
		response.Message = "参数错误!"
		return response
	}
	err = deleteUserGroupAndUser(parentId)
	if err != nil {
		fmt.Println("DeleteUserGroup:", err)
		response.Rtn = -1
		response.Message = "参数错误!"
		return response
	}
	response.Rtn = 0
	response.Message = "删除成功!"
	return response
}
func deleteUserGroupAndUser(parentId int) error {
	groups, err := userGroupLogic.SelectByParentId(parentId)
	if err != nil {
		return err
	}
	if len(groups) > 0 {
		for _, v := range groups {
			groupId := v.Id
			//2--表示删除 删除组
			userGroupLogic.UpdateStatus(groupId, 2)
			//2--表示删除 删除用户
			userLogic.UpdateStatusByGroupId(groupId, 2)
			deleteUserGroupAndUser(groupId)
		}
	}
	return nil
}
