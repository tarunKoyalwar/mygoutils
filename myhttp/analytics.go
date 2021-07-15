package myhttp

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// This is a function to measure response time of request
// Note this works good only in single thread and is not thread safe
// do not use this with goroutines (Even if you want try threshold > 40)
// inputs
// client -> a http client instance,
// req  -> A Http request to send,
// duration -> The Threshold value to get response,
// returns
// a boolean value depending on the threshold
func CompareRespTime(client *http.Client, req *http.Request, duration int) bool {
	start := time.Now()
	_, err := client.Do(req)
	HandleError(err, "Request Sending Failed => "+req.RequestURI)
	eval := time.Since(start).Seconds()
	if eval > float64(duration) {
		return true
	} else {
		return false
	}

}

// This is a function to measure response time of request
// Note this works good only in single thread and is not thread safe
// do not use this with goroutines (Even if you want try threshold > 40)
// inputs
// client -> a http client instance,
// req  -> A Http request to send,
// This returns response time of the request
func GetRespTime(client *http.Client, req *http.Request) float64 {
	start := time.Now()
	_, err := client.Do(req)
	HandleError(err, "Request Sending Failed => "+req.RequestURI)
	eval := time.Since(start).Seconds()
	return float64(eval)
}

// This function is used to create Urls from base and relative paths
// It does not matter if it has / or not it adds itself
// It also automatically handles schema and defaults to https
// inputs
// site -> base path,
// path -> relative path
func CreateURL(site string, path string) string {
	if strings.HasPrefix(site, "http://") || strings.HasPrefix(site, "https://") {
		base, err1 := url.Parse(site)
		HandleError(err1, "Could not Parse url -> "+site)
		new, err2 := base.Parse(path)
		HandleError(err2, "Could not parse relative url ->"+path)
		return new.String()
	} else {
		base, err1 := url.Parse("https://" + site)
		HandleError(err1, "Could not Parse url -> "+site)
		new, err2 := base.Parse(path)
		HandleError(err2, "Could not parse relative url ->"+path)
		return new.String()
	}
}

//This function is used default schema of urls to https
// It handles schema adds if not  present
// Along with normal it also remove port 80 and replaces it with 443
func ForceHTTPs(urls *[]string) {
	for i, v := range *urls {
		if strings.Contains(v, "http:") {
			// fmt.Println(v)
			(*urls)[i] = strings.Replace(v, "http:", "https:", -1)
		} else if strings.Contains(v, "https:") {
			_ = i
			continue
		} else {
			(*urls)[i] = "https://" + v
		}
	}
	for i, v := range *urls {
		if strings.Contains(v, ":80") {
			(*urls)[i] = strings.Replace(v, ":80", ":443", -1)
		}
	}
}

//Just a Wrapper function to Handle Error
// with message
func HandleError(er error, msg string) {
	if er != nil {
		fmt.Println(msg)
		panic(er)
	}
}
