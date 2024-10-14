package client

import (
	"hscli/logging"
	"net/http"
)

type WithUARoundTripper struct {
	r http.RoundTripper
}

// Decorator to add user-agent to requests
func (urt WithUARoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("User-Agent", ProgramName+"/"+ProgramVersion)
	return urt.r.RoundTrip(r)
}

type WithLoggingRoundTripper struct {
	r http.RoundTripper
}

// Decorator to log outgoing requests
func (lrt WithLoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	logging.LogDebug("Outgoing request %s %s", r.Method, r.URL)
	return lrt.r.RoundTrip(r)
}
