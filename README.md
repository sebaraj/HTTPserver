# HTTP Server

A simple HTTP 1.1(ish) server written in Go, constructed from RFC 9112 without referencing other implementations. While functional, this is neither a production-ready, nor complete implementation of the HTTP 1.1 protocol.

### To-do

- [x] Process path variables in router

- [x] Debugging functions to print server.Request and server.Response structs to return on error

- [x] Handle all carriage returns (/r) in hashtable

- [x] Handle basic persistent/closed connections

- [x] Handle more errors for system resiliency

- [ ] Add documentation in pkg.go.dev (need tagged/stable version)

- [ ] Full persistent connections/chunked transfer, as specified in HTTP 1.1 (race condition > 408 request timeout) and full HTTP 1.1 (options method, caching support, 100 continue status)

- [ ] Handle query parameters in router

- [ ] Limit number of open connections from a single client (characteristic of denial of service)

