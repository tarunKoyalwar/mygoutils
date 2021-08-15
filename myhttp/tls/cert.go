package tls

import (
	"crypto/tls"
	"strconv"
)

//Fetches Common Name From TLS Certificate of given Hostname
// if function returns 0 then it could not grab certificate
func GetCN(hostname string) string {
	conn, er1 := tls.Dial("tcp", hostname+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if er1 != nil {
		return "0"
	}
	dbyte := conn.ConnectionState().PeerCertificates[0].Subject
	// fmt.Println(dbyte)
	return dbyte.CommonName
}

//Fetches Common Name From TLS Certificate of given Hostname and PORT
// if function returns 0 then it could not grab certificate
func GetCNWithPort(hostname string, port int) string {
	conn, er1 := tls.Dial("tcp", hostname+":"+strconv.Itoa(port), &tls.Config{
		InsecureSkipVerify: true,
	})
	if er1 != nil {
		return "0"
	}
	dbyte := conn.ConnectionState().PeerCertificates[0].Subject
	// fmt.Println(dbyte)
	return dbyte.CommonName
}

//Performs BULK Extraction of CN (COMMON NAME)
func GetBulkCN(hostnames *[]string) *[]string {
	arr := make([]string, len(*hostnames))
	for k, v := range *hostnames {
		arr[k] = GetCN(v)
	}
	return &arr
}
