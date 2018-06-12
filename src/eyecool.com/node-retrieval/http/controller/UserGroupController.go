package controller

import (
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"eyecool.com/node-retrieval/model"
	"strconv"
	"eyecool.com/node-retrieval/http/buz"
)

type UserGroupController struct {
}
type UserGroupRequest struct {
	Id             string //修改部门Id
	Predecessor_id string //前继id($编号@$集群号)
	Name           string //部门名字
	Extra_meta     string //额外信息
}

var userGroupService = new(buz.UserGroupService)

//新建或者修改部门::POST请求  /group
func (this *UserGroupController) InsertOrUpdateUserGroup(req *restful.Request, rsp *restful.Response) {
	response := new(buz.UserGroupResponse)
	request := new(UserGroupRequest)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, request)
	if err != nil {
		fmt.Println("InsertOrUpdateUserGroup Unmarshal userGroup err : ", err)
		response.Rtn = -1
		response.Message = err.Error()
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		return
	}
	userGroup := new(model.UserGroup)
	if request.Id != "" {
		userGroup.Id, err = strconv.Atoi(request.Id)
		if err != nil {
			response.Rtn = -1
			response.Message = "id不正确!"
			responseBytes, _ := json.Marshal(response)
			rsp.ResponseWriter.Write(responseBytes)
			return
		}
	}
	userGroup.Predecessor_id = request.Predecessor_id
	userGroup.Name = request.Name
	userGroup.ExtraMeta = request.Extra_meta
	sessionId := req.HeaderParameter("session_id")
	sessionUser := cacheMap.GetUserSession(sessionId)
	if sessionUser != nil {
		if userGroup.Id != 0 {
			userGroupService.UpdateUserGroup(userGroup, response)
		} else {
			userGroupService.InsertUserGroup(userGroup, response)
		}
	} else {
		response.Rtn = -1
		response.Message = "用户未登录!"
	}

	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

//删除部门
func (this *UserGroupController) DeleteUserGroup(req *restful.Request, rsp *restful.Response) {
	response := new(buz.UserGroupResponse)
	params := req.Request.URL.Query()
	idStr := params.Get("id")
	sessionId := req.HeaderParameter("session_id")
	sessionUser := cacheMap.GetUserSession(sessionId)
	if sessionUser != nil {
		response = userGroupService.DeleteUserGroup(idStr)
	} else {
		response.Rtn = -1
		response.Message = "用户未登录!"
	}

	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}
