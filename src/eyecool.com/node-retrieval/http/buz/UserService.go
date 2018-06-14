package buz

import (
	"github.com/satori/go.uuid"
	"strings"
	"fmt"
	"strconv"
	"time"
	"eyecool.com/node-retrieval/logic"
	"eyecool.com/node-retrieval/model"
	"eyecool.com/node-retrieval/utils"
	"log"
)

type UserService struct {
}
type UserResponse struct {
	Rtn       int                `json:"rtn"`                  //接收状态。0表示接收正常，非0表示接收异常（<0表示错误，>0表示警告）
	Message   string             `json:"message,omitempty"`    //接收状态描述
	Id        string             `json:"id,omitempty"`         //用户id
	Name      string             `json:"name,omitempty"`       //用户姓名
	Sets      []*model.UserGroup `json:"sets,omitempty"`       //用户组数组
	Users     []*model.User      `json:"users,omitempty"`      //用户数组
	SessionId string             `json:"session_id,omitempty"` //登录的session_id
	ExtraMeta string             `json:"extra_meta,omitempty"` //额外信息
}

var userLogic = new(logic.UserLogic)
var cacheMap = new(utils.CacheMap)

func (this *UserService) Login(user *model.User) map[string]interface{} {
	result := make(map[string]interface{})
	loginF := 0
	logger := logic.GetLogger(nil)
	if user != nil && user.Name != "" && user.Password != "" {
		existUser, has, err := userLogic.FindUserByName(user.Name)
		if err != nil {
			logger.Errorln("find user error:", err)
			return nil
		}
		if has {
			if existUser.Password == user.Password {
				loginF = 1
				result["user"] = existUser
			} else {
				loginF = 2
			}
		} else {
			loginF = 0
		}
	} else {
		loginF = 3
	}
	result["loginF"] = loginF
	return result
}

func (this *UserService) RespLoginResult(result map[string]interface{}) *UserResponse {
	response := new(UserResponse)
	loginF := result["loginF"]
	switch loginF {
	case 0:
		response.Rtn = -1
		response.Message = "用户不存在!"
	case 1:
		uuid, _ := uuid.NewV4()
		sessionId := strings.Replace(uuid.String(), "-", "", -1)
		existUser, _ := result["user"].(*model.User)
		cacheMap.SetUserSession(sessionId, existUser)
		fmt.Println("----------------- CacheMap:", cacheMap.GetInstance())
		response.Rtn = 0
		response.SessionId = sessionId
		response.Message = "登入成功!"
	case 2:
		response.Rtn = -1
		response.Message = "密码错误!"
	case 3:
		response.Rtn = -1
		response.Message = "参数错误!"
	}
	return response
}

func (this *UserService) GetTopUserAndTopGroup(user *model.User) *UserResponse {
	response := new(UserResponse)
	users := make([]*model.User, 0)
	userGroups := userGroupLogic.FindAllTopUserGroup(new(model.UserGroup))
	if len(userGroups) > 0 {
		for _, v := range userGroups {
			ids := userGroupLogic.FindPredecessorIds(v.ParentId)
			v.PredecessorIds = ids
			userList, err := userLogic.FindUsersByGroupId(v.Id)
			if err != nil {
				log.Println("find user error:", err)
				continue
			}
			if len(userList) > 0 {
				for _, existUser := range userList {
					users = append(users, existUser)
				}
			}
		}
	}
	response.Rtn = 0
	response.Message = "获取成功!"
	response.Sets = userGroups
	response.Users = users
	return response
}

func (this *UserService) GetDepthUserAndUserGroup(idStr, depthStr, allStr string) *UserResponse {
	response := new(UserResponse)
	logger := logic.GetLogger(nil)
	groupId, _, err := utils.GetIdAndClusterId(idStr)
	if err != nil {
		response.Rtn = -1
		response.Message = "id参数不合格!"
		return response
	}
	depth, err2 := strconv.Atoi(depthStr)
	all, err3 := strconv.Atoi(allStr)
	if err2 != nil {
		logger.Errorln("param conv error:", err2)
		response.Rtn = -1
		response.Message = "depth参数不合格!"
		return response
	}
	if err3 != nil {
		logger.Errorln("param conv error:", err3)
		response.Rtn = -1
		response.Message = "all参数不合格!"
		return response
	}
	fillDepthInfoResult(groupId, all, depth, response)
	return response
}
func fillDepthInfoResult(id int, all int, depth int, response *UserResponse) {
	fmt.Println("ididididididididididididididididididid", id)
	has, groupLevel := userGroupLogic.FindGroupLevelById(id)
	if !has {
		response.Rtn = -1
		response.Message = "该数据不存在!"
		return
	}
	switch all {
	case 0:
		userGroups := make([]*model.UserGroup, 0)
		users := make([]*model.User, 0)
		for i := 0; i <= depth; i++ {
			userGroupList := userGroupLogic.FindUserGroupByLevel(groupLevel + i)
			for _, v := range userGroupList {
				userGroups = append(userGroups, v)
			}
			userList := userLogic.FindUserByLevel(groupLevel + i)
			for _, v := range userList {
				users = append(users, v)
			}
		}
		response.Sets = userGroups
		response.Users = users
		response.Rtn = 0
		response.Message = "查询成功!"
	case 1:
		groupLevel = groupLevel + depth
		userGroupList := userGroupLogic.FindUserGroupByLevel(groupLevel)
		userList := userLogic.FindUserByLevel(groupLevel)
		response.Sets = userGroupList
		response.Users = userList
		response.Rtn = 0
		response.Message = "查询成功!"
	default:
		response.Rtn = -1
		response.Message = "参数错误!"
	}
}

func (this *UserService) UpdateUser(user *model.User, response *UserResponse) {
	has := userLogic.FindUserById(user.Id)
	if !has {
		response.Rtn = -1
		response.Message = "该用户不存在!"
		return
	}
	err := userLogic.UpdateUser(user)
	if err != nil {
		log.Println("update user error:", err)
		response.Rtn = -1
		response.Message = "update user error!"
		return
	}
	response.Rtn = 0
	response.Message = "修改成功!"
}

func (this *UserService) InsertUser(user *model.User, response *UserResponse) {
	if user.Name == "" {
		response.Rtn = -1
		response.Message = "name不能为空!"
		return
	}
	if user.Password == "" {
		response.Rtn = -1
		response.Message = "密码不能为空!"
		return
	}
	groupId, clusterId, err := utils.GetIdAndClusterId(user.Predecessor_id)
	if err != nil {
		response.Rtn = -1
		response.Message = "Predecessor_id不合格!"
		return
	}
	_, has, _ := userLogic.FindUserByName(user.Name)
	if has {
		response.Rtn = -1
		response.Message = "该用户已存在!"
		return
	}
	fmt.Println("groupIdgroupIdgroupIdgroupIdgroupIdgroupIdgroupId", groupId)
	has, globalLevel := userGroupLogic.FindGroupLevelById(groupId)
	if has {
		user.Password = utils.MD5(user.Password)
		user.UserLevel = globalLevel
		user.UpdateTime = time.Now()
		//user.ExtraMeta = user.Extra_meta
		user.ClusterId = clusterId
		user.Status = 0
		user.CreateTime = time.Now().Unix()
		err = userLogic.InsertUser(user)
		if err != nil {
			log.Println("insert user error:", err)
			response.Rtn = -1
			response.Message = "insert user error!"
			return
		}
		response.Rtn = 0
		response.Id = strconv.Itoa(user.Id)
		response.Message = "添加成功!"
		return
	}
	response.Rtn = -1
	response.Message = "用户组不存在,添加失败!"
	return

}

func (this *UserService) DeleteUser(idStr string, response *UserResponse) {
	logger := logic.GetLogger(nil)
	userId, _, err := utils.GetIdAndClusterId(idStr)
	if err != nil {
		response.Rtn = -1
		response.Message = "参数不合格!"
		return
	}
	user := new(model.User)
	user.Status = 2
	user.Id = userId
	err = userLogic.DeleteUser(user)
	if err != nil {
		logger.Errorln("delete user error:", err)
		response.Rtn = -1
		response.Message = "delete user error!"
		return
	}
	response.Rtn = 0
	response.Message = "删除成功!"
}
