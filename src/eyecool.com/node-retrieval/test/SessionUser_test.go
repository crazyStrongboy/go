package main

import (
	"testing"
	"fmt"
)

//var cache = new(utils.CacheMap)
func TestGetEngine(t *testing.T) {
	/*user1:=new(model.User)
	user1.Name="aa"
	user1.ClusterId=5
	cache.SetUserSession("12",user1)
	user2:=new(model.User)
	user2.Name="bb"
	cache.SetUserSession("13",user2)
	cacheMap := cache.GetInstance()
	fmt.Println(cacheMap)
	a:=cache.CheckSession("14")
	fmt.Println(a)
	cache.ClearSession("13")
	fmt.Println(cacheMap)
	user :=cache.GetUserSession("12")
	fmt.Println(user)
	clusteId := cache.GetUserClusterId("12")
	fmt.Println(clusteId)
	cache.ClearSession("12")
	isEmpty := cache.CheckMap()
	fmt.Println(isEmpty)

	fmt.Println(uuid.NewV4())*/
}

func TestGetIdAndClusterId(t *testing.T) {
	//fmt.Println(utils.MD5("aaa"))
	//a, b, err := utils.GetIdAndClusterId("@13")
	//fmt.Println(a, b,err)
	userGroupLogic := new(UserGroupLogic)
	a := userGroupLogic.SelectGroupLevelById(5)
	fmt.Println(a)
}
