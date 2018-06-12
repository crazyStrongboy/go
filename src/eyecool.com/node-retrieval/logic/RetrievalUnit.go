package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
)

type RetrievalUnitLogic struct{}

var DefaultRetrievalUnit = RetrievalUnitLogic{}

func (self RetrievalUnitLogic) Insert(unit *model.RetrievalUnit) error {
	logger := GetLogger(nil)
	_, err := MasterDB.Insert(unit)
	if err != nil {
		logger.Errorln("insert RetrievalUnit error:", err)
		return err
	}
	return nil
}

func (self RetrievalUnitLogic) UpdateRetrievalUnitResults(unitId, dealNum int, results string) {
	go func() {
		MasterDB.Table(new(model.RetrievalUnit)).Where("id=?", unitId).Update(map[string]interface{}{
			"results":  results,
			"deal_num": dealNum,
		})
	}()
}
func (unitLogic *RetrievalUnitLogic) FindByRetrievalId(retrievalId int64) []*model.RetrievalUnit {
	retrievalUnits := make([]*model.RetrievalUnit, 0)
	MasterDB.Table(new(model.RetrievalUnit)).Where("retrieval_id=?", retrievalId).Find(&retrievalUnits)
	return retrievalUnits
}
