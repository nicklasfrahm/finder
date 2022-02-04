package main

// This implements a patch to fix some inconsistencies
// when running the application via WSL.

/*
#cgo linux pkg-config: gtk+-3.0

#if defined(__APPLE__) || defined(_WIN32)

void* window_new(void) {
  // Create no-op function on other OSes.
  return NULL;
}

#else

#include <gtk/gtk.h>

void* window_new(void) {

  // Initialize GTK to load user theme.
  gtk_init(NULL, NULL);

  // Create a new header bar to ensure that minimize, maximize and close button are visible.
  GtkWidget* header_bar = gtk_header_bar_new();
  gtk_header_bar_set_show_close_button(GTK_HEADER_BAR(header_bar), TRUE);
  gtk_header_bar_set_decoration_layout(GTK_HEADER_BAR(header_bar), "icon:minimize,maximize,close");

  // Create a new window and configure it.
  GtkWidget* window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
  gtk_window_set_titlebar(GTK_WINDOW(window), header_bar);
  gtk_window_set_default_size(GTK_WINDOW(window), 1280, 720);
  gtk_window_maximize(GTK_WINDOW(window));

  return (void*)window;
}

#endif
*/
import "C"

import (
	"runtime"

	"github.com/webview/webview"
)

// NewWebView creates a new maximized window.
func NewWebView(debug bool) webview.WebView {
	// TODO: I could not make this work with build
	// tags, so for now I will use this workaround.
	// The disadvantage here is increased binary size
	// and most likely also runtime overhead.
	if runtime.GOOS == "linux" {
		window := C.window_new()
		return webview.NewWindow(debug, window)
	}

	return webview.New(debug)
}
