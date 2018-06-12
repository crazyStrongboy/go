package controller

import (
	"github.com/emicklei/go-restful"
	"eyecool.com/node-retrieval/utils"
	"eyecool.com/node-retrieval/http/buz"
	"io/ioutil"
	"encoding/json"
	"log"
	"fmt"
)

type VerifyController struct{}

func (this *VerifyController) FaceVerify(req *restful.Request, res *restful.Response) {
	log.Print("Received VideoController.QueryVideo API request : ", req.Request.RemoteAddr)
	sessionId := req.HeaderParameter("session_id")
	cacheMap := utils.CacheMap{}
	//判断用户是否登陆
	flag := cacheMap.CheckSession(sessionId)
	flag = true
	result := &buz.VerifyReponse{}
	if flag {
		verify := buz.VerifyRequest{}
		body, _ := ioutil.ReadAll(req.Request.Body)
		err := json.Unmarshal(body, &verify)
		if err != nil {
			log.Println("FaceVerify err:", err)

			result.Rtn = -1
			result.Message = "参数错误！"

		}
		result = buz.FaceVerify(&verify)
	} else {
		result.Rtn = -1
		result.Message = "用户未登录"
	}

	fmt.Println(req.Request.Method)
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
	res.Header().Set("Access-Control-Allow-Headers", "x-requested-with");
	res.Header().Set("Access-Control-Max-Age", "1800"); //30 min
	responseBytes, _ := json.Marshal(result)
	res.ResponseWriter.Write(responseBytes)
}
