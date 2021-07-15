// This is a different program
// This program has methods to serailze and deserialize request
// In lame terms we can create a request just by copying the plain burpsuite editor data
// The important thing is it does not matter whether it is GET/POST or has body etc
// Note the url is constructed using Host header and relative path
// If you want to spoof host header use different fuctions
//Ex:
// POST /cart/checkout HTTP/1.1
// Host: ac591f241ece64df80d508f2008900a7.web-security-academy.net
// Cookie: session=MHoKMUS1Jrq2mFisuar8SRnrV5R7QgRW
// Content-Length: 37
// Cache-Control: max-age=0
// Upgrade-Insecure-Requests: 1
// Origin: https://ac5c1f981eff870a800e417c001d0091.web-security-academy.net
// Content-Type: application/x-www-form-urlencoded
// User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
// Sec-Gpc: 1
// Sec-Fetch-Site: same-origin
// Sec-Fetch-Mode: navigate
// Sec-Fetch-User: ?1
// Sec-Fetch-Dest: document
// Referer: https://ac5c1f981eff870a800e417c001d0091.web-security-academy.net/cart
// Accept-Encoding: gzip, deflate
// Accept-Language: en-GB,en-US;q=0.9,en;q=0.8
// Connection: close

// csrf=MBw8zEPAnPLXok64Fn9QzX4UPRY0Igpc

package myhttp

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Raw_request struct {
	// This is bare metal structure of a http request
	// which is extracted from bytes / string
	// Note  the url is constructed using Host header and relative path
	// There is no need to use this struct unless you want to make changes to request
	// after it is converted from bytes to struct
	// These values are filled from bytestostruct function
	// There are two methods to parse cookies in golang
	// If we use *http.cookie it will occasionally had "" around the cookie
	// if it has sepecial characters
	// Another way is adding it like a header instead of cookie
	Verb         string            // HTTP verb
	URI          string            // This is relative path
	baseuri      *url.URL          // This is base uri constructed from host header
	cookies      *[]http.Cookie    // This contains cookies
	UserAgent    string            // User Agent of request
	Headers      map[string]string // A map with headers and its values forbidden headers are omitted
	Body         []byte            // Byte array of the body
	Host         string            // A Host header
	CookieSwitch bool              // Cookie Switch to switch between cookie parsing methods
}

// This methods are not implemented and are automatically handled by golang
var ForbiddenHeaders string = " Connection Content-Length Transfer-Encoding Trailer "

// This is internal function to seperate request headers and body
func getreqandbody(testx []byte) *[][]byte {
	sep1 := []byte{13, 10, 13, 10}
	sep := []byte{10, 10}
	res1 := bytes.SplitN(testx, sep1, 2)
	if len(res1) < 2 {
		res1 = nil
		res2 := bytes.SplitN(testx, sep, 2)
		return &res2
	}
	return &res1
}

// This is intenal function to fetch cookies from bytes
// This is only helful if you want to parse cookies using *http.cookie
func getCookiesfromLine(line string, raw *Raw_request) {
	cookies := strings.Split(line, ";")
	var cookiekv []http.Cookie
	for _, v := range cookies {
		rawcookie := strings.SplitN(v, "=", 2)
		rawcookie[1] = strings.Trim(rawcookie[1], "\"")
		name := strings.TrimSpace(rawcookie[0])
		value := strings.TrimSpace(rawcookie[1])
		cookiekv = append(cookiekv, http.Cookie{
			Name:  name,
			Value: value,
		})
	}
	raw.cookies = &cookiekv
}

// A simple function to Dump Request using httputil
func RequesttoBytes(req *http.Request) *[]byte {
	dump, err := httputil.DumpRequestOut(req, true)
	HandleError(err, "Failed to dump request")
	return &dump

}

// A simple function to DUmp response using httputil with option
// to dump body or not
func ResponsetoBytes(resp *http.Response, body bool) *[]byte {
	dump, err := httputil.DumpResponse(resp, body)
	HandleError(err, "Failed to dump Response")
	return &dump
}

// This function is helful to Create Struct from a http request
func SerializeRequest(req *http.Request) *Raw_request {
	data := RequesttoBytes(req)
	mystruct := BytestoStruct(data, "")
	return mystruct
}

// This function is helpful to Create a http request from bare metal struct
func DeserializeRequest(z Raw_request) *http.Request {
	urlstring, err1 := z.baseuri.Parse(z.URI)
	HandleError(err1, "Failed to parse Complete URL")

	request, err2 := http.NewRequest(z.Verb, urlstring.String(), bytes.NewReader(z.Body))
	HandleError(err2, "Failed to Create New Request")

	request.Host = z.Host
	request.Header.Set("User-Agent", z.UserAgent)

	//set remaining headers
	for k, v := range z.Headers {
		request.Header.Set(k, v)
	}

	//set Cookies
	if z.CookieSwitch {
		for _, v := range *z.cookies {
			request.AddCookie(&v)
		}
	}

	return request
}

// This converts bytes to http struct .
// In Golang cookies are parsed in two ways,
// If we parse it using http.cookies duoble quotes are added to cookies if it has other symbols in it.
// If we parse it as header it is parsed as it is no change will take place.
// If CookieSwitch = true it is parsed as a  cookie and not as a header.
// input a url string without schema to use as baseuri instead constructing from hostline
func BytestoStruct(reqdata *[]byte, urlx string) *Raw_request {

	var z Raw_request
	headerx := make(map[string]string)

	reqbytes := getreqandbody(*reqdata)
	z.Body = (*reqbytes)[1]

	//split line by line
	mystring := string((*reqbytes)[0])
	mystring = strings.ReplaceAll(mystring, "\r", "")
	head := strings.Split(mystring, "\n")

	//parse request verb and uri
	header := strings.Fields(head[0])
	z.Verb = header[0]
	z.URI = header[1]
	hostline := strings.Split(head[1], ":")
	hostheader := strings.TrimSpace(hostline[1])

	z.Host = hostheader
	if urlx == "" {
		base, err := url.Parse("https://" + hostheader)
		HandleError(err, "Failed to parse host line")
		z.baseuri = base
	} else {
		base, err := url.Parse("https://" + urlx)
		HandleError(err, "Failed to parse url from urlx string")
		z.baseuri = base
	}

	//Extract remaining Headers
	head = head[2:]
	for _, v := range head {
		line := strings.SplitN(v, ":", 2)
		if z.CookieSwitch && strings.Contains(line[0], "Cookie") {
			getCookiesfromLine(line[1], &z)
			continue
		} else if strings.Contains(line[0], "User-Agent") {
			z.UserAgent = line[1]
			continue
		} else if strings.Contains(ForbiddenHeaders, line[0]) {
			continue
		} else {
			headerx[line[0]] = line[1]
		}
	}

	z.Headers = headerx

	return &z
}

//This function it used to get a request from bytes
func BytestoRequest(reqdata *[]byte) *http.Request {
	mystruct := BytestoStruct(reqdata, "")
	myrequest := DeserializeRequest(*mystruct)
	return myrequest
}

//This function it used to get a request from bytes
// and uses urlx to construct base uri instead of host header
func BytestoRequestwithurl(reqdata *[]byte, urlx string) *http.Request {
	mystruct := BytestoStruct(reqdata, urlx)
	myrequest := DeserializeRequest(*mystruct)
	return myrequest
}
