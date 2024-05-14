package server

import (
	"github.com/bryanwsebaraj/httpserver/tcpsocket"
)

func CreateServer() {
	tcpsocket.CreateSocket()
	println("server created!")

}
