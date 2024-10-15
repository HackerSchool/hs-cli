package client

import (
	"hscli/logging"
	"net/http"
	"time"
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

// Decorator to log outgoing requests and incoming responses
func (lrt WithLoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	logging.LogDebug("Outgoing request %s %s", r.Method, r.URL)
	start := time.Now()
	rsp, err := lrt.r.RoundTrip(r)
	if err != nil {
		return rsp, err
	}
	logging.LogDebug("Incoming response %d, took %0.2fs", rsp.StatusCode, time.Since(start).Seconds())
	return rsp, err
}
