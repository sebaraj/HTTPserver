to-do

handle path variables in router - done 

helper function to get path variables in handler - done

helper function to easily print server.Request and predefined server.Response structs to return on error - done

handle all carriage returns (/r) in hashtables - done 

basic persistent/closed connections - done

handle more errors for system resiliency - done

add documentation in pkg.go.dev (need tagged/stable version)

add debug mode 

persistent connections/chunked transfer, as specified in HTTP 1.1 (race condition > 408 request timeout) and full HTTP 1.1 (host field, options method, caching support, 100 continue status)

handle query parameters in router

limit number of open connections from a single client (characteristic of denial of service)