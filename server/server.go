package server

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/bryanwsebaraj/httpserver/socket"
)

type Server struct {
	Socket socket.Socket
	Routes map[Route]RouteHandler
}

func (s *Server) StartServer(port uint16, addr uint32, bcklog int) error {
	//sock := new(socket.Socket)
	// Initialize initializes a socket connection using the specified address family and service.
	// The address family is typically syscall.AF_INET for IPv4 or syscall.AF_INET6 for IPv6.
	// The service parameter specifies the specific service or protocol to use with the socket.
	// Returns the initialized socket connection and an error, if any.
	var err error = errors.New("")
	err = s.Socket.CreateSocket(syscall.AF_INET, syscall.SOCK_STREAM, 0, port, addr, bcklog)
	if err != nil {
		println("create failed")
		return err
	}
	//sock.CloseSocket()

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

func (s *Server) InitializeRoutes(r *map[Route]RouteHandler) error {
	s.Routes = *r
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
		//println("waiting to stop")
		<-done
		exitFlag = true
		s.Socket.CloseSocket()
		//print("closed socket ctrlc")
		return
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
			//println("accept failed")
			//println(newFd)
			//println(sock.GetFD())
			s.Socket.CloseSocket()
			//println("closed socket from err")
			return err
		}
		//println("new connection accepted")
		go handleRequest(newFd, s.Routes)

	}
	//sock.CloseSocket()
	//println("closed socket at end")
	return nil

}
