package utils

import (
	"strings"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"errors"
	"github.com/satori/go.uuid"
	"fmt"
)

func GetIdAndClusterId(id string) (groupId int, clusterId int, err error) {
	var err1, err2 error
	splitId := strings.Split(id, "@")
	if len(splitId) == 2 {
		if splitId[0] == "" {
			groupId = -2 //代表空串
		} else {
			groupId, err1 = strconv.Atoi(splitId[0])
		}
		if splitId[1] == "" {
			clusterId = -2 //代表空串
		} else {
			clusterId, err2 = strconv.Atoi(splitId[1])
		}
		if err1 != nil || err2 != nil {
			return -1, -1, errors.New("GetIdAndClusterId error")
		}
		return groupId, clusterId, nil
	}
	return -1, -1, errors.New("GetIdAndClusterId error")
}

func GetClusterIdAndId(str string) (int, int, error) {
	m := strings.Split(str, "@")
	id, err := strconv.Atoi(m[0])
	clusterId, err := strconv.Atoi(m[1])
	return id, clusterId, err

}

func MD5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func UUID() string {
	uuid, _ := uuid.NewV4()
	uuidStr := strings.Replace(uuid.String(), "-", "", -1)
	return uuidStr
}

func SplitArrayByComma(data []string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(data), "[]"), " ", ",", -1)
}
