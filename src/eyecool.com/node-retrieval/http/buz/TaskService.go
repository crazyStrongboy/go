package buz

import (
	"eyecool.com/node-retrieval/model"
	"fmt"
	"strconv"
	"eyecool.com/node-retrieval/logic"
	"time"
	"eyecool.com/node-retrieval/utils"
	"encoding/json"
)

type TaskRequest struct {
	Requests []Request
	Name string
}

type Request struct {
	CameraId string `json:"camera_id"`
	RepositoryId string `json:"repository_id"`
	Threshold float64 `json:"threshold"`
	ExtraMeta string `json:"extra_meta"`
}

type InsertTaskResponse struct {
	Id string `json:"id"`
	Rtn int `json:"rtn"`
	Message string `json:"message"`
}

type Param struct {
	PkId int
	Type int
}

type TaskResponse struct {
	Tasks []*logic.Task `json:"tasks"`
	Rtn int `json:"rtn"`
	Message string `json:"message"`
}


type TaskUpdateRequest struct {
	Id string
	CameraId string `json:"camera_id"`
	RepositoryId string`json:"repository_id"`
	Name string
	Threshold float64
}






func InsertTask(t *TaskRequest) *InsertTaskResponse{
	result:=&InsertTaskResponse{}
	requests:=t.Requests
	for _,v:=range requests{
		cameraId,_,err:=utils.GetIdAndClusterId(v.CameraId)
		repositoryId,_,err:=utils.GetIdAndClusterId(v.RepositoryId)
		if err !=nil{
			result.Rtn=-1
			result.Message="参数错误"
			return result
		}
		flag:=logic.DefaultRepository.SelectRepositoryById(repositoryId)
		if !flag{
			result.Rtn=-1
			result.Message="库不存在"
			return result
		}
		flag,_=logic.DefaultCamera.SelectCameraById(cameraId)
		if !flag{
			result.Rtn=-1
			result.Message="摄像机不存在"
			return result
		}
	}
	task:=&model.Task{
		CreateTime:time.Now().Unix(),
		UpdateTime:time.Now(),
		Status:0,
		//集群号先写死 1
		ClusterId:1,
		Name:t.Name,
	}
	//入库
	err:=logic.DefaultTask.InsertTask(task)
	if err!=nil{
		result.Rtn=-1
		result.Message="插入失败"
		return result
	}
	pkId:=strconv.Itoa(task.PkId)
	clusterId:=strconv.Itoa(1)
	tu:=&model.Task{
		PkId:task.PkId,
		Id:pkId+"@"+clusterId,
	}
	err=logic.DefaultTask.UpdateTask(tu)
	if err!=nil{
		result.Rtn=-1
		result.Message="插入失败"
		return result
	}

	//插入子任务
	for _,v:=range requests{
		//判断子任务是否存在
		flag,taskChildren:=logic.DefaultTaskChildren.SelectTaskChildren(v.CameraId,v.RepositoryId)
		cameraId,cId,_:=utils.GetIdAndClusterId(v.CameraId)
		if flag{
			//存在
			taskChildren.Status=0
			taskChildren.Threshold=v.Threshold
			err:=logic.DefaultTaskChildren.UpdateTaskChildren(taskChildren)
			if err!=nil{
				result.Rtn=-1
				result.Message="插入失败"
				return result
			}
			param:=&Param{
				PkId:taskChildren.PkId,
				Type:2,
			}
			//向go发送请求
			bytes,_:=json.Marshal(param)
			utils.DoBytesPost("http://192.168.0.192:8091/surveillance/api/task/lifecycle",bytes)
		}else{
			//不存在
			taskChildren:=&model.TaskChildren{}
			falg,camera:=logic.DefaultCamera.SelectCameraById(cameraId)
			if falg{
				taskChildren.CameraIp=camera.Ip
			}
			taskChildren.TaskId=tu.Id
			taskChildren.ClusterId=cId
			taskChildren.CameraId=v.CameraId
			taskChildren.RepositoryId=v.RepositoryId
			taskChildren.ExtraMeta=v.ExtraMeta
			taskChildren.Status=0
			taskChildren.CreateTime=time.Now().Unix()
			taskChildren.UpdateTime=time.Now()
			taskChildren.Threshold=v.Threshold
			err:=logic.DefaultTaskChildren.InsertTaskChildren(taskChildren)
			pkId:=strconv.Itoa(taskChildren.PkId)
			cId:=strconv.Itoa(cId)
			cu:=&model.TaskChildren{
				PkId:taskChildren.PkId,
				Id:pkId+"@"+cId,
			}
			err=logic.DefaultTaskChildren.UpdateTaskChildren(cu)
			if err!=nil{
				fmt.Println("TaskService InserTaskChildren err:",err)
				result.Rtn=-1
				result.Message="插入失败"
				return result
			}
			param:=&Param{
				PkId:taskChildren.PkId,
				Type:0,
			}
			//向go发送请求
			bytes,_:=json.Marshal(param)
			utils.DoBytesPost("http://192.168.0.192:8091/surveillance/api/task/lifecycle",bytes)
		}
	}

	result.Id=pkId+"@"+clusterId
	result.Rtn=0
	result.Message="插入成功"
	return result
}

func DeleteTask(id string)*model.RespMsg{
	result:=&model.RespMsg{}
	taskId,_,err:=utils.GetIdAndClusterId(id)
	if err!=nil{
		result.Rtn=-1
		result.Message="参数错误"
		return result
	}
	task:=&model.Task{
		Status:2,
		PkId:taskId,
		UpdateTime:time.Now(),
	}
	//删除任务库
	err=logic.DefaultTask.UpdateTask(task)
	if err !=nil{
		result.Rtn=-1
		result.Message="删除失败"
		return result
	}
	//删除子任务
	taskChildren:=&model.TaskChildren{
		TaskId:id,
		Status:2,
	}
	err=logic.DefaultTaskChildren.DeleteTaskChildren(taskChildren)

	//查询该taskid下的所有taskChildren
	tc,err:=logic.DefaultTaskChildren.SelectTaskChildrenByTaskId(id)
	if err!=nil{
		fmt.Println("DeleteTask SelectTaskChildrenByTaskId err:",err)
		result.Rtn=-1
		result.Message="删除失败"
		return result
	}
	for _,v:=range tc{
		param:=&Param{
			PkId:v.PkId,
			Type:1,
		}
		//向go发送请求
		bytes,_:=json.Marshal(param)
		utils.DoBytesPost("http://192.168.0.192:8091/surveillance/api/task/lifecycle",bytes)
	}
	result.Rtn=0
	result.Message="删除成功"
	return result
}

func QueryTask()*TaskResponse{
	taskResponse:=&TaskResponse{}
	tasks,err:=logic.DefaultTask.QueryTask()
	if err!=nil{
		taskResponse.Rtn=-1
		taskResponse.Message="查询失败"
		return taskResponse
	}
	taskResponse.Tasks=tasks
	taskResponse.Rtn=0
	taskResponse.Message="查询成功"
	return taskResponse

}

func UpdateTask(request *TaskUpdateRequest)*model.RespMsg{
	result:=&model.RespMsg{}
	pkId,_,err:=utils.GetIdAndClusterId(request.Id)
	if err!=nil{
		result.Rtn=-1
		result.Message="参数错误"
		return result
	}
	//更新任务
	task:=&model.Task{
		PkId:pkId,
		UpdateTime:time.Now(),
		Name:request.Name,
	}
	err=logic.DefaultTask.UpdateTask(task)
	if err !=nil{
		fmt.Println("UpdateTask err :",err)
		result.Rtn=-1
		result.Message="更新失败"
		return result
	}
	//更新子任务
	taskChildren:=&model.TaskChildren{
		TaskId:request.Id,
		UpdateTime:time.Now(),
		Threshold:request.Threshold,
		CameraId:request.CameraId,
		RepositoryId:request.RepositoryId,
	}
	err=logic.DefaultTaskChildren.UpdateTaskChildrenByCameraId(taskChildren)
	if err!=nil{
		fmt.Println("UpdateTaskChildren err :",err)
		result.Rtn=-1
		result.Message="更新失败"
		return result
	}

	//查询该taskid下的所有taskChildren
	tc,err:=logic.DefaultTaskChildren.SelectTaskChildrenByTaskId(request.Id)
	if err!=nil{
		fmt.Println("UpdateTask SelectTaskChildrenByTaskId err:",err)
		result.Rtn=-1
		result.Message="更新失败"
		return result
	}
	for _,v:=range tc{
		param:=&Param{
			PkId:v.PkId,
			Type:2,
		}
		//向go发送请求
		bytes,_:=json.Marshal(param)
		utils.DoBytesPost("http://192.168.0.192:8091/surveillance/api/task/lifecycle",bytes)
	}

	result.Rtn=0
	result.Message="更新成功"
	return result
}



