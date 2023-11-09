package server

import (
	"log"
	"net/http"
	"os"
	"strings"
	"yeetfile/web/utils"
)

type Route struct {
	Path   string
	Method string
}

type router struct {
	routes   map[Route]http.HandlerFunc
	reserved []string
}

// ServeHTTP finds the proper routing handler for the provided path.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for el, handler := range r.routes {
		if r.matchPath(el.Path, req.URL.Path) && el.Method == req.Method {
			if os.Getenv("YEETFILE_DEBUG") == "1" {
				log.Printf("%s %s\n", req.Method, req.URL)
			}
			handler(w, req)
			return
		}
	}

	log.Printf("Error: %s %s", req.Method, req.URL)
	http.NotFound(w, req)
}

// matchPath takes a URL path and determines if it's a match for a particular
// handler. This allows differentiating between two paths where the only
// difference is a wildcard (i.e. "/u" and "/u/*" for uploadInit and uploadData)
func (r *router) matchPath(pattern, path string) bool {
	parts := strings.Split(pattern, "/")
	segments := strings.Split(path, "/")

	isWildcard := parts[1] == "*"
	isEndpoint := utils.Contains(r.reserved, segments[1])

	if len(parts) != len(segments) || (isWildcard && isEndpoint) {
		return false
	}

	for i, part := range parts {
		if part == "*" && len(path) > 1 {
			continue
		}

		if part != segments[i] {
			return false
		}
	}

	return true
}
