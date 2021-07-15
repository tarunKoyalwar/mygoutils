//This program consists of wrapper to efficiently create
// A HTTP client with simple and lame terms rather than RFC standard names
package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type CustomHttpClient struct {
	// This is a struct of http client
	// This consists of simple http client configurations
	// In plain english words
	ValidateCertificate   bool   // Whether to validate ssl certificate
	FollowRedirect        bool   // Whether to follow redirects all 30x
	NewConForEveryRequest bool   // Whether to make new connection for every request This has a tcp handshake overhead if true
	MaxIdleConnections    int    // Maximum number of IDLE connections
	IdleConnectionTimeout int    // Waiting time before a connection is made idle
	Timeout               int    // A normal http connection timeout
	ProxyUrl              string // Proxy URL ex http://127.0.0.1:8080
}

//This is just a function to create appropriate client based
// on the customhttpclient parameters
//This returns a *http.client instance
func (C *CustomHttpClient) Create() *http.Client {

	if C.MaxIdleConnections == 0 {
		C.MaxIdleConnections = 30
	}
	if C.IdleConnectionTimeout == 0 {
		C.IdleConnectionTimeout = 1
	}
	if C.Timeout == 0 {
		C.Timeout = 30
	}

	tempfunc := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	var proxyur func(*http.Request) (*url.URL, error)
	tlsconfig := &tls.Config{InsecureSkipVerify: true}

	if C.FollowRedirect {
		tempfunc = nil
	}
	if C.ValidateCertificate {
		tlsconfig = nil
	}
	if C.ProxyUrl == "" {
		fmt.Println("Not using Proxy")
	} else {
		proxyurl, err := url.Parse(C.ProxyUrl)
		HandleError(err, "Failed to Parse Url")
		proxyur = http.ProxyURL(proxyurl)
	}

	m := &http.Client{
		Timeout:       time.Duration(C.Timeout) * time.Second,
		CheckRedirect: tempfunc,
		Transport: &http.Transport{
			MaxIdleConns:      C.MaxIdleConnections,
			IdleConnTimeout:   time.Duration(C.IdleConnectionTimeout) * time.Second,
			DisableKeepAlives: C.NewConForEveryRequest,
			TLSClientConfig:   tlsconfig,
			Proxy:             proxyur,
		},
	}
	return m

}
