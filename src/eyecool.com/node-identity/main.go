package main

import (
	"eyecool.com/node-identity/http"
	_ "eyecool.com/node-identity/timer"
)

func main() {
	http.StartWebService()
}
