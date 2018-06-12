package logic

import (
	"eyecool.com/node-retrieval/model"
	. "eyecool.com/node-retrieval/db"
	"strconv"
)

type VideoLogic struct{}

var DefaultVideo = VideoLogic{}

type Video struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`
	PermissionMap string `json:"permission_map"`
	Enabled int `json:"enabled"`
	RecParams string `json:"rec_params"`
	CreatorId string `json:"creator_id"`
	ExtraMeta string `json:"extra_meta"`
}

func (VideoLogic) QueryVideo() ([]*Video,error){
	objLog := GetLogger(nil)
	videos:=make([]*Video, 0)
	video := make([]*model.Video, 0)
	err := MasterDB.Find(&video)
	if err != nil {
		objLog.Errorln("VideoLogic QueryVideo PageAd error:", err)
		return nil,err
	}
	for _,v:=range video{
		vi:=&Video{}
		pkId:=strconv.Itoa(v.PkId)
		cluster:=strconv.Itoa(v.ClusterId)
		creatorId:=strconv.FormatInt(v.CreatorId,10)
		vi.Id=pkId+"@"+cluster
		vi.Name=v.Name
		vi.Url=v.Url
		vi.PermissionMap=v.PermissionMap
		vi.Enabled=v.Enabled
		vi.RecParams=v.RecParams
		vi.CreatorId=creatorId
		vi.ExtraMeta=v.ExtraMeta
		videos=append(videos, vi)
	}
	return videos,nil
}

func (VideoLogic)InsertVideo(video *model.Video)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Insert(video)
	if err != nil {
		session.Rollback()
		logger.Errorln("insert video error:", err)
		return  err
	}
	session.Commit()
	return nil
}

func (VideoLogic)UpdateVideo(video *model.Video)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("pk_id=?",video.PkId).Update(video)
	if err != nil {
		session.Rollback()
		logger.Errorln("update video error:", err)
		return  err
	}
	session.Commit()
	return nil
}

func (VideoLogic)DeleteVideo(video *model.Video)error{
	logger := GetLogger(nil)
	session := MasterDB.NewSession()
	defer session.Close()
	session.Begin()
	_, err := MasterDB.Where("pk_id=?",video.PkId).Delete(video)
	if err != nil {
		session.Rollback()
		logger.Errorln("delete video error:", err)
		return  err
	}
	session.Commit()
	return nil
}
