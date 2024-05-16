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

func AboutHandler(req *Request) *Response {
	return &Response{
		Version:    req.Version,
		StatusCode: OK,
		StatusText: "OK",
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       "<h1>About</h1>",
	}

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
	/*sock := new(socket.Socket)
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
	*/
	//ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	//	defer stop()

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

func handleRequest(newFd int, routes map[Route]RouteHandler) {
	// in the future, handle persistent connections (keep-alive). currently closing after every response (HTTP 1.0)
	// convert this to loop to read until end, not just 10000
	buffer := make([]byte, 10000)
	valread, err := syscall.Read(newFd, buffer)
	if err != nil {
		println("read failed")
		return //err
	}
	println("receiving request:")
	for i := 0; i < valread; i++ {
		print(string(buffer[i]))
	}

	// parse request
	req := parseRequest(buffer)
	// handle route
	res := HandleRoute(req, routes)
	bufferOut := parseResponse(res)
	//res = []byte(res)
	println("sending response:")
	for i := 0; i < 30; i++ {
		print(string(bufferOut[i]))
	}
	syscall.Write(newFd, bufferOut)
	syscall.Close(newFd)
}

func parseResponse(res *Response) []byte {
	//bufferOut := make([]byte, 10000)
	// write response
	// strconv.FormatUint((uint64(res.StatusCode)), 3)
	bufferOut := []byte(res.Version + " " + "200" + " " + res.StatusText + "\n" + "Content-Type: text/html\n\n" + res.Body + "\n")
	//bufferOut = []byte("HTTP/1.1 200 OK\nContent-Type: text/html\n\n<h1>Hello, World!</h1>")
	//bufferOut = append(bufferOut, []byte(res.Version+" "+string(res.StatusCode)+" "+res.StatusText+"\r\n")+[]byte("Content-Type: "+res.Headers["Content-Type"]+"\r\n")+[]byte("Content-Length: "+string(len(res.Body))+"\r\n\r\n")+[]byte(res.Body)+"\r\n") // add headers
	return bufferOut
}

func parseRequest(buffer []byte) *Request {
	// parse request
	req := new(Request)
	// parse method
	req.Method = GET
	req.URI = "/about"
	req.Version = "HTTP/1.1"
	req.Headers = map[string]string{"Content-Type": "text/html"}
	req.Body = "<h1></h1>"

	return req
}

func HandleRoute(req *Request, routes map[Route]RouteHandler) *Response {
	path := req.URI
	handler, ok := routes[Route(path)]
	if !ok {
		return &Response{
			Version:    req.Version,
			StatusCode: HttpStatusCode(404),
			StatusText: "Not Found",
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       "<h1>404 Not Found</h1>",
		}
	}
	return handler(req)
	// break down reponse as bytestream to write via syscall
}

type HttpMethod uint16

const (
	GET HttpMethod = iota
	POST
	PUT
	DELETE
)

type HttpStatusCode uint16

const (
	OK                  HttpStatusCode = 200
	BadRequest                         = 400
	Unauthorized                       = 401
	Forbidden                          = 403
	NotFound                           = 404
	MethodNotAllowed                   = 405
	NotAcceptable                      = 406
	Conflict                           = 409
	InternalServerError                = 500
	NotImplemented                     = 501
	BadGateway                         = 502
	ServiceUnavailable                 = 503
	// Add more status codes here
)

type Route string

type RouteHandler func(req *Request) *Response

type Request struct {
	Method  HttpMethod
	URI     string
	Version string
	Headers map[string]string
	Body    string
}

type Response struct {
	Version    string
	StatusCode HttpStatusCode
	StatusText string
	Headers    map[string]string
	Body       string
}
