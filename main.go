package main

import (
	"os"

	"github.com/bryanwsebaraj/httpserver/server"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		println("Cannot load env. Through %v", err)
	}
	println("http server on!")
	s := new(server.Server)
	s.StartServer(os.Getenv("SERVER_PORT"), os.Getenv("SERVER_IP"), 10)

	sampleMap := map[server.Route]server.RouteHandler{
		"/":            HomeHandler,
		"/about/about": AboutHandler,
		"/test/{path}": PathHandler,
		"/favicon.ico": IconHandler,
		"/json":        JsonExample,
	}
	s.InitializeRoutes(&sampleMap)
	s.ListenAndServe()

	println("\nhttp server off!")

}
