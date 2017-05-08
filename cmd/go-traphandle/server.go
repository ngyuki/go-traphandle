package main

import (
	"io/ioutil"
	"log"
	"net"
)

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

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Printf("accept error ... %v", err)
			continue
		}

		go func(conn *net.TCPConn) {
			defer conn.Close()

			input, err := ioutil.ReadAll(conn)
			if err != nil {
				log.Printf("accept %v recv error %v", conn.RemoteAddr().String(), err)
				return
			}

			log.Printf("accept %v recv %v bytes", conn.RemoteAddr().String(), len(input))

			callback(input)

		}(conn)
	}
}
