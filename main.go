package main

import (
	"log"
	"net/http/cgi" //#nosec G504 -- Use Go versions > 1.17
	"os"

	"golang.org/x/net/webdav"
)

func main() {
	dir := os.Getenv("GODAV_DIR")
	if dir == "" {
		dir = os.Getenv("DOCUMENT_ROOT")
	}

	prefix := os.Getenv("GODAV_URL_PREFIX")
	if prefix == "" {
		prefix = os.Getenv("SCRIPT_NAME")
	}
	if prefix == "" {
		prefix = "/"
	}

	// Validate directory
	if dir == "" {
		log.Fatal("GODAV_DIR or DOCUMENT_ROOT environment variable must be set")
	}

	// Log startup
	log.Printf("Starting godavd with DIR=%s, PREFIX=%s", dir, prefix)

	davHandler := &webdav.Handler{
		Prefix:     prefix,
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	}

	if err := cgi.Serve(davHandler); err != nil {
		log.Fatal(err)
	}
}
