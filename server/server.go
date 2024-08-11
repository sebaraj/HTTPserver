package server

import (
	"encoding/binary"
	"errors"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sebaraj/httpserver/socket"
)

// GOMAXPROCS is set to number of CPU cores available by default

type Server struct {
	Socket    socket.Socket
	Routes    map[Route]RouteHandler
	RouteTree *RouteNode
}

func (s *Server) StartServer(port string, addr string, bcklog int) error {
	// Initialize initializes a socket connection using the specified address family and service.
	// The address family is typically syscall.AF_INET for IPv4 or syscall.AF_INET6 for IPv6.
	// The service parameter specifies the specific service or protocol to use with the socket.
	// Returns the initialized socket connection and an error, if any.

	var err error = errors.New("")
	portUint, err := strconv.ParseInt(port, 0, 16)
	if err != nil {
		println("port conversion failed")
		return err
	}
	addrUint := binary.BigEndian.Uint32(net.ParseIP(addr).To4())
	err = s.Socket.CreateSocket(syscall.AF_INET, syscall.SOCK_STREAM, 0, uint16(portUint), addrUint, bcklog)
	if err != nil {
		println("create failed")
		return err
	}

	err = s.Socket.BindSocket()
	if err != nil {
		println("bind failed")
		return err
	}

	err = s.Socket.ListenSocket()
	if err != nil {
		println("listening failed")
		return err
	}
	return nil
}

func (s *Server) ListenAndServe() error {
	// Initialize initializes a socket connection using the specified address family and service.
	// The address family is typically syscall.AF_INET for IPv4 or syscall.AF_INET6 for IPv6.
	// The service parameter specifies the specific service or protocol to use with the socket.
	// Returns the initialized socket connection and an error, if any.
	done := make(chan os.Signal, 1)
	exitFlag := false
	signal.Notify(done, os.Interrupt)

	go func(s *Server) {
		<-done
		exitFlag = true
		s.Socket.CloseSocket()
	}(s)

out:
	for {
		if exitFlag {
			break out
		}
		// accept connection
		newFd, _, err := s.Socket.AcceptSocket()
		if err != nil {
			// handle this better so it doesn't always close the socket/server unless desired
			s.Socket.CloseSocket()
			return err
		}
		syscall.SetNonblock(newFd, true)
		go handleRequest(newFd, s.Routes, s.RouteTree)

	}
	return nil
}
