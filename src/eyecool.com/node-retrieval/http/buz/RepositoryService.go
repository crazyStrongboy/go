package buz

import (
	"eyecool.com/node-retrieval/logic"
	"eyecool.com/node-retrieval/model"
	"time"
	"strconv"
	"strings"
	"eyecool.com/node-retrieval/utils"
	"eyecool.com/node-retrieval/http/service"
	"fmt"
	"encoding/json"
)

type RepositoryRequest struct {
	Id        string
	Name      string
	ExtraMeta string `json:"extra_meta"`
}

type RepositoryResponse struct {
	Rtn     int             `json:"rtn"`
	Message string          `json:"message"`
	Results []*logic.Result `json:"results"`
}

type InsertRepositoryResponse struct {
	Id      string `json:"id"`
	Rtn     int    `json:"rtn"`
	Message string `json:"message"`
}

var repositoryLogic = new(logic.RepositoryLogic)

func QueryRepository() *RepositoryResponse {
	result := &RepositoryResponse{}
	results, err := logic.DefaultRepository.QueryRepository()
	if err != nil {
		result.Rtn = -1
		result.Message = "查询失败！"
		return result
	}
	result.Results = results
	result.Rtn = 0
	result.Message = "查询成功!"
	return result
}

func InsertRepository(r *RepositoryRequest, user *model.User) *InsertRepositoryResponse {
	result := &InsertRepositoryResponse{}
	flag := logic.DefaultRepository.FindByName(r.Name)
	if flag {
		result.Rtn = -1
		result.Message = "该库已存在"
		return result
	}
	repository := &model.Repository{
		ExtraMeta: r.ExtraMeta,
		//集群号先写死
		ClusterId: 1,
		//CreatorId:user.Id,
		CreatorId:  1,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now(),
		Name:       r.Name,
	}
	err := logic.DefaultRepository.InsertRepository(repository)
	if err != nil {
		result.Rtn = -1
		result.Message = "插入失败"
		return result
	}
	pkId := strconv.Itoa(repository.PkId)
	clusterId := strconv.Itoa(repository.ClusterId)
	ru := &model.Repository{
		PkId: repository.PkId,
		Id:   pkId + "@" + clusterId,
	}
	err = logic.DefaultRepository.UpdateRepository(ru)
	if err != nil {
		result.Rtn = -1
		result.Message = "插入失败"
		return result
	}

	//库同步
	lifecycleRequest := &model.LifecycleRequest{
		RepositoryId: ru.Id,
		Type:         0, //0--增加 1--删除
	}
	response := &model.LifecycleResponse{}
	err = service.RepositoryLifecycle(nil, lifecycleRequest, response)
	if err != nil {
		fmt.Println("InsertRepository synchronized err :", err)
		result.Rtn = -1
		result.Message = "插入失败"
		return result
	}

	result.Id = pkId
	result.Rtn = 0
	result.Message = "插入成功"
	return result
}

func UpdateRepository(r *RepositoryRequest) *model.RespMsg {
	result := &model.RespMsg{}
	if length := strings.Count(r.Name, "") - 1; length > 128 {
		result.Message = "库名不能大于128个字符"
		result.Rtn = -1
		return result
	}
	pkId, _, err := utils.GetIdAndClusterId(r.Id)
	if err != nil {
		result.Rtn = -1
		result.Message = "参数错误!"
		return result
	}
	flag, repository := logic.DefaultRepository.FindByPrimaryKey(pkId)
	if !flag {
		result.Rtn = -1
		result.Message = "要修改的库不存在!"
		return result
	}
	if r.Name != "" {
		repository.Name = r.Name
	}
	repository.ExtraMeta = r.ExtraMeta
	err = logic.DefaultRepository.UpdateRepository(repository)
	if err != nil {
		result.Message = "更新失败!"
		result.Rtn = -1
		return result
	}
	//库同步
	//lifecycleRequest := &model.LifecycleRequest{
	//	RepositoryId: r.Id,
	//	Type:         0, //0--增加 1--删除
	//}
	//response := &model.LifecycleResponse{}
	//err = service.RepositoryLifecycle(nil, lifecycleRequest, response)
	//if err != nil {
	//	fmt.Println("InsertRepository synchronized err :", err)
	//	result.Rtn = -1
	//	result.Message = "更新失败"
	//	return result
	//}

	result.Message = "更新成功"
	result.Rtn = 0
	return result
}

//删除人像库
func DeleteRepository(repositoryId string) *model.RespMsg {
	result := &model.RespMsg{}
	if repositoryId == "" {
		result.Rtn = -1
		result.Message = "repositoryId不能为空!"
		return result
	}
	id, _, err := utils.GetIdAndClusterId(repositoryId)
	if err != nil || id == -2 {
		result.Rtn = -1
		result.Message = "参数错误!"
		return result
	}
	has, _ := logic.DefaultRepository.FindByPrimaryKey(id)
	if !has {
		result.Rtn = -1
		result.Message = "要删除的人像库不存在!"
		return result
	}
	//删除人像库
	err = logic.DefaultRepository.DeleteRepository(id)
	if err != nil {
		result.Rtn = -1
		result.Message = "删除失败"
		return result
	}
	//将与库关联的人状态置为2
	err = peopleLogic.UpdateStatusByRepositoryId(2, repositoryId)
	//将与库关联的特征值状态置为2
	err = featureLogic.UpdateStatusByRepositoryId(2, repositoryId)
	//将与库关联的子任务置为2
	err = logic.TaskChildrenLogic{}.UpdateStatusByRepositoryId(2, repositoryId)
	var taskChildrens []*model.TaskChildren
	taskChildrens, err = logic.TaskChildrenLogic{}.FindTaskChildrenByRepositoryId(repositoryId)
	if err != nil {
		result.Rtn = -1
		result.Message = err.Error()
		return result
	}
	//向布控系统发送请求,type : type 0 : ADD  1 : 册除 2 : 修改....删除这些子任务!!
	if taskChildrens != nil && len(taskChildrens) > 0 {
		for _, taskChild := range taskChildrens {
			param := &Param{
				PkId: taskChild.PkId,
				Type: 1,
			}
			bytes, _ := json.Marshal(param)
			utils.DoBytesPost(task_url, bytes)
		}
	}
	//库同步
	lifecycleRequest := &model.LifecycleRequest{
		RepositoryId: repositoryId,
		Type:         1, //0--增加 1--删除
	}
	response := &model.LifecycleResponse{}
	err = service.RepositoryLifecycle(nil, lifecycleRequest, response)
	if err != nil {
		fmt.Println("InsertRepository synchronized err :", err)
		result.Rtn = -1
		result.Message = "删除失败"
		return result
	}

	result.Rtn = 0
	result.Message = "删除成功"
	return result
}
