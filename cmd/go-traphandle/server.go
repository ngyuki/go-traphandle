package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

type Handler struct {
	callback func([]byte)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s %s %s %s \"%s\": ", r.RemoteAddr, r.Method, r.URL, r.Proto, r.UserAgent())

	if strings.ToUpper(r.Method) != "POST" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("request error ... ", err)
		return
	}

	h.callback(body)
	log.Println("request complete")
}

func startServer(server string, callback func([]byte)) {

	addr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		panic(err)
	}

	listen, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}

	log.Printf("listen ... %v", addr)

	http.Handle("/", &Handler{callback})
	log.Println(http.Serve(listen, nil))
}
