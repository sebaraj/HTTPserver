package server

/*
	func GetStandardResponse(statusCode HttpStatusCode) Response {
		return Response{
			Version:    "HTTP/1.1",
			StatusCode: statusCode,
			StatusText: string(getHttpMethodAsString(statusCode)[4:]),
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       string("<h1>" + getHttpMethodAsString(statusCode) + "</h1>"),
		}
	}

	func getHttpMethodAsString(statusCode HttpStatusCode) string {
		switch statusCode {
		case OK:
			return "200 OK"
		case BadRequest:
			return "400 Bad Request"
		case Unauthorized:
			return "401 Unauthorized"
		case Forbidden:
			return "403 Forbidden"
		case NotFound:
			return "404 Not Found"
		case MethodNotAllowed:
			return "405 Method Not Allowed"
		case NotAcceptable:
			return "406 Not Acceptable"
		case Conflict:
			return "409 Conflict"
		case InternalServerError:
			return "500 Internal Server Error"
		case NotImplemented:
			return "501 Not Implemented"
		case BadGateway:
			return "502 Bad Gateway"
		case ServiceUnavailable:
			return "503 Service Unavailable"
		default:
			return "200 OK"
		}
	}
*/
var ResponseBadRequest = Response{
	Version:    "HTTP/1.1",
	StatusCode: BadRequest,
	StatusText: "Bad Request",
	Headers:    map[string]string{},
	Body:       "",
}

var ResponseUnauthorized = Response{
	Version:    "HTTP/1.1",
	StatusCode: Unauthorized,
	StatusText: "Unauthorized",
	Headers:    map[string]string{},
	Body:       "",
}

var ResponseForbidden = Response{
	Version:    "HTTP/1.1",
	StatusCode: Forbidden,
	StatusText: "Forbidden",
	Headers:    map[string]string{},
	Body:       "",
}

var ResponseNotFound = Response{
	Version:    "HTTP/1.1",
	StatusCode: NotFound,
	StatusText: "Not Found",
	Headers:    map[string]string{},
	Body:       "",
}

var ResponseMethodNotAllowed = Response{
	Version:    "HTTP/1.1",
	StatusCode: MethodNotAllowed,
	StatusText: "Method Not Allowed",
	Headers:    map[string]string{},
	Body:       "",
}

var ResponseInternalServerError = Response{
	Version:    "HTTP/1.1",
	StatusCode: InternalServerError,
	StatusText: "Internal Server Error",
	Headers:    map[string]string{},
	Body:       "",
}

var ResponseNotImplemented = Response{
	Version:    "HTTP/1.1",
	StatusCode: NotImplemented,
	StatusText: "Not Implemented",
	Headers:    map[string]string{},
	Body:       "",
}

var ResponseBadGateway = Response{
	Version:    "HTTP/1.1",
	StatusCode: BadGateway,
	StatusText: "Bad Gateway",
	Headers:    map[string]string{},
	Body:       "",
}

var ResponseServiceUnavailable = Response{
	Version:    "HTTP/1.1",
	StatusCode: ServiceUnavailable,
	StatusText: "Service Unavailable",
	Headers:    map[string]string{},
	Body:       "",
}
