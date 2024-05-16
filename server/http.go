package server

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
