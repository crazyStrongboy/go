package buz

import (
	"eyecool.com/node-retrieval/model"
	"eyecool.com/node-retrieval/logic"
	"time"
	"strconv"
	"eyecool.com/node-retrieval/utils"
)
//查询摄像机返回对象
type CameraResponse struct {
	Id string `json:"id,omitempty"`
	Cameras []*logic.Cameras `json:"cameras,omitempty"`
	Rtn int `json:"rtn"`
	Message string `json:"message"`
}
//插入返回
type InsertCameraResponse struct {
	Id string `json:"id"`
	Rtn int `json:"rtn"`
	Message string `json:"message"`
}
//请求对象
type CameraRequest struct {
	Id string 	`json:"id"`
	Name string	`json:"name"`
	Url string	`json:"url"`
	Ip string	`json:"ip"`
	Enabled int `json:"enabled"`
	PredecessorId string `json:"predecessor_id"`
	RecParams string	`json:"rec_params"`
	ExtraMeta string `json:"extra_meta"`
}

//查询摄像机
func CameraQuery()*CameraResponse{
	cameras:=logic.DefaultCamera.CameraQuery()
	result:=&CameraResponse{}
	result.Cameras=cameras
	result.Rtn=0
	result.Message="查询成功！"
	return result
}
//插入摄像机
func InsertCamera(camera *CameraRequest,user *model.User)*InsertCameraResponse{
	result:=&InsertCameraResponse{}
	flag:=logic.DefaultCamera.FindIP(camera.Ip)
	if !flag{
		regionId,_:=strconv.Atoi(camera.PredecessorId)
		c:=&model.Camera{
			PredecessorId:camera.PredecessorId,
			RegionId: regionId,
			ExtraMeta:camera.ExtraMeta,
			Ip:camera.Ip,
			//集群号写死 1
			ClusterId:1,
			CreateTime:time.Now(),
			UpdateTime:time.Now(),
			CreatorId:1,
			Name:camera.Name,
			Url:camera.Url,
		}
		//插入camera表
		err:=logic.DefaultCamera.Insert(c)
		if err!=nil{
			result.Rtn=-1
			result.Message="新增失败！"
			return result
		}
		id:=strconv.Itoa(c.PkId)
		cu:=&model.Camera{
			PkId:c.PkId,
			Id:id+"@"+"1",
		}
		//更新camera中id
		err=logic.DefaultCamera.Update(cu)
		if err!=nil{
			result.Rtn=-1
			result.Message="新增失败！"
			return result
		}

		vc:=&model.VideoCamera{
			CameraId:c.Id,
			Status:c.Status,
			Param1:c.Ip,
			RtspUrl:c.Url,
		}
		//新增video_camera表
		err=logic.DefaultCamera.InsertVideoCamera(vc)
		if err!=nil{
			result.Rtn=-1
			result.Message="新增失败！"
			return result
		}

		//	result=c.PkId
		result.Rtn=0
		result.Message="新增成功！"
		return result
	}else{
		result.Rtn=-1
		result.Message="Ip已经存在！"
		return result
	}

}
//删除摄像机
func DeleteCamera(id string)*model.RespMsg{
	result:=&model.RespMsg{}
	cameraId,_,err:=utils.GetIdAndClusterId(id)
	if err!=nil{
		result.Rtn=-1
		result.Message="参数错误！"
		return result
	}
	err=logic.DefaultCamera.DeleteCamera(cameraId)
	if err!=nil{
		result.Rtn=-1
		result.Message="删除失败！"
		return result
	}
	result.Rtn=0
	result.Message="删除成功！"
	return result

}



//更新摄像机
func UpdateCamera(camera *CameraRequest) *model.RespMsg{
	result:=&model.RespMsg{}
	id,_,err:=utils.GetIdAndClusterId(camera.Id)
	//clusterId,err:=utils.GetClusterId(camera.Id)
	if err!=nil{
		result.Rtn=-1
		result.Message="参数错误！"
		return result
	}
	c:=&model.Camera{
		PkId:id,
		Status:camera.Enabled,
		RecParams:camera.RecParams,
		Name:camera.Name,
		Url:camera.Url,
		ExtraMeta:camera.ExtraMeta,
	}
	//更新camera表
	err=logic.DefaultCamera.Update(c)
	if err!=nil{
		result.Rtn=-1
		result.Message="更新失败！"
		return result
	}
	//更新video_camera表
	vc:=&model.VideoCamera{
		CameraId:camera.Id,
		Status:c.Status,
	}
	err=logic.DefaultCamera.UpdateVideoCamera(vc)
	if err!=nil{
		result.Rtn=-1
		result.Message="更新失败！"
		return result
	}

	result.Rtn=0
	result.Message="更新成功！"
	return result
}
