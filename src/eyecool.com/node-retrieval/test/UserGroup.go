package main

import (
	. "eyecool.com/node-retrieval/db"
	"fmt"
)

type UserGroupLogic struct {
}

func (this *UserGroupLogic) SelectGroupLevelById(id int) int {
	groupLevel := 0
	has, err := MasterDB.Table("buz_user_group").Cols("group_level").Where("id = ?", id).Get(&groupLevel)
	fmt.Println(has, err)
	return groupLevel
}
