package http

import (
	"eyecool.com/node-retrieval/http/controller"
	"github.com/emicklei/go-restful"
	. "github.com/polaris1119/config"
	"net/http"
	"log"
	"fmt"
)

func StartWebService() {

	// Create RESTful handler
	accessctl := new(controller.AccessController)
	facectl := new(controller.FaceController)
	retrievalctl := new(controller.RetrievalController)
	ws := new(restful.WebService)
	wc := restful.NewContainer()
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/access")
	ws.Route(ws.GET("/").To(accessctl.Anything))
	ws.Route(ws.GET("/req").To(accessctl.Req))
	wc.Add(ws)

	//wc.Filter(container_filter_A)

	ws1 := new(restful.WebService)
	ws1.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws1.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws1.Path("/business/api")
	ws1.Route(ws1.GET("/").To(accessctl.Anything))
	//ws1.Route(ws1.POST("/face/verify").To(facectl.Verify))
	ws1.Route(ws1.POST("/face/insert_orig_image").To(facectl.InsertOrigImage))
	//ws1.Route(ws1.POST("/retrieval").To(retrievalctl.Retrieval))
	ws1.Route(ws1.POST("/retrieval/repository_feature_insert").To(retrievalctl.RepositoryFeatureInsert))
	ws1.Route(ws1.POST("/retrieval/repository_lifecycle").To(retrievalctl.RepositoryLifecycle))
	ws1.Route(ws1.POST("/retrieval/repository_target").To(retrievalctl.RetrievalRepositoryTarget))
	ws1.Route(ws1.POST("/retrieval/camera_target").To(retrievalctl.RetrievalCameraTarget))

	userController := new(controller.UserController)
	userGroupController := new(controller.UserGroupController)
	clusterController := new(controller.ClusterController)
	peopleController := new(controller.PeopleController)
	origController := new(controller.OrigImageController)
	retrievalController := new(controller.RetrievalController)
	imageFailController := new(controller.ImageFailController)
	alarmController := new(controller.AlarmController)
	ws1.Route(ws1.POST("/login").To(userController.GetSelf))
	ws1.Route(ws1.GET("/user/self").To(userController.GetSelfInfo))
	ws1.Route(ws1.GET("/user/top").To(userController.GetTopUserAndTopGroup))
	ws1.Route(ws1.GET("/user").To(userController.GetDepthUserAndTopGroup))
	ws1.Route(ws1.POST("/user").To(userController.InsertOrUpdateUser))
	ws1.Route(ws1.DELETE("/user").To(userController.DeleteUser))
	ws1.Route(ws1.POST("/group").To(userGroupController.InsertOrUpdateUserGroup))
	ws1.Route(ws1.DELETE("/group").To(userGroupController.DeleteUserGroup))
	ws1.Route(ws1.GET("/self/cluster_id").To(clusterController.GetSelfClusterId))
	ws1.Route(ws1.GET("/cluster_ids").To(clusterController.GetClusterIds))
	ws1.Route(ws1.POST("/repository/picture/synchronized").To(peopleController.PictureSynchronized))
	ws1.Route(ws1.POST("/face/update").To(peopleController.FaceUpdate))
	ws1.Route(ws1.POST("/face/delete").To(peopleController.FaceDelete))
	ws1.Route(ws1.POST("/capture/fetch").To(origController.GetCaptureImage))
	ws1.Route(ws1.POST("/single/image").To(origController.GetSingleImage))
	ws1.Route(ws1.POST("/retrieval").To(retrievalController.PictureSynchronized))
	ws1.Route(ws1.POST("/condition/query").To(retrievalController.ConditionQuery))
	ws1.Route(ws1.POST("/repository/picture/failed").To(imageFailController.GetFailImage))
	ws1.Route(ws1.POST("/hit/alert").To(alarmController.HitAlert))

	camera := new(controller.CameraController)
	region := new(controller.RegionController)
	video := new(controller.VideoController)
	task := new(controller.TaskController)
	verify := new(controller.VerifyController)
	repository := new(controller.RepositoryController)
	ws1.Route(ws1.POST("/face/verify").To(verify.FaceVerify))

	ws1.Route(ws1.GET("/camera").To(camera.CameraQuery))
	ws1.Route(ws1.POST("/camera").To(camera.InsertCamera))
	ws1.Route(ws1.DELETE("/camera").To(camera.DeleteCamera))
	ws1.Route(ws1.PUT("/camera").To(camera.UpdateCamera))

	ws1.Route(ws1.POST("/camera/region").To(region.InsertRegion))
	ws1.Route(ws1.GET("/camera/region").To(region.QueryRegion))
	ws1.Route(ws1.DELETE("/camera/region").To(region.DeleteRegion))
	ws1.Route(ws1.PUT("/camera/region").To(region.UpdateRegion))

	ws1.Route(ws1.GET("/repository").To(repository.QueryRepository))
	ws1.Route(ws1.POST("/repository").To(repository.InsertRepository))
	ws1.Route(ws1.PUT("/repository").To(repository.UpdateRepository))
	ws1.Route(ws1.DELETE("/repository").To(repository.DeleteRepository))

	ws1.Route(ws1.GET("/video").To(video.QueryVideo))
	ws1.Route(ws1.POST("/video").To(video.InsertVideo))
	ws1.Route(ws1.DELETE("/video").To(video.DeleteVideo))
	ws1.Route(ws1.PUT("/video").To(video.UpdateVideo))

	ws1.Route(ws1.GET("/surveillance/task").To(task.QueryTask))
	ws1.Route(ws1.POST("/surveillance/task").To(task.InsertTask))
	ws1.Route(ws1.DELETE("/surveillance/task").To(task.DeleteTask))
	ws1.Route(ws1.PUT("/surveillance/task").To(task.UpdateTask))
	ws1.Route(ws1.POST("/surveillance/task/children/delete").To(task.DeleteChildTask))
	wc.Add(ws1)

	host, err := ConfigFile.GetValue("listen", "host")
	port, err := ConfigFile.GetValue("listen", "port")
	if err != nil {
		log.Fatalf(" env.ini not found http_listen_addr ")
	}
	fmt.Printf(" ListenAndServe %s:%s start !!! \n", host, port)
	// Run server
	log.Fatal(http.ListenAndServe(host+":"+port, wc))

}

func container_filter_A(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("addr %s url path:%v\n", req.Request.RemoteAddr, req.Request.URL)
	trace("container_filter_A: before", 1)
	chain.ProcessFilter(req, resp)
	trace("container_filter_A: after", -1)
}

var indentLevel int

func trace(what string, delta int) {
	indented := what
	if delta < 0 {
		indentLevel += delta
	}
	for t := 0; t < indentLevel; t++ {
		indented = "." + indented
	}
	log.Printf("%s", indented)
	if delta > 0 {
		indentLevel += delta
	}
}
