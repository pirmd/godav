package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/cgi" //#nosec G504 -- Use Go versions > 1.17
	"os"

	"golang.org/x/net/webdav"
)

// enableCORS adds CORS headers to allow cross-origin requests from web frontends
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, PROPFIND, MKCOL, COPY, MOVE, LOCK, UNLOCK, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, Authorization, Destination, Overwrite")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Content-Length, ETag")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Parse command-line flags
	var (
		dirFlag    = flag.String("dir", "", "Directory to expose (overrides GODAV_DIR and DOCUMENT_ROOT)")
		prefixFlag = flag.String("prefix", "", "URL prefix (overrides GODAV_URL_PREFIX and SCRIPT_NAME)")
	)
	flag.Parse()

	// Configure directory (priority: flag > GODAV_DIR > DOCUMENT_ROOT)
	dir := *dirFlag
	if dir == "" {
		dir = os.Getenv("GODAV_DIR")
	}
	if dir == "" {
		dir = os.Getenv("DOCUMENT_ROOT")
	}

	// Configure URL prefix (priority: flag > GODAV_URL_PREFIX > SCRIPT_NAME > "/")
	prefix := *prefixFlag
	if prefix == "" {
		prefix = os.Getenv("GODAV_URL_PREFIX")
	}
	if prefix == "" {
		prefix = os.Getenv("SCRIPT_NAME")
	}
	if prefix == "" {
		prefix = "/"
	}

	// Validate directory
	if dir == "" {
		log.Fatal("DIR must be set via -dir flag or GODAV_DIR/DOCUMENT_ROOT environment variable")
	}

	// Log startup
	log.Printf("Starting godavd with DIR=%s, PREFIX=%s", dir, prefix)

	davHandler := enableCORS(&webdav.Handler{
		Prefix:     prefix,
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	})

	if err := cgi.Serve(davHandler); err != nil {
		log.Fatal(err)
	}
}
