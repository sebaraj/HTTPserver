package main

import (
	"encoding/binary"
	"log"
	"net"
	"os"

	"strconv"

	"github.com/bryanwsebaraj/httpserver/server"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Cannot load env. Through %v", err)
	}
	ip := binary.BigEndian.Uint32(net.ParseIP(os.Getenv("IP_ADDRESS")).To4())
	port, _ := strconv.ParseInt(os.Getenv("PORT"), 0, 16)
	println("http server on!")
	s := new(server.Server)
	s.StartServer(uint16(port), ip, 10)

	//server.ListenAndServe(uint16(port), ip, 10)

	sampleMap := map[server.Route]server.RouteHandler{
		"/":      HomeHandler,
		"/about": AboutHandler,
	}
	s.InitializeRoutes(&sampleMap)
	s.ListenAndServe()

	println("\nhttp server off!")

}
