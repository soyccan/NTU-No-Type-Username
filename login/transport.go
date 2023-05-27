package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

const logHTTPMsg = false

// implements http.RoundTripper
type ntuTransport struct{}

func stripDefaultPort(s string) string {
	// Workaround: the web server rejects port number in Host field of HTTP header
	// so we tamper url host by removing default port number
	host, port, cut := strings.Cut(s, ":")
	if cut && (port == "80" || port == "443") {
		return host
	}
	return s
}

// hook before sending data
// log all http requests for debugging
func (*ntuTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.URL.Host = stripDefaultPort(req.URL.Host)

	var logData []byte

	if logHTTPMsg {
		reqData, _err := httputil.DumpRequestOut(req, true)
		if _err != nil {
			err = _err
			return
		}
		logData = append(logData, reqData...)
	}

	resp, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return
	}

	if logHTTPMsg {
		respData, _err := httputil.DumpResponse(resp, true)
		if _err != nil {
			err = _err
			return
		}
		logData = append(logData, respData...)
		fmt.Printf("%s\n", logData)
	}

	return
}
