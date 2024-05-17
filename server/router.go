package server

import (
	"strings"
)

type Route string

type RouteHandler func(req *Request) *Response

type PathVars map[string]string

func (s *Server) InitializeRoutes(r *map[Route]RouteHandler) error {
	s.RouteTree = new(RouteNode)
	s.RouteTree.CurrPath = "/"
	s.RouteTree.Children = make(map[string]*RouteNode)
	s.RouteTree.HasPathVarChild = false
	s.RouteTree.HandlerFunc = nil
	root := s.RouteTree
	currentNode := root
	// handler is _
	for route, handler := range *r {
		// split route
		routeParts := strings.Split(string(route), "/")
		if routeParts[1] == "" {
			routeParts = []string{"/"}
		} else {
			routeParts = append([]string{"/"}, routeParts[1:]...)
		}

		//routeParts = append([]string("/"), routeParts)
		//println(routeParts)
		// add route parts to tree
		// second _ is part
		for index, part := range routeParts {
			println("path:", part)
			println("index:", index)
			if part == "" || part == "\n" {
				break
			}
			// first check if part is path var {}
			if string(part[0]) == "{" && string(part[len(part)-1]) == "}" {
				// if path var, add to path vars
				print("adding path var")
				// go to path var child node
				// if path var child node does not exist, create it
				if currentNode.HasPathVarChild == false {
					currentNode.HasPathVarChild = true
					child := new(RouteNode)
					child.CurrPath = "PATHVARIABLE"
					child.Children = make(map[string]*RouteNode)
					child.HasPathVarChild = false
					child.HandlerFunc = nil
					currentNode.Children["PATHVARIABLE"] = child
					println("added path var child", currentNode.Children["PATHVARIABLE"].CurrPath)
				}
				currentNode = currentNode.Children["PATHVARIABLE"]

				// if path var child node exists, go to it

			} else {
				//println("is error here")
				//if currentNode.Children[part] == nil {
				//	print("nil found")
				//}
				//_, err := currentNode.Children[part]
				//println(err)

				// if part is not in tree, add it

				if currentNode.Children[part] == nil {
					//println("is error here pt 2")
					println("adding new child", string(part))
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
				println("current Node current path:", currentNode.CurrPath)
			}

			if index == len(routeParts)-1 {
				// if last part, add handler function
				println("assigned handler function to", currentNode.CurrPath)
				currentNode.HandlerFunc = handler
				break
			}
			println(" ")
			//if  == nil {
			//	currentNode.Children = make(map[string]*RouteNode)
			//}
		}
		currentNode = root
		println(" ")

	}

	// old implementation
	s.Routes = *r
	return nil
}

func HandleRoute(req *Request, routes map[Route]RouteHandler, rtTree RouteNode) *Response {
	path := req.URI
	//handler, ok := routes[Route(path)]
	//CurrPathVars := make(map[string]string)
	//CurrPathVars["hi"] = "bye"
	// walk s.RouteTree to check
	println("currPath", rtTree.CurrPath)
	/*
		if !ok {
			return &Response{
				Version:    req.Version,
				StatusCode: NotFound,
				StatusText: "Not Found",
				Headers:    map[string]string{"Content-Type": "text/html"},
				Body:       "<h1>404 Not Found</h1>",
			}
		}
	*/
	currentNode := rtTree
	// split path
	pathParts := strings.Split(path, "/")
	//pathParts = append([]string{"/"}, pathParts[0:]...)
	if pathParts[1] == "" {
		pathParts = []string{"/"}
	} else {
		pathParts = append([]string{"/"}, pathParts[1:]...)
	}
	for index, part := range pathParts {
		println("path part:", part)
		println("index:", index)
		if index == len(pathParts)-1 {
			for key := range currentNode.Children {
				println("Key:", key)
			}
			if currentNode.Children[part] == nil {
				currentNode = *currentNode.Children["PATHVARIABLE"]
			} else {
				currentNode = *currentNode.Children[part]
			}
			//currentNode = *currentNode.Children[part]
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

		// this is bad logic
		/*
			if part[0] == '{' && part[len(part)-1] == '}' {
				// if path var, add to path vars
				// go to path var child node
				// if path var child node does not exist, create it
				if currentNode.HasPathVarChild == false {
					return &Response{
						Version:    req.Version,
						StatusCode: NotFound,
						StatusText: "Not Found",
						Headers:    map[string]string{"Content-Type": "text/html"},
						Body:       "<h1>404 Not Found</h1>",
					}
				} else {
					currentNode = *currentNode.Children["PATHVARIABLE"]
				}
				// if path var child node exists, go to it

			} else {
				if currentNode.Children[part] == nil {
					return &Response{
						Version:    req.Version,
						StatusCode: NotFound,
						StatusText: "Not Found",
						Headers:    map[string]string{"Content-Type": "text/html"},
						Body:       "<h1>404 Not Found</h1>",
					}
				}
				currentNode = *currentNode.Children[part]
			}
		*/
	}
	println("current node path:", currentNode.CurrPath)
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

type RouteNode struct {
	CurrPath        string
	HasPathVarChild bool
	Children        map[string]*RouteNode
	HandlerFunc     RouteHandler
}
