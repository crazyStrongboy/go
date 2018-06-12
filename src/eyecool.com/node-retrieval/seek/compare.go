package seek

import (
	"eyecool.com/node-retrieval/utils"
	"eyecool.com/node-retrieval/model"
	"eyecool.com/node-retrieval/global"
	"time"
	"log"
)

func CompareTarget(repositoryId string, feature []byte, topk int) (utils.FeaturePairList, error) {
	s := time.Now().UnixNano()
	compareNum, scores := global.G_ChlFaceX.ChlFaceSdkListCompare(global.BuildReposKey(repositoryId), 0, feature, 0, -1)
	log.Printf("=====repos id : [%s]===========scores : [%v]====", global.BuildReposKey(repositoryId), scores)
	featureTopNList := utils.NewFeatureTopNList(topk)
	maxMatchScore, scoresLen := 0, len(scores)
	if setcoll, ok := global.G_ReposHandleMap[global.BuildReposKey(repositoryId)]; ok {
		log.Printf("x ok : [%t]  set [%v]", ok, setcoll)
		setcoll.Each(func(e interface{}) bool {
			entry := e.(*model.FeatureEntry)
			log.Println("entry : ", entry)
			if (entry.Status == 0) && entry.Pos < scoresLen {
				score := int(scores[entry.Pos])
				if score > maxMatchScore {
					maxMatchScore = score
				}
				featureTopNList.Put(&utils.FeaturePair{Similarity: int(score), Tmpl: entry})
			}
			return true
		})
	} else {
		log.Printf("ok : [%t]  set [%v]", ok, setcoll)
	}
	//log.Info("compareTarget result : ",maxMatchScore,featureTopNList.TopNList[0].Tmpl.PersonId,featureTopNList.TopNList[0].Tmpl.Pos)
	//	log.Info("compareTarget scores:",scores)
	spendTime := (time.Now().UnixNano() - s) / 1000000
	log.Printf("====================ChlFaceSdkListCompare spend time: %d  compareNum: %d scores.len: %d: maxMatchScore：%d \n", spendTime, compareNum, len(scores), maxMatchScore)
	return featureTopNList.TopNList, nil
}

