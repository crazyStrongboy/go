package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"fmt"
)

type RetrievalLogic struct{}

var DefaultRetrieval = RetrievalLogic{}

func (self RetrievalLogic) UpdateRetrievalResults(retrievalId int, results string) {
	go func() {
		MasterDB.Table(new(model.Retrieval)).Where("id=?", retrievalId).Update(map[string]interface{}{
			"results": results})
	}()
}

func (self *RetrievalLogic) UpdateRetrievalById(retrieval *model.Retrieval) {
	MasterDB.Table(new(model.Retrieval)).ID(retrieval.Id).Update(retrieval)
}

func (self *RetrievalLogic) SelectRetrievalById(retrievalId int64) (bool, *model.Retrieval) {
	retrieval := new(model.Retrieval)
	has, _ := MasterDB.Table(new(model.Retrieval)).ID(retrievalId).Get(retrieval)
	return has, retrieval
}
func (retrievalLogic *RetrievalLogic) Insert(retrieval *model.Retrieval) {
	_, err := MasterDB.Table(new(model.Retrieval)).Insert(retrieval)
	if err!=nil {
		fmt.Println(err)
	}
}
