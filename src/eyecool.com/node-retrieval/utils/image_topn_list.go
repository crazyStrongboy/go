package utils

import (
	"eyecool.com/node-retrieval/model"
	"fmt"
	"sort"
)

type   OrigImagePair struct {
	Similarity int
	Tmpl  *model.OrigImage
}

type OrigImagePairList []*OrigImagePair

func (list OrigImagePairList) Len() int {
	return len(list)
}

func (list OrigImagePairList) Less(i, j int) bool {
	if list[i]==nil  {
		return false
	}
	if  list[j]==nil{
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

func (list OrigImagePairList) Swap(i, j int) {
	var temp *OrigImagePair = list[i]
	list[i] = list[j]
	list[j] = temp
}
func (this *OrigImagePairList) PrintString()  {
	for i, e := range *this {
		fmt.Println(" =====index: ",i," =======score: ",e.Similarity ,"========feature id:" ,e.Tmpl.Id)
	}
}


type OrigImageTopNList struct {
	MaxTopN  int
	TopNList OrigImagePairList
}

func NewOrigImageTopNList(topN int) *OrigImageTopNList {
	if topN<=0 {
		topN=5
	}
	return &OrigImageTopNList{
		MaxTopN:  topN,
		TopNList: make([]*OrigImagePair, 0,topN),
	}
}

func (this *OrigImageTopNList) Put(pair *OrigImagePair) {
	if this.isNotFull() {
		if e := this.isExist(pair); e != -1 {
			fmt.Println("discovery the same feature tmpl :", *pair)
		} else {
			this.TopNList = append(this.TopNList, pair)
			sort.Sort(this.TopNList)
			fmt.Println("!!!!!!!!! this.TopNList.Len(): ",this.TopNList.Len())
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
		this.TopNList=this.TopNList[0:this.MaxTopN]
		this.TopNList.PrintString()
		fmt.Println("!!!!!!!!! this.TopNList.Len() > this.MaxTopN", this.TopNList.Len(), this.MaxTopN)
	}
}

func (this *OrigImageTopNList) isExist(pair *OrigImagePair) int {
	for idx, e := range this.TopNList {
		if e != nil && e.Tmpl.Id == pair.Tmpl.Id {
			return idx
		}
	}
	return -1
}

func (this *OrigImageTopNList) isNotFull() bool {
	if  this.TopNList.Len()<this.MaxTopN{
		return true
	}
	for i, e := range this.TopNList {
		if e == nil && i < this.MaxTopN {
			return true
		}
	}
	return false
}
