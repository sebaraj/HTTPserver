package server

import (
	"strconv"
	"strings"
	"syscall"
	"time"
)

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

const ()

func (req *Request) AsString() string {
	method := ""
	switch req.Method {
	case GET:
		method = "GET"
	case POST:
		method = "POST"
	case PUT:
		method = "PUT"
	case DELETE:
		method = "DELETE"
	case PATCH:
		method = "PATCH"
	case OPTIONS:
		method = "OPTIONS"
	case HEAD:
		method = "HEAD"
	case TRACE:
		method = "TRACE"
	case CONNECT:
		method = "CONNECT"
	}
	headersString := "Headers: \n"
	for key, value := range req.Headers {
		headersString += key + ": " + value + "\n"
	}

	//bodyString := "Body: \n"
	//for key, value := range req.Body {
	//	bodyString += key + ": " + value + "\n"
	//}
	return string("Method: " + method + "\nURI: " + req.URI + "\nVersion: " + req.Version + "\n" + headersString + "\nBody: " + req.Body)

}

/* non-persistent connection handling. working


func handleRequest(newFd int, routes map[Route]RouteHandler, rtTree *RouteNode) {
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
	bufferOut := parseResponse(HandleRoute(req, routes, *rtTree))
	println("sending response:")
	for i := 0; i < len(bufferOut); i++ {
		print(string(bufferOut[i]))
	}
	syscall.Write(newFd, bufferOut)
	syscall.Close(newFd)
}*/

func handleRequest(newFd int, routes map[Route]RouteHandler, rtTree *RouteNode) {
	// in the future, handle persistent connections (keep-alive). currently closing after every response (HTTP 1.0)
	// convert this to loop to read until end, not just 10000
	buffer := make([]byte, 10000)
	// while has time, try read. if read == 0 and time out, close connection
	valread := 0
	syscall.SetNonblock(newFd, true)
	err := error(nil)
	timeOut := time.Now().Add(5 * time.Second)
	for time.Now().Before(timeOut) {
		println("trying read before timeout")
		//syscall.Write(newFd, []byte("$"))
		//buffer = buffer[:10000]
		valread, err = syscall.Read(newFd, buffer) // is this nonblocking? this might be blocking https://stackoverflow.com/questions/36112445/will-go-block-the-current-thread-when-doing-i-o-inside-a-goroutine
		//println("valread: ", valread)
		//println("err: ", err)

		//if err != nil {
		//	println("read failed")
		//	return //err
		//} else
		if valread > 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if valread == -1 || valread == 0 {
		println("timeout")
		syscall.Close(newFd)
		return
	}
	//valread, err := syscall.Read(newFd, buffer)
	if err != nil {
		println("read failed")
		//	return //err
	}
	println("receiving request:")
	for i := 0; i < valread; i++ {
		print(string(buffer[i]))
	}

	// parse request
	req := parseRequest(buffer[1:valread])
	res := new(Response)
	if req.Headers["Connection"] == "close" || (req.Headers["Connection"] != "keep-alive" && req.Version == "HTTP/1.0") {
		res = HandleRoute(req, routes, *rtTree)
		if int(int(res.StatusCode)/100) != 1 {
			res.Headers["Connection"] = "close"
		}
		bufferOut := parseResponse(res)
		println("sending response:")
		//for i := 0; i < len(bufferOut); i++ {
		//	print(string(bufferOut[i]))
		//}
		syscall.Write(newFd, bufferOut)
		syscall.Close(newFd)
	} else {
		// keep-alive
		println("keep-alive")
		//println("!!!!!!!!!!!!!!!keep-alive")
		res = HandleRoute(req, routes, *rtTree)
		if int(int(res.StatusCode)/100) != 1 {
			res.Headers["Connection"] = "keep-alive"
			println("added keep-alive connection header: ")
		}
		bufferOut := parseResponse(res)
		println("sending response:")
		//for i := 0; i < len(bufferOut); i++ {
		//	print(string(bufferOut[i]))
		//}
		syscall.Write(newFd, bufferOut)
		handleRequest(newFd, routes, rtTree)
		//syscall.Close(newFd)
	}

	// handle route
	//bufferOut := parseResponse(res)
	//println("sending response:")
	//for i := 0; i < len(bufferOut); i++ {
	//	print(string(bufferOut[i]))
	//}
	//syscall.Write(newFd, bufferOut)
	//syscall.Close(newFd)
}

func parseRequest(buffer []byte) *Request {
	req := new(Request)
	firstRow := strings.Split(string(buffer), "\n")[0]

	req.Method = getHttpMethod(strings.Split(firstRow, " ")[0])
	req.URI = strings.Split(firstRow, " ")[1]
	req.Version = strings.Split(firstRow, " ")[2]
	req.Headers = map[string]string{}
	//req.Body = ""

	//inBody := false
	count := 0
	for _, line := range strings.Split(string(buffer), "\n")[1:] {
		count++
		if line == "\r" {
			break
		} else {
			header := strings.Split(line, ": ")
			//strings.TrimRight(req.Headers["Connection"], "\r\n")
			req.Headers[header[0]] = strings.TrimRight(header[1], "\r\n")
		}
	}
	count++
	req.Body = strings.Join(strings.Split(string(buffer), "\n")[count:], "")

	return req
}

func parseResponse(res *Response) []byte {
	// res.Version hardcoded for now (all responses are HTTP/1.1ish anyways)
	headersString := ""
	for key, value := range res.Headers {
		headersString += key + ": " + value + "\n"
	}
	bufferOut := []byte("HTTP/1.1" + " " + strconv.Itoa(int(res.StatusCode)) + " " + res.StatusText + "\n" + headersString + "\n" + res.Body + "\n")
	return bufferOut
}
