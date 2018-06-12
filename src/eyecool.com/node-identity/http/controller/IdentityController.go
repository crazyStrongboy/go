package controller

import (
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"encoding/json"
	"log"
)

type IdentityController struct {
}
type IdCardRequest struct {
	IdNum string `json:"idNum,omitempty"`
}

func (this *IdentityController) UploadIdCard(req *restful.Request, rsp *restful.Response) {
	oi := IdCardRequest{}
	body, _ := ioutil.ReadAll(req.Request.Body)
	err := json.Unmarshal(body, &oi)
	log.Println(err, oi)
}
