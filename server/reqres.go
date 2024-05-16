package server

import (
	"strconv"
	"strings"
	"syscall"
)

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
	//res := HandleRoute(req, routes)
	bufferOut := parseResponse(HandleRoute(req, routes))
	println("sending response:")
	for i := 0; i < len(bufferOut); i++ {
		print(string(bufferOut[i]))
	}
	syscall.Write(newFd, bufferOut)
	syscall.Close(newFd)
}

func parseRequest(buffer []byte) *Request {
	req := new(Request)
	firstRow := strings.Split(string(buffer), "\n")[0]

	req.Method = getHttpMethod(strings.Split(firstRow, " ")[0])
	req.URI = strings.Split(firstRow, " ")[1]
	req.Version = strings.Split(firstRow, " ")[2]
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
			} else if string(lineTrimmed[0]) == "}" {
				break
			} else if inBody {
				body := strings.Split(string(lineTrimmed), ": ")
				trimVal := strings.Trim(body[1], ",")
				req.Body[strings.Trim(body[0], "\"")] = strings.Trim(trimVal, "\"")
			}
		}

	}

	return req
}

func parseResponse(res *Response) []byte {
	// res.Version is being a little buggy, so hardcoding for now (all responses are HTTP/1.1ish anyways)
	headersString := ""
	for key, value := range res.Headers {
		headersString += key + ": " + value + "\n"
	}
	bufferOut := []byte("HTTP/1.1" + " " + strconv.Itoa(int(res.StatusCode)) + " " + res.StatusText + "\n" + headersString + "\n" + res.Body + "\n")
	return bufferOut
}
