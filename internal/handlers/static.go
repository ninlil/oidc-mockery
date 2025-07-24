package handlers

import (
	"net/http"
	"strings"
)

// HandleStatic serves static files from the static directory
func HandleStatic(w http.ResponseWriter, r *http.Request) {
	// Remove /static/ prefix from the path
	path := strings.TrimPrefix(r.URL.Path, "/static/")

	// Serve files from the static directory
	http.ServeFile(w, r, "static/"+path)
}
