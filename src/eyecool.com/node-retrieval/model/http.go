package model

type OrigImageRequest struct {
	Uuid             string `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
	CameraId         string `protobuf:"bytes,2,opt,name=cameraId" json:"cameraId,omitempty"`
	CameraIp         string
	ImageName        string `protobuf:"bytes,3,opt,name=imageName" json:"imageName,omitempty"`
	ImageRealPath    string `protobuf:"bytes,4,opt,name=imageRealPath" json:"imageRealPath,omitempty"`
	FaceNum          int32  `protobuf:"varint,5,opt,name=faceNum" json:"faceNum,omitempty"`
	FeatList         string `protobuf:"bytes,6,opt,name=featList" json:"featList,omitempty"`
	FaceRect         string `protobuf:"bytes,7,opt,name=FaceRect" json:"FaceRect,omitempty"`
	FaceProp         string `protobuf:"bytes,7,opt,name=FaceProp" json:"FaceProp,omitempty"`
	Timestamp        int64  `protobuf:"varint,8,opt,name=Timestamp" json:"Timestamp,omitempty"`
	ImageContextPath string `protobuf:"bytes,9,opt,name=imageContextPath" json:"imageContextPath,omitempty"`
	FaceImageUri     string `protobuf:"bytes,10,opt,name=faceImageUri" json:"faceImageUri,omitempty"`
}

type RetrievalRequest struct {
	Uuid          string  `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
	RetrievalId   int32   `protobuf:"varint,2,opt,name=retrievalId" json:"retrievalId,omitempty"`
	FaceImageId   string  `protobuf:"bytes,3,opt,name=faceImageId" json:"faceImageId,omitempty"`
	CameraIds     string  `protobuf:"bytes,4,opt,name=cameraIds" json:"cameraIds,omitempty"`
	RepositoryIds string  `protobuf:"bytes,5,opt,name=repositoryIds" json:"repositoryIds,omitempty"`
	VideoIds      string  `protobuf:"bytes,6,opt,name=videoIds" json:"videoIds,omitempty"`
	Threshold     float64 `protobuf:"fixed64,7,opt,name=threshold" json:"threshold,omitempty"`
	Topk          int32   `protobuf:"varint,8,opt,name=topk" json:"topk,omitempty"`
	Async         bool    `protobuf:"varint,9,opt,name=async" json:"async,omitempty"`
	Params        string  `protobuf:"bytes,10,opt,name=params" json:"params,omitempty"`
	Feats         []byte  `protobuf:"bytes,11,opt,name=feats,proto3" json:"feats,omitempty"`
}

type RetrievalResponse struct {
	Msg         string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	RetrievalId string `protobuf:"bytes,2,opt,name=retrievalId" json:"retrievalId,omitempty"`
	Total       int32  `protobuf:"varint,3,opt,name=total" json:"total,omitempty"`
}

type RetrievalFeatureRequest struct {
	Uuid         string `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
	PeopleId     string `protobuf:"bytes,2,opt,name=peopleId" json:"peopleId,omitempty"`
	CameraId     string `protobuf:"bytes,3,opt,name=cameraId" json:"cameraId,omitempty"`
	RepositoryId string `protobuf:"bytes,4,opt,name=repositoryId" json:"repositoryId,omitempty"`
	Id           string `protobuf:"bytes,5,opt,name=id" json:"id,omitempty"`
	Type         int32  `protobuf:"varint,6,opt,name=type" json:"type,omitempty"`
	FeatNum      int32  `protobuf:"varint,7,opt,name=featNum" json:"featNum,omitempty"`
	Feats        []byte `protobuf:"bytes,8,opt,name=feats,proto3" json:"feats,omitempty"`
}

type RetrievalFeatureResponse struct {
	Msg   string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	Pos   int32  `protobuf:"varint,2,opt,name=pos" json:"pos,omitempty"`
	Total int32  `protobuf:"varint,3,opt,name=total" json:"total,omitempty"`
}

type LifecycleRequest struct {
	Uuid         string `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
	CameraId     string `protobuf:"bytes,2,opt,name=cameraId" json:"cameraId,omitempty"`
	RepositoryId string `protobuf:"bytes,3,opt,name=repositoryId" json:"repositoryId,omitempty"`
	Type         int32  `protobuf:"varint,4,opt,name=type" json:"type,omitempty"`
}

type LifecycleResponse struct {
	Msg   string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	Total int32  `protobuf:"varint,3,opt,name=total" json:"total,omitempty"`
}

type RespMsg struct{
	Rtn int `json:"rtn"`
	Message string `json:"message"`
}