package server

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
		Headers:    map[string]string{"Content-Type": "text/html", "Date": time.Now().String()},
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
	for i := 0; i < len(bufferOut); i++ {
		print(string(bufferOut[i]))
	}
	syscall.Write(newFd, bufferOut)
	syscall.Close(newFd)
}

func parseResponse(res *Response) []byte {
	//bufferOut := make([]byte, 10000)
	// write response
	// strconv.FormatUint((uint64(res.StatusCode)), 3)
	//println(res.Version)
	// figure out res.Version and print all headers
	headersString := ""
	for key, value := range res.Headers {
		headersString += key + ": " + value + "\n"
	}

	bufferOut := []byte("HTTP/1.0 " + "200" + " " + res.StatusText + "\n" + headersString + "\n" + res.Body + "\n")

	//bufferOut = []byte("HTTP/1.1 200 OK\nContent-Type: text/html\n\n<h1>Hello, World!</h1>")
	//bufferOut = append(bufferOut, []byte(res.Version+" "+string(res.StatusCode)+" "+res.StatusText+"\r\n")+[]byte("Content-Type: "+res.Headers["Content-Type"]+"\r\n")+[]byte("Content-Length: "+string(len(res.Body))+"\r\n\r\n")+[]byte(res.Body)+"\r\n") // add headers
	return bufferOut
}

func getHttpMethod(method string) HttpMethod {
	switch method {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	default:
		return GET
	}

}

func parseRequest(buffer []byte) *Request {
	req := new(Request)
	firstRow := strings.Split(string(buffer), "\n")[0]

	req.Method = getHttpMethod(strings.Split(firstRow, " ")[0])
	req.URI = strings.Split(firstRow, " ")[1]
	req.Version = strings.Split(firstRow, " ")[2]

	//req.Method = GET
	//req.URI = "/about"
	//req.Version = "HTTP/1.1"

	// parse headers

	req.Headers = map[string]string{}
	req.Body = map[string]string{}

	inBody := false
	count := 0
	for _, line := range strings.Split(string(buffer), "\n")[1:] {
		count++
		if line == "\r" {
			break
		} else {
			header := strings.Split(line, ": ")
			println(header[0] + ": " + header[1])
			req.Headers[header[0]] = header[1]
		}
	}
	count++

	// naive implementation, but it works for now for json. update later
	if count < len(strings.Split(string(buffer), "\n")) {
		// theres a body
		for _, line := range strings.Split(string(buffer), "\n")[count:] {
			if line == "\r" {
				break
			}
			lineTrimmed := strings.TrimLeft(line, " ")
			if string(lineTrimmed[0]) == "{" {
				inBody = true
				println("in body")
			} else if string(lineTrimmed[0]) == "}" {
				break
			} else if inBody {
				body := strings.Split(string(lineTrimmed), ": ")
				println(body[0] + ": " + body[1])
				// need to strip quotes
				trimVal := strings.Trim(body[1], ",")
				req.Body[strings.Trim(body[0], "\"")] = strings.Trim(trimVal, "\"")
			}
		}

	}
	/* for debugging, delete
	println(req.Method)
	for key := range req.Body {
		println(key)
		println(req.Body[key])
	}
	println(req.Body["college"])
	println(req.Headers["Host"])
	*/
	// parse body (for PUT/POST)
	//req.Body = map[string]string{"html": "<h><h>"}

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
	Body    map[string]string
}

type Response struct {
	Version    string
	StatusCode HttpStatusCode
	StatusText string
	Headers    map[string]string
	Body       string
}
