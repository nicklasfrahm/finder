package main

import (
	"github.com/nicklasfrahm/finder/pkg/appview"
)

// url can be overriden during development
// to use the development server of a React
// or Svelte app. See the Makefile for more
// information.
var url = "file:///web/build"

func main() {
	// Wrap the creation of the webview
	// to work around some quirks on WSL.
	webView := NewWebView(false)

	app := appview.New(url, appview.WithWebView(webView))
	defer app.Destroy()

	app.Navigate()

	app.Run()
}
