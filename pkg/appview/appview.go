package appview

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"

	"github.com/markbates/pkger"
	"github.com/webview/webview"
)

// TODO: I am not sure if it is smart to hide the
// WebView API "just" because I don't want users
// to call `Navigate(url)`. Any thoughts are
// appreciated.

// AppView is a high-level API that allows you
// to display your React, Svelte, etc. app easily.
type AppView struct {
	URL string

	server  *httptest.Server
	webView webview.WebView
}

// New creates a new AppView instance.
func New(rawUrl string, options ...Option) *AppView {
	// Apply all user-provided options to override the defaults.
	opts, err := GetDefaultOptions().Apply(options...)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// Create default webview if none was provided.
	if opts.WebView == nil {
		opts.WebView = webview.New(false)
	}

	// Parse URL to figure out if local server should
	// be set up. The local server is recommened during
	// production, so when shipping your app.
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// If the scheme points to a folder. The local HTTP
	// server will be started.
	if parsedUrl.Scheme == "file" {
		webroot := parsedUrl.Path

		// Configure pkger to inline folder.
		_ = pkger.Include(webroot)

		return &AppView{
			URL:     rawUrl,
			webView: opts.WebView,
			server:  NewInlineFilesystem(webroot).Server(),
		}
	}

	return &AppView{
		URL:     rawUrl,
		webView: opts.WebView,
		server:  nil,
	}
}

// Destroy stops the application by shutting down the HTTP server
// if one was started and destroying the WebView.
func (a *AppView) Destroy() {
	if a.server != nil {
		a.server.Close()
	}

	a.webView.Destroy()
}

// Navigate navigates to the provided remote URL
// or the index.html if a local server is used.
func (a *AppView) Navigate() {
	url := a.URL

	// Overwrite URL if local server is used.
	if a.server != nil {
		url = fmt.Sprintf("%s/index.html", a.server.URL)
	}

	a.webView.Navigate(url)
}

// Run blocks until the application is closed.
func (a *AppView) Run() {
	a.webView.Run()
}
