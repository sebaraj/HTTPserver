package server

import (
	"errors"
	"os"
	"os/signal"
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
	//sock.CloseSocket()

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

	//ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	//	defer stop()

	done := make(chan os.Signal, 1)
	exitFlag := false
	signal.Notify(done, os.Interrupt)

	go func(sock *socket.Socket) {
		//println("waiting to stop")
		<-done
		exitFlag = true
		sock.CloseSocket()
		//print("closed socket ctrlc")
		return
	}(sock)

out:
	for {
		if exitFlag {
			break out
		}
		// accept connection
		newFd, _, err := sock.AcceptSocket()
		if err != nil {
			//println("accept failed")
			//println(newFd)
			//println(sock.GetFD())
			sock.CloseSocket()
			//println("closed socket from err")
			return err
		}
		//println("new connection accepted")
		go handleRequest(newFd)

	}
	//sock.CloseSocket()
	//println("closed socket at end")
	return nil

}

func handleRequest(newFd int) {
	// in the future, handle persistent connections (keep-alive). currently closing after every response (HTTP 1.0)
	// convert this to loop to read until end, not just 10000
	buffer := make([]byte, 10000)
	valread, err := syscall.Read(newFd, buffer)
	if err != nil {
		println("read failed")
		return //err
	}
	for i := 0; i < valread; i++ {
		print(string(buffer[i]))
	}
	syscall.Write(newFd, []byte("HTTP/1.1 200 OK\nContent-Type: text/html\n\n<h1>Hello, World!</h1>"))
	syscall.Close(newFd)
}

type Request struct {
	Method  int // HttpMethod
	URI     string
	Version string
	Headers map[string]string
	Body    string
}

type Response struct {
	Version    string
	StatusCode uint16
	StatusText string
	Headers    map[string]string
	Body       string
}
