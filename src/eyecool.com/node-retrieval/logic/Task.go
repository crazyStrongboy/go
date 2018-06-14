package logic

import (
	"eyecool.com/node-retrieval/model"
	"fmt"
	. "eyecool.com/node-retrieval/db"
	"strconv"
	"log"
	"errors"
)

type TaskLogic struct{}
type TaskChildrenLogic struct{}

var DefaultTask = TaskLogic{}
var DefaultTaskChildren = TaskChildrenLogic{}

type Task struct {
	Name          string          `json:"name"`
	Id            string          `json:"id"`
	Surveillances []*Surveillance `json:"surveillances"`
}

type Surveillance struct {
	CameraId     string  `json:"camera_id"`
	RepositoryId string  `json:"repository_id"`
	Threshold    float64 `json:"threshold"`
	ExtraMeta    string  `json:"extra_meta"`
}

func (TaskLogic) FindTaskTupleMap() map[string][]*model.TaskTuple {

	fmt.Println("start FindTaskTupleMap....")
	taskTuples := make([]*model.TaskTuple, 0)
	err := MasterDB.SQL("SELECT t.id TaskId, t.`name` TaskName ,tc.pk_id TaskChildrenPkId,tc.camera_id CameraId,tc.repository_id RepositoryId,tc.threshold Threshold FROM buz_task  t JOIN buz_task_children tc ON t.pk_id=tc.task_id WHERE t.`status`=0").Find(&taskTuples)
	if err != nil {
		log.Println("TaskLogic FindTask  error:", err)
		//fmt.Println("TaskLogic taskTuples error : ",err)
		return nil
	}
	fmt.Println("TaskLogic taskTuples: ", taskTuples)
	result := make(map[string][]*model.TaskTuple, 0)
	for i, v := range taskTuples {
		fmt.Println("taskTuple:", *v)
		if col, ok := result[v.CameraId]; ok {
			result[v.CameraId] = append(col, taskTuples[i])
		} else {
			col := make([]*model.TaskTuple, 0)
			col = append(col, taskTuples[i])
			result[v.CameraId] = col
		}
	}
	log.Println("FindTaskTupleMap : ", result)
	return result
}

func (TaskLogic) FindTaskById(id string) (bool, *model.Task) {
	task := new(model.Task)
	log.Println("insert FindTaskById ,id : ", id)
	has, err := MasterDB.Table(task).Where("id = ?", id).Get(task)
	if err != nil {
		log.Println("FindTaskById err :", err)
	}
	return has, task
}

func (TaskLogic) InsertTask(task *model.Task) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Insert(task)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("insert Task error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (TaskLogic) UpdateTask(task *model.Task) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("pk_id=?", task.PkId).Update(task)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("Update Task error:", err)
		return err
	}
	session.Commit()
	return nil
}

//判断子任务是否存在
func (TaskChildrenLogic) FindTaskChildrenByCameraIpAndRepId(cameraId string, repositortId string) (bool, *model.TaskChildren) {
	taskChildren := model.TaskChildren{}
	flag, err := MasterDB.Where("camera_id=? and repository_id=? and status!=2", cameraId, repositortId).Get(&taskChildren)
	if err != nil {
		log.Println("TaskChildrenLogic selectTaskChildren PageAd error:", err)
		return false, nil
	}
	return flag, &taskChildren
}

//更新子任务
func (TaskChildrenLogic) UpdateTaskChildren(taskChildren *model.TaskChildren) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("pk_id=?", taskChildren.PkId).Update(taskChildren)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("Update taskChildren error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (TaskChildrenLogic) InsertTaskChildren(task *model.TaskChildren) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	//fmt.Println("logic task :", task)
	_, err := MasterDB.Insert(task)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("#################insert Task error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (TaskChildrenLogic) DeleteTaskChildren(task *model.TaskChildren) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("task_id=?", task.TaskId).Update(task)
	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("DeleteTaskChildren  error:", err)
		return err
	}
	session.Commit()
	return nil
}

func (TaskChildrenLogic) FindTaskChildrenByTaskId(taskId string) ([]*model.TaskChildren, error) {
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	taskChildrens := make([]*model.TaskChildren, 0)
	err := MasterDB.Where("task_id=?", taskId).Find(&taskChildrens)
	if err != nil {
		session.Rollback()
		logger.Errorln("DeleteTaskChildren  error:", err)
		return nil, err
	}
	session.Commit()
	return taskChildrens, nil
}

func (TaskLogic) QueryTask() ([]*Task, error) {
	tasks := make([]*model.Task, 0)
	task := make([]*Task, 0)
	err := MasterDB.Where("status=0").Find(&tasks)
	if err != nil {
		log.Println("AdLogic FindAll PageAd error:", err)
		return nil, err
	}
	for _, v := range tasks {
		t := &Task{}
		taskPkId := strconv.Itoa(v.PkId)
		t.Name = v.Name
		t.Id = taskPkId
		Surveillances := make([]*Surveillance, 0)
		taskChildrens := make([]*model.TaskChildren, 0)
		err := MasterDB.Where("task_id=? and status=0", v.Id).Find(&taskChildrens)
		if err != nil {
			log.Println("QueryTask  error:", err)
			return nil, err
		}
		for _, v := range taskChildrens {
			surveillance := &Surveillance{}
			surveillance.CameraId = v.CameraId
			surveillance.RepositoryId = v.RepositoryId
			surveillance.Threshold = v.Threshold
			surveillance.ExtraMeta = v.ExtraMeta
			Surveillances = append(Surveillances, surveillance)
		}
		t.Surveillances = Surveillances
		task = append(task, t)

	}
	return task, nil
}

//根据cameraid repositoryId taskId更新子任务
func (TaskChildrenLogic) UpdateTaskChildrenByCameraId(children *model.TaskChildren) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	var err error
	if children.RepositoryId == "" || children.CameraId == "" {
		_, err = MasterDB.Where("task_id=? ", children.TaskId).Update(children)
	} else {
		_, err = MasterDB.Where("task_id=? and repository_id=? and camera_id=?", children.TaskId, children.RepositoryId, children.CameraId).Update(children)
	}

	if err != nil {
		session.Rollback()
		session.Commit()
		log.Println("UpdateTaskChildrenByCameraId task err : ", err)
		return err
	}
	session.Commit()
	return nil
}

func (TaskChildrenLogic) UpdateStatusByRepositoryId(status int, repositoryId string) error {
	_, err := MasterDB.Exec("update buz_task_children set status = ? where repository_id = ?", status, repositoryId)
	if err != nil {
		log.Println("UpdateStatusByRepositoryId buz_task_children err: ", err)
		return errors.New("删除子任务失败!")
	}
	return nil
}
func (childrenLogic TaskChildrenLogic) FindTaskChildrenByRepositoryId(repositoryId string) ([]*model.TaskChildren, error) {
	taskChildrens := make([]*model.TaskChildren, 0)
	err := MasterDB.Table(new(model.TaskChildren).TableName()).Where("repository_id = ?", repositoryId).Find(&taskChildrens)
	if err != nil {
		log.Println("FindTaskChildrenByRepositoryId buz_task_children err: ", err)
		return nil, errors.New("查询子任务失败!")
	}
	return taskChildrens, nil
}

func (TaskChildrenLogic) UpdateStatusByCameraIdAndRepositoryId(status int, cameraId, repositoryId string) error {
	_, err := MasterDB.Exec("update buz_task_children set status = ? where repository_id = ? and camera_id = ?", status, repositoryId, cameraId)
	if err != nil {
		log.Println("FindTaskChildrenByRepositoryId buz_task_children err: ", err)
		return errors.New("删除子任务失败!")
	}
	return nil
}
