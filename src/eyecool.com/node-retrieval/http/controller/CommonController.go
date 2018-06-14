package controller

import "github.com/emicklei/go-restful"

func SetResponse(res *restful.Response) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	res.Header().Set("Access-Control-Max-Age", "1800"); //30 min
}
