package logic

import (
	"eyecool.com/node-retrieval/model"
	"fmt"
	. "eyecool.com/node-retrieval/db"
	"strconv"
	"log"
	"github.com/polaris1119/logger"
)


type TaskLogic struct{}
type TaskChildrenLogic struct{}

var DefaultTask =TaskLogic{}
var DefaultTaskChildren =TaskChildrenLogic{}

type Task struct {
	Name string
	Id string
	Surveillances []*Surveillance
}

type Surveillance struct {
	CameraId string
	RepositoryId string
	Threshold float64
	ExtraMeta string
}





func (TaskLogic) FindTaskTupleMap( ) map[string][]*model.TaskTuple {

	fmt.Println("start FindTaskTupleMap....")
	taskTuples := make([]*model.TaskTuple, 0)
	err := MasterDB.SQL("SELECT t.id TaskId, t.`name` TaskName ,tc.pk_id TaskChildrenPkId,tc.camera_id CameraId,tc.repository_id RepositoryId,tc.threshold Threshold FROM buz_task  t JOIN buz_task_children tc ON t.pk_id=tc.task_pk_id WHERE t.`status`=0").Find(&taskTuples)
	if err != nil {
		log.Println("TaskLogic FindTask  error:", err)
		//fmt.Println("TaskLogic taskTuples error : ",err)
		return nil
	}
	fmt.Println("TaskLogic taskTuples: ",taskTuples)
	result:=make(map[string][]*model.TaskTuple,0)
	for i,v :=range  taskTuples{
		fmt.Println("taskTuple:",*v)
		if  col,ok:=result[v.CameraId];ok{
			result[v.CameraId]=append(col,taskTuples[i])
		}else{
			col :=make([]*model.TaskTuple,0)
			col=append(col, taskTuples[i])
			result[v.CameraId]=col
		}
	}
	log.Println("FindTaskTupleMap : ",result)
	return result
}


func (TaskLogic)InsertTask(task *model.Task)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Insert(task)
	if err != nil {
		session.Rollback()
		logger.Errorln("insert Task error:", err)
		return  err
	}
	session.Commit()
	return nil
}

func (TaskLogic)UpdateTask(task *model.Task)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("pk_id=?",task.PkId).Update(task)
	if err != nil {
		session.Rollback()
		logger.Errorln("Update Task error:", err)
		return  err
	}
	session.Commit()
	return nil
}

//判断子任务是否存在
func (TaskChildrenLogic)SelectTaskChildren(cameraId string  ,repositortId string )(bool,*model.TaskChildren){
	objLog := GetLogger(nil)
	taskChildren:=model.TaskChildren{}
	flag,err:=MasterDB.Where("camera_id=? and repository_id=? and status!=2",cameraId,repositortId).Get(&taskChildren)
	if err != nil {
		objLog.Errorln("TaskChildrenLogic selectTaskChildren PageAd error:", err)
		return false,nil
	}
	return flag,&taskChildren
}

//更新子任务
func (TaskChildrenLogic)UpdateTaskChildren(taskChildren *model.TaskChildren)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("pk_id=?",taskChildren.PkId).Update(taskChildren)
	if err != nil {
		session.Rollback()
		logger.Errorln("Update taskChildren error:", err)
		return  err
	}
	session.Commit()
	return nil
}

func (TaskChildrenLogic)InsertTaskChildren(task *model.TaskChildren)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Insert(task)
	if err != nil {
		session.Rollback()
		logger.Errorln("insert Task error:", err)
		return  err
	}
	session.Commit()
	return nil
}

func (TaskChildrenLogic)DeleteTaskChildren(task *model.TaskChildren)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("task_id=?",task.TaskId).Update(task)
	if err != nil {
		session.Rollback()
		logger.Errorln("DeleteTaskChildren  error:", err)
		return  err
	}
	session.Commit()
	return nil
}

func (TaskChildrenLogic)SelectTaskChildrenByTaskId(taskId string)([]*model.TaskChildren,error){
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	taskChildrens := make([]*model.TaskChildren, 0)
	err := MasterDB.Where("task_id=?",taskId).Find(&taskChildrens)
	if err != nil {
		session.Rollback()
		logger.Errorln("DeleteTaskChildren  error:", err)
		return  nil,err
	}
	session.Commit()
	return taskChildrens,nil
}

func (TaskLogic) QueryTask()([]*Task,error){
	objLog := GetLogger(nil)
	tasks := make([]*model.Task, 0)
	task:=make([]*Task,0)
	err := MasterDB.Where("status=0").Find(&tasks)
	if err != nil {
		objLog.Errorln("AdLogic FindAll PageAd error:", err)
		return nil,err
	}
	for _,v:=range tasks{
		t:=&Task{}
		taskPkId:=strconv.Itoa(v.PkId)
		t.Name=v.Name
		t.Id=taskPkId
		Surveillances:=make([]*Surveillance,0)
		taskChildrens := make([]*model.TaskChildren, 0)
		err := MasterDB.Where("task_id=? and status=0",v.Id).Find(&taskChildrens)
		if err!=nil{
			logger.Errorln("QueryTask  error:", err)
			return  nil,err
		}
		for _,v:=range taskChildrens{
			surveillance:=&Surveillance{}
			surveillance.CameraId=v.CameraId
			surveillance.RepositoryId=v.RepositoryId
			surveillance.Threshold=v.Threshold
			surveillance.ExtraMeta=v.ExtraMeta
			Surveillances=append(Surveillances,surveillance )
		}
		task=append(task, t)

	}
	return task,nil
}

//根据cameraid repositoryId taskId更新子任务
func (TaskChildrenLogic )UpdateTaskChildrenByCameraId(children *model.TaskChildren) error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	var err error
	if children.RepositoryId==""|| children.CameraId==""{
		_, err = MasterDB.Where("task_id=? ",children.TaskId).Update(&children)
	}else{
		_, err = MasterDB.Where("task_id=? and repository_id=? and camera_id=?",children.TaskId,children.RepositoryId,children.CameraId).Update(&children)
	}

	if err != nil {
		session.Rollback()
		logger.Errorln("Update taskChildren error:", err)
		return  err
	}
	session.Commit()
	return nil
}