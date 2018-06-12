package controller

import (
"github.com/emicklei/go-restful"
"log"
"eyecool.com/node-retrieval/logic"
"eyecool.com/node-retrieval/model"

)

type AccessController struct {
}

func (this *AccessController)Req(eq *restful.Request, rsp *restful.Response){
	log.Print("Received AccessController.Req API request")
	rsp.WriteEntity(map[string]string{
		"message": "Hi, this is the Req API",
	})
	var repos	[]*model.Repository =logic.DefaultRepository.FindAll();
	//global.G_Dispatcher.Track("cccccccccc")
	for i,v:= range repos{
		log.Printf(" DefaultRepository : %v %v %d",v.Id,v.Name ,i)
	}
}

func (s *AccessController) Anything(req *restful.Request, rsp *restful.Response) {
	log.Print("Received AccessController.Verify API request")
	rsp.WriteEntity(map[string]string{
		"message": "Hi, this is the Verify API",
	})
	rsp.ResponseWriter.Write([]byte("xxxxxxxxResponseWriter Verify xxxxxxxx"))
}



