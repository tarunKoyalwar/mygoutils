package myhttp

import (
	"io/ioutil"
	"net/http"
)

func GetWebPage(url string) *[]byte {
	resp, err := DefaultClient.Get(url)
	HandleError(err, "Failed to fetch response of "+url)
	data, err2 := ioutil.ReadAll(resp.Body)
	HandleError(err2, "Failed to Get body of "+url)
	return &data
}

func GetWebPageWithHeader(url string, header_name string, header_value string) *[]byte {
	req, err0 := http.NewRequest("GET", url, nil)
	HandleError(err0, "Failed to create Request for "+url)
	req.Header.Add(header_name, header_value)
	resp, err := DefaultClient.Do(req)
	HandleError(err, "Failed to fetch response of "+url)
	data, err2 := ioutil.ReadAll(resp.Body)
	HandleError(err2, "Failed to Get body of "+url)
	return &data
}
