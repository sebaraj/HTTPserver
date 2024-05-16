package server

type Route string

type RouteHandler func(req *Request) *Response

func HandleRoute(req *Request, routes map[Route]RouteHandler) *Response {
	path := req.URI
	handler, ok := routes[Route(path)]
	if !ok {
		return &Response{
			Version:    req.Version,
			StatusCode: NotFound,
			StatusText: "Not Found",
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       "<h1>404 Not Found</h1>",
		}
	}
	return handler(req)
}
