package appview

import (
	"errors"

	"github.com/webview/webview"
)

// Options define the configuration of an AppView.
type Options struct {
	WebView webview.WebView
}

// Apply applies the option functions to the current set of options.
func (o *Options) Apply(options ...Option) (*Options, error) {
	for _, option := range options {
		if err := option(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

// Option defines the function signature to set
// an option for the operations of this library.
type Option func(*Options) error

// GetDefaultOptions returns the default options
// for all operations of this library.
func GetDefaultOptions() *Options {
	return &Options{
		WebView: nil,
	}
}

// WithWebView can be used to provide a custom
// webview instead of using the default.
func WithWebView(webView webview.WebView) Option {
	return func(o *Options) error {
		if webView == nil {
			return errors.New("no webview provided")
		}
		o.WebView = webView
		return nil
	}
}
