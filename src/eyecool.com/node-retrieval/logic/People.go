package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"fmt"
	"errors"
	"log"
)

type PeopleLogic struct {
}

func (this *PeopleLogic) Insert(people *model.People) (*model.People, error) {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(people.TableName()).InsertOne(people)
	session.Commit()
	return people, err
}

func (this *PeopleLogic) FindPeopleById(id int64) (bool, *model.People) {
	people := new(model.People)
	has, _ := MasterDB.Table(people.TableName()).ID(id).Get(people)
	return has, people
}

func (this *PeopleLogic) UpdateById(people *model.People) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(people.TableName()).ID(people.Id).Update(people)
	session.Commit()
	return err
}
func (this *PeopleLogic) DeleteById(people *model.People) error {
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Table(people.TableName()).ID(people.Id).Cols("people_status", "update_time").Update(people)
	session.Commit()
	return err
}
func (peopleLogic *PeopleLogic) FindByRetrieval(retrieval *model.Retrieval) []*model.People {
	peoples := make([]*model.People, 0)
	session := MasterDB.Table(new(model.People).TableName())
	if retrieval.PersonId != "" {
		session.Where("person_id = ?", retrieval.PersonId)
	}
	if retrieval.Name != "" {
		session.Where("name = ?", retrieval.Name)
	}
	if retrieval.RepositoryId != "" {
		session.Where("repository_id = ?", retrieval.RepositoryId)
	}
	if retrieval.Timestamp != "" {
		session.Where("UNIX_TIMESTAMP(create_time) > ?", retrieval.Timestamp)
	}
	err := session.Find(&peoples)
	if err != nil {
		fmt.Println(err)
	}
	return peoples
}
func (peopleLogic *PeopleLogic) UpdateStatusByRepositoryId(status int, repositoryId string) error {
	_, err := MasterDB.Exec("update buz_people set status = ? where repository_id = ?", status, repositoryId)
	if err != nil {
		log.Println("UpdateStatusByRepositoryId buz_people err :",err)
		return errors.New("删除people失败!")
	}
	return nil
}
