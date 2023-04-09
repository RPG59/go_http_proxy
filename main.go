package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)


type HttpProxy struct {
	client *http.Client
}

var dataMap sync.Map

func (p *HttpProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var key = strings.Split(req.URL.Path, "/")[1]
	data, ok := dataMap.Load(key)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var url = fmt.Sprintf("http://localhost/%s/test.jpg", data)

	res, err := p.client.Get(url)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.Copy(w, res.Body)
	res.Body.Close()
}

func (p *HttpProxy) init() {
	p.client = &http.Client{}
}


func main() {
	dataMap.LoadOrStore("foo", "testImg")
	var httpProxy = HttpProxy{}

	httpProxy.init()

	err := http.ListenAndServe("127.0.0.1:6969",  &httpProxy)

	if err != nil {
		log.Fatal("Server Error: ", err)
	}

}