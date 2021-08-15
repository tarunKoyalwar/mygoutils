package pulse

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

type Pulse struct {
	Re      *regexp.Regexp
	Timeout time.Duration
}

func NewPulse() *Pulse {
	return &Pulse{
		Re:      regexp.MustCompile(`//(.*?)/`),
		Timeout: time.Duration(5) * time.Second,
	}
}

func (P *Pulse) stat(hostname string) bool {
	var hostaddress string
	if strings.Contains(hostname, ":") {
		hostaddress = hostname
	} else {
		hostaddress = hostname + ":80"
	}

	_, er1 := net.DialTimeout("tcp", hostaddress, P.Timeout)
	if er1 != nil {
		return false
	} else {
		return true
	}
}

func (P *Pulse) GetStatus(url string) bool {
	t := P.Re.FindStringSubmatch(url)
	result := P.stat(t[1])
	if result {
		fmt.Printf("%v is Alive\n", url)
	} else {
		fmt.Printf("%v is Dead\n", url)
	}
	return result
}

func (P *Pulse) Getstats(list []string) *[]bool {
	arr := make([]bool, len(list))
	for k, v := range list {
		arr[k] = P.GetStatus(v)
	}

	return &arr
}
