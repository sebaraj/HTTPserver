package server

import (
	"strings"
)

type Route string

type RouteHandler func(req *Request) *Response

type PathVars map[string]string

type RouteNode struct {
	CurrPath        string
	HasPathVarChild bool
	Children        map[string]*RouteNode
	HandlerFunc     RouteHandler
}

func (s *Server) InitializeRoutes(r *map[Route]RouteHandler) error {
	s.RouteTree = new(RouteNode)
	s.RouteTree.CurrPath = "/"
	s.RouteTree.Children = make(map[string]*RouteNode)
	s.RouteTree.HasPathVarChild = false
	s.RouteTree.HandlerFunc = nil
	root := s.RouteTree
	currentNode := root
	for route, handler := range *r {
		// split route
		routeParts := strings.Split(string(route), "/")
		if routeParts[1] == "" {
			routeParts = []string{"/"}
		} else {
			routeParts = append([]string{"/"}, routeParts[1:]...)
		}

		// add route parts to tree
		for index, part := range routeParts {
			if part == "" || part == "\n" {
				break
			}
			// first check if part is path var {}
			if string(part[0]) == "{" && string(part[len(part)-1]) == "}" {
				// if path var child node does not exist, create it
				if currentNode.HasPathVarChild == false {
					currentNode.HasPathVarChild = true
					child := new(RouteNode)
					child.CurrPath = "PATHVARIABLE"
					child.Children = make(map[string]*RouteNode)
					child.HasPathVarChild = false
					child.HandlerFunc = nil
					currentNode.Children["PATHVARIABLE"] = child
				}
				currentNode = currentNode.Children["PATHVARIABLE"]

			} else {
				// if part is not in tree, add it

				if currentNode.Children[part] == nil {
					child := new(RouteNode)
					child.CurrPath = part
					child.Children = make(map[string]*RouteNode)
					child.HasPathVarChild = false
					child.HandlerFunc = nil
					currentNode.Children[part] = child
					currentNode.HasPathVarChild = false

				}
				child := currentNode.Children[part]
				currentNode = child
			}

			if index == len(routeParts)-1 {
				// if last part, add handler function
				currentNode.HandlerFunc = handler
				break
			}
		}
		currentNode = root
	}
	return nil
}

func HandleRoute(req *Request, routes map[Route]RouteHandler, rtTree RouteNode) *Response {
	path := req.URI
	currentNode := rtTree
	// split path
	pathParts := strings.Split(path, "/")
	if pathParts[1] == "" {
		pathParts = []string{"/"}
	} else {
		pathParts = append([]string{"/"}, pathParts[1:]...)
	}
	for index, part := range pathParts {
		if index == len(pathParts)-1 {
			if currentNode.Children[part] == nil {
				currentNode = *currentNode.Children["PATHVARIABLE"]
			} else {
				currentNode = *currentNode.Children[part]
			}
			break
		}

		if currentNode.Children[part] != nil {
			currentNode = *currentNode.Children[part]
		} else {
			if currentNode.HasPathVarChild == true {
				currentNode = *currentNode.Children["PATHVARIABLE"]
			} else {
				return &Response{
					Version:    req.Version,
					StatusCode: NotFound,
					StatusText: "Not Found",
					Headers:    map[string]string{"Content-Type": "text/html"},
					Body:       "<h1>404 Not Found</h1>",
				}
			}
		}

	}
	handler := currentNode.HandlerFunc
	if handler == nil {
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

// prepopulate bad return responses to call in handlers

func GetPathVars(req *Request, path string) PathVars {
	pathVars := make(map[string]string)
	pathParts := strings.Split(path, "/")
	reqParts := strings.Split(req.URI, "/")
	for index, part := range pathParts {
		if part == "" {
			continue
		}
		if string(part[0]) == "{" && string(part[len(part)-1]) == "}" {
			pathVars[part[1:len(part)-1]] = reqParts[index]
		}
	}
	return pathVars
}
