package utils

import (
	"fmt"
	"sort"
)

type FeaturePairList []*FeaturePair

func (list FeaturePairList) Len() int {
	return len(list)
}

func (list FeaturePairList) Less(i, j int) bool {
	if list[i] == nil {
		return false
	}
	if list[j] == nil {
		return true
	}
	if list[i].Similarity < list[j].Similarity {
		return false
	} else if list[i].Similarity > list[j].Similarity {
		return true
	} else {
		return list[i].Similarity > list[j].Similarity
	}
}

func (list FeaturePairList) Swap(i, j int) {
	var temp *FeaturePair = list[i]
	list[i] = list[j]
	list[j] = temp
}
func (this *FeaturePairList) PrintString() {
	for i, e := range *this {
		fmt.Println(" =====index: ", i, " =======score: ", e.Similarity, "========feature id:", e.Tmpl.FaceImageId)
	}
}

func (this *FeaturePairList) Max() int {
	if len(*this) > 0 {
		return (*this)[0].Similarity
	}
	return 0
}

type FeatureTopNList struct {
	MaxTopN  int
	TopNList FeaturePairList
}

func NewFeatureTopNList(topN int) *FeatureTopNList {
	if topN <= 0 {
		topN = 5
	}
	return &FeatureTopNList{
		MaxTopN:  topN,
		TopNList: make([]*FeaturePair, 0, topN),
	}
}

func (this *FeatureTopNList) Put(pair *FeaturePair) {
	if this.isNotFull() {
		if e := this.isExist(pair); e != -1 {
			fmt.Println("discovery the same feature tmpl :", *pair)
		} else {
			this.TopNList = append(this.TopNList, pair)
			sort.Sort(this.TopNList)
			//fmt.Println("!!!!!!!!! this.TopNList.Len(): ", this.TopNList.Len())
		}
	} else {
		if last := this.TopNList[this.MaxTopN-1]; last != nil && last.Similarity < pair.Similarity {
			if e := this.isExist(pair); e == -1 {
				this.TopNList[this.MaxTopN-1] = pair
				sort.Sort(this.TopNList)
			} else {
				//this.TopNList[this.MaxTopN-1] = pair
				//sort.Sort(this.TopNList)
				fmt.Println("discovery the same feature tmpl :", *pair)
			}
		} else {
			//	fmt.Println("!!!!!!!!! when is full array  last score :", last.Score," > ",pair.Score ," do nothing !!!")
		}
	}
	if this.TopNList.Len() > this.MaxTopN {
		this.TopNList = this.TopNList[0:this.MaxTopN]
		this.TopNList.PrintString()
		fmt.Println("!!!!!!!!! this.TopNList.Len() > this.MaxTopN", this.TopNList.Len(), this.MaxTopN)
	}
}

func (this *FeatureTopNList) isExist(pair *FeaturePair) int {
	for idx, e := range this.TopNList {
		if e != nil && e.Tmpl.FaceImageId == pair.Tmpl.FaceImageId {
			return idx
		}
	}
	return -1
}

func (this *FeatureTopNList) isNotFull() bool {
	if this.TopNList.Len() < this.MaxTopN {
		return true
	}
	for i, e := range this.TopNList {
		if e == nil && i < this.MaxTopN {
			return true
		}
	}
	return false
}
