package server

import (
	"errors"
	"syscall"

	"github.com/bryanwsebaraj/httpserver/socket"
)

func ListenAndServe(port uint16, addr uint32, bcklog int) error {
	sock := new(socket.Socket)
	// Initialize initializes a socket connection using the specified address family and service.
	// The address family is typically syscall.AF_INET for IPv4 or syscall.AF_INET6 for IPv6.
	// The service parameter specifies the specific service or protocol to use with the socket.
	// Returns the initialized socket connection and an error, if any.
	var err error = errors.New("")
	err = sock.CreateSocket(syscall.AF_INET, syscall.SOCK_STREAM, 0, port, addr, bcklog)
	if err != nil {
		println("create failed")
		return err
	}

	err = sock.BindSocket()
	if err != nil {
		println("bind failed")
		return err
	}

	err = sock.ListenSocket()
	if err != nil {
		println("listening failed")
		return err
	}

	for {
		// accept connection
		println("accepting connection on port", sock.GetSockaddr().Port)
		newFd, _, err := sock.AcceptSocket()
		if err != nil {
			println("accept failed")
			return err
		}
		println("new connection accepted")

		//go func(newFd int) {
		// FILEPATH: /Users/bryansebaraj/Workspace/HTTPserver/server/server.go
		buffer := make([]byte, 10000)
		valread, err := syscall.Read(newFd, buffer)
		if err != nil {
			println("read failed")
			return err
		}
		for i := 0; i < valread; i++ {
			print(string(buffer[i]))
		}
		syscall.Write(newFd, []byte("HTTP/1.1 200 OK\nContent-Type: text/html\n\n<h1>Hello, World!</h1>"))
		println("response sent")
		syscall.Close(newFd)
		//}(newFd)

	}
	// close connection
	return nil

}
