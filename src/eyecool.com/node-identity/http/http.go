package http

import (
	"github.com/emicklei/go-restful"
	. "github.com/polaris1119/config"
	"net/http"
	"log"
	"fmt"
	"eyecool.com/node-identity/http/controller"
)

func StartWebService() {

	wc := restful.NewContainer()
	identityController := new(controller.IdentityController)
	heartController := new(controller.HeartController)
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/identity")
	ws.Route(ws.POST("/verification").To(identityController.UploadIdCard))
	ws.Route(ws.POST("/heartbeat").To(heartController.HeartBeat))
	wc.Add(ws)

	host, err := ConfigFile.GetValue("listen", "host")
	port, err := ConfigFile.GetValue("listen", "port")
	if err != nil {
		log.Fatalf(" env.ini not found http_listen_addr ")
	}
	fmt.Printf(" ListenAndServe %s:%s start !!! \n", host, port)
	// Run server
	log.Fatal(http.ListenAndServe(host+":"+port, wc))

}
