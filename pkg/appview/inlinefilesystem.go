package appview

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/markbates/pkger"
)

// InlineFilesystem implements the http.Filesystem interface,
// which allows it to load assets from within the binary.
type InlineFilesystem struct {
	SourceFolder string
}

// NewInlineFilesystem creates a new inline filesystem based on
// the provided URL.
func NewInlineFilesystem(sourceFolder string) *InlineFilesystem {
	return &InlineFilesystem{
		SourceFolder: sourceFolder,
	}
}

// Open  Will only be called if the local HTTP server is active.
func (ifs *InlineFilesystem) Open(name string) (http.File, error) {
	return pkger.Open(fmt.Sprintf("%s/%s", ifs.SourceFolder, name))
}

// Server configures a local HTTP server to serve
// inline assets of the provided file system.
func (ifs *InlineFilesystem) Server() *httptest.Server {
	return httptest.NewServer(http.FileServer(ifs))
}
