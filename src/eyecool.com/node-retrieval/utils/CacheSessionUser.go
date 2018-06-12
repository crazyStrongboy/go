package utils

import (
	"eyecool.com/node-retrieval/model"
)

var cacheMap = make(map[string]*model.User, 0)

type CacheMap struct {

}

func (this *CacheMap)GetInstance() map[string]*model.User{
	return cacheMap
}

func (this *CacheMap)SetUserSession(sessionId string, user *model.User) {
	cacheMap[sessionId] = user
}

func (this *CacheMap)CheckSession(sessionId string) bool {
	if _, f := cacheMap[sessionId]; f {
		return true
	}
	return false
}

func (this *CacheMap)ClearSession(sessionId string) {
	delete(cacheMap, sessionId)
}

func (this *CacheMap)GetUserSession(sessionId string) *model.User {
	if v, f := cacheMap[sessionId]; f {
		return v
	}
	return nil
}

func (this *CacheMap)GetUserClusterId(sessionId string) int {
	if v, f := cacheMap[sessionId]; f {
		return v.ClusterId
	}
	return 0
}

func (this *CacheMap)CheckMap() bool {
	return len(cacheMap) == 0
}
