package controller

import (
	"github.com/emicklei/go-restful"
	"eyecool.com/node-retrieval/model"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"eyecool.com/node-retrieval/utils"
	"strconv"
	"eyecool.com/node-retrieval/http/buz"
)

type UserController struct {
}

var userService = new(buz.UserService)
var cacheMap = new(utils.CacheMap)
//登入接口::POST请求 /login
func (this *UserController) GetSelf(req *restful.Request, rsp *restful.Response) {
	response := new(buz.UserResponse)
	user := new(model.User)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, user)
	if err != nil {
		fmt.Println("GetSelf Unmarshal User err : ", err, ":", user)
		response.Rtn = -1
		response.Message = err.Error()
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		return
	}
	result := userService.Login(user)
	if result == nil {
		response.Rtn = -1
		response.Message = "系统内部错误!"
	} else {
		response = userService.RespLoginResult(result)
	}

	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

//查询用户信息:GET请求 /user/self
func (this *UserController) GetSelfInfo(req *restful.Request, rsp *restful.Response) {
	response := new(buz.UserResponse)
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response.Id = strconv.Itoa(user.Id)
		response.Name = user.Name
		response.ExtraMeta = user.ExtraMeta
		response.Rtn = 0
		response.Message = "查询成功!"
	} else {
		response.Rtn = -1
		response.Message = "用户未登入!"
	}

	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

//获取能看到的最上层用户及用户组:GET请求 /user/top
func (this *UserController) GetTopUserAndTopGroup(req *restful.Request, rsp *restful.Response) {
	response := new(buz.UserResponse)
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response = userService.GetTopUserAndTopGroup(user)
	} else {
		response.Rtn = -1
		response.Message = "用户未登入!"
	}

	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

//获取获取某个用户组下一定深度的节点信息:GET请求 /user
func (this *UserController) GetDepthUserAndTopGroup(req *restful.Request, rsp *restful.Response) {
	response := new(buz.UserResponse)

	params := req.Request.URL.Query()
	id := params.Get("id")
	depth := params.Get("depth")
	all := params.Get("all")
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		response = userService.GetDepthUserAndUserGroup(id, depth, all)
	} else {
		response.Rtn = -1
		response.Message = "用户未登入!"
	}

	rsp.Header().Set("Access-Control-Allow-Origin", "*")
	rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(response)
	rsp.ResponseWriter.Write(responseBytes)
}

//新建或修改用户 /user POST
func (this *UserController) InsertOrUpdateUser(req *restful.Request, rsp *restful.Response) {
	response := new(buz.UserResponse)

	var params map[string]string
	user := new(model.User)
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &params)
	if v, f := params["id"]; f {
		user.Id, err = strconv.Atoi(v)
		if err != nil {
			response.Rtn = -1
			response.Message = "id不正确!"
			rsp.Header().Set("Access-Control-Allow-Origin", "*")
			rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
			rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
			rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
			responseBytes, _ := json.Marshal(response)
			rsp.ResponseWriter.Write(responseBytes)
			return
		}
	}
	user.Name = params["name"]
	user.ExtraMeta = params["extra_meta"]
	user.Predecessor_id = params["predecessor_id"]
	user.Password = params["password"]
	if err != nil {
		fmt.Println("InsertOrUpdateUser Unmarshal User err : ", err)
		response.Rtn = -1
		response.Message = err.Error()
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
		rsp.Header().Set("Access-Control-Max-Age", "1800"); //30 min
		responseBytes, _ := json.Marshal(response)
		rsp.ResponseWriter.Write(responseBytes)
		return
	}
	sessionId := req.HeaderParameter("session_id")
	sessionUser := cacheMap.GetUserSession(sessionId)
	if sessionUser != nil {
		if user.Id != 0 {
			//修改用户
			userService.UpdateUser(user, response)
		} else {
			//新增用户
			userService.InsertUser(user, response)
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

//删除用户  /user  DELETE
func (this *UserController) DeleteUser(req *restful.Request, rsp *restful.Response) {
	response := new(buz.UserResponse)
	params := req.Request.URL.Query()
	idStr := params.Get("id")
	sessionId := req.HeaderParameter("session_id")
	user := cacheMap.GetUserSession(sessionId)
	if user != nil {
		userService.DeleteUser(idStr, response)
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
