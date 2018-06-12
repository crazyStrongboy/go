package buz

import (
	"eyecool.com/node-retrieval/logic"
	"eyecool.com/node-retrieval/model"
	"time"
	"strconv"
	"strings"
	"eyecool.com/node-retrieval/utils"
)

type RepositoryRequest struct {
	Id string
	Name string
	ExtraMeta string `json:"extra_meta"`
}

type RepositoryResponse struct {
	Rtn int `json:"rtn"`
	Message string `json:"message"`
	Results []*logic.Result `json:"results"`
}

type InsertRepositoryResponse struct {
	Id string `json:"id"`
	Rtn int `json:"rtn"`
	Message string `json:"message"`
}

var repositoryLogic = new(logic.RepositoryLogic)

func QueryRepository() *RepositoryResponse{
	result:=&RepositoryResponse{}
	results,err:=logic.DefaultRepository.QueryRepository()
	if err!=nil{
		result.Rtn=-1
		result.Message="查询失败！"
		return result
	}
	result.Results=results
	result.Rtn=0
	result.Message="查询成功!"
	return result
}

func InsertRepository(r *RepositoryRequest,user *model.User) *InsertRepositoryResponse{
	result:=&InsertRepositoryResponse{}
	flag:=logic.DefaultRepository.SelectByName(r.Name)
	if flag{
		result.Rtn=-1
		result.Message="该库已存在"
		return result
	}
	repository:=&model.Repository{
		ExtraMeta:r.ExtraMeta,
		//集群号先写死
		ClusterId:1,
		//CreatorId:user.Id,
		CreatorId:1,
		CreateTime:time.Now().Unix(),
		UpdateTime:time.Now(),
		Name:r.Name,
	}
	err:=logic.DefaultRepository.InsertRepository(repository)
	if err!=nil{
		result.Rtn=-1
		result.Message="插入失败"
		return result
	}
	pkId:=strconv.Itoa(repository.PkId)
	clusterId:=strconv.Itoa(repository.ClusterId)
	ru:=&model.Repository{
		PkId:repository.PkId,
		Id:pkId+"@"+clusterId,
	}
	err=logic.DefaultRepository.UpdateRepository(ru)
	if err!=nil{
		result.Rtn=-1
		result.Message="插入失败"
		return result
	}

	//库同步
	/*lifecycleRequest := &model.LifecycleRequest{
		RepositoryId:ru.Id,
		Type:0,//0--增加 1--删除
	}
	response := &model.LifecycleResponse{}
	err = service.RepositoryLifecycle(nil, lifecycleRequest, response)
	if err!=nil{
		fmt.Println("InsertRepository synchronized err :",err)
		result.Rtn=-1
		result.Message="插入失败"
		return result
	}*/


	result.Id=pkId
	result.Rtn=0
	result.Message="插入成功"
	return result
}

func UpdateRepository( r*RepositoryRequest)*model.RespMsg{
	result:=&model.RespMsg{}
	if length:=strings.Count(r.Name,"")-1;length>128{
		result.Message="库名不能大于128个字符"
		result.Rtn=-1
		return result
	}
	flag:=logic.DefaultRepository.SelectByName(r.Name)
	if flag{
		result.Rtn=-1
		result.Message="该库已存在"
		return result
	}
	pkId,_,err:=utils.GetClusterIdAndId(r.Id)
	if err!=nil{
		result.Rtn=-1
		result.Message="参数错误"
		return result
	}
	repository:=&model.Repository{
		PkId:pkId,
		ExtraMeta:r.ExtraMeta,
		Name:r.Name,
	}
	err=logic.DefaultRepository.UpdateRepository(repository)
	if err!=nil{
		result.Message="更新失败"
		result.Rtn=-1
		return result
	}
	//库同步
	/*lifecycleRequest := &model.LifecycleRequest{
		RepositoryId:r.Id,
		Type:0,//0--增加 1--删除
	}
	response := &model.LifecycleResponse{}
	err = service.RepositoryLifecycle(nil, lifecycleRequest, response)
	if err!=nil{
		fmt.Println("InsertRepository synchronized err :",err)
		result.Rtn=-1
		result.Message="更新失败"
		return result
	}*/

	result.Message="更新成功"
	result.Rtn=0
	return result
}

//删除人像库
func DeleteRepository (repositoryId string)*model.RespMsg{
	result:=&model.RespMsg{}
	id,_,err:=utils.GetClusterIdAndId(repositoryId)
	//	clusterId,err:=utils.GetClusterId(regionId)
	if err!=nil{
		result.Rtn=-1
		result.Message="参数错误"
		return result
	}
	//删除人像库
	err=logic.DefaultRepository.DeleteRepository(id)
	if err !=nil{
		result.Rtn=-1
		result.Message="删除失败"
		return result
	}
	//库同步
/*	lifecycleRequest := &model.LifecycleRequest{
		RepositoryId:repositoryId,
		Type:1,//0--增加 1--删除
	}
	response := &model.LifecycleResponse{}
	err = service.RepositoryLifecycle(nil, lifecycleRequest, response)
	if err!=nil{
		fmt.Println("InsertRepository synchronized err :",err)
		result.Rtn=-1
		result.Message="删除失败"
		return result
	}*/

	result.Rtn=0
	result.Message="删除成功"
	return result
}