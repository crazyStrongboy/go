package main

import (
	"encoding/json"
	"fmt"
)

func main()  {

	type Rectx struct{
		X int `json:"x"`
		Y int
		T int
		B int
	}

	rect:="[{\"1\":{\"x\":780,\"y\":216,\"t\":972,\"b\":462}}]"

	rmap:=make([]map[int]Rectx,0)

	err := json.Unmarshal([]byte(rect), &rmap)
	if err != nil {
		fmt.Println("Unmarshal LifecycleRequest err : ", err)
	}

	fmt.Println(rmap)

}
