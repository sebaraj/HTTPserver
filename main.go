package main

import (
	"encoding/binary"
	"net"

	"github.com/bryanwsebaraj/httpserver/server"
)

func main() {
	ip := net.ParseIP("127.0.0.1")
	ipUint32 := binary.BigEndian.Uint32(ip.To4())
	server.ListenAndServe(9954, ipUint32, 10)
	println("http server created!")

}
