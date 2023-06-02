package main

import (
	"net/http"
	"net/http/cookiejar"
)

type Loginable interface {
	Login() error
}

func Login(l Loginable) error {
	return l.Login()
}

// implements NoTypeUsernameLoginner
type NoTypeUsername struct {
	CredPath   string
	CookiePath string
	client     *http.Client
	noRedirect bool
}

func (ntu *NoTypeUsername) initHTTPClient() error {
	if ntu.client != nil {
		return nil
	}
	ntu.client = &http.Client{
		Jar: func(j *cookiejar.Jar, _ error) *cookiejar.Jar { return j }(cookiejar.New(nil)),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if ntu.noRedirect {
				// Disable automatic redirect by returning http.ErrUseLastResponse
				ntu.noRedirect = false
				return http.ErrUseLastResponse
			}
			return nil
		},
		Transport: &ntuTransport{},
	}
	if ntu.client == nil {
		return &ClientInitError{}
	}
	return nil
}

type ClientInitError struct{}

func (*ClientInitError) Error() string { return "HTTP client init failed" }
