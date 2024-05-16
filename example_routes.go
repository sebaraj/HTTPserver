package main

import (
	"time"

	"github.com/bryanwsebaraj/httpserver/server"
)

func AboutHandler(req *server.Request) *server.Response {
	return &server.Response{
		Version:    req.Version,
		StatusCode: server.OK,
		StatusText: "OK",
		Headers:    map[string]string{"Content-Type": "text/html", "Date": time.Now().String()},
		Body:       "<h1>About</h1>",
	}

}

func HomeHandler(req *server.Request) *server.Response {
	return &server.Response{
		Version:    req.Version,
		StatusCode: server.OK,
		StatusText: "OK",
		Headers:    map[string]string{"Content-Type": "text/html", "Date": time.Now().String()},
		Body:       "<h1>Home</h1>",
	}
}
