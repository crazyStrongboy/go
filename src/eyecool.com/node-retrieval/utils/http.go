package utils

import (
	"log"
	"io/ioutil"
	"bytes"
	"net/http"
)

var POST_METHOD  = "POST"

//body提交二进制数据
func DoBytesPost(url string, data []byte) ([]byte, error) {

	body := bytes.NewReader(data)
	request, err := http.NewRequest(POST_METHOD, url, body)
	if err != nil {
		log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "Keep-Alive")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("http.Do failed,[err=%s][url=%s]", err, url)
	}
	log.Println("http.Do finish,[err=%s][url=%s][resp =%s]", err, url,string(b))

	return b, err
}
