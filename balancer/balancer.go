package balancer

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// LB is used to keep the current
// load balancing configurations.
type LB struct {
	algorithm   string
	endpoint    []string
	LBAlgorithm LBAlgorithm
}

//noinspection GoNilness
func (lb *LB) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	start := time.Now()

	url, err := lb.LBAlgorithm.Apply()

	//fmt.Println(url)
	if err != nil {
		log.Println(err)
	}

	req.Host = url.Host
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme

	req.RequestURI = ""

	remoteAddressHost, _, _ := net.SplitHostPort(req.RemoteAddr)
	req.Header.Set("X-Forwarded-For", remoteAddressHost)
	req.Header.Set("X-Forwarded-Host", remoteAddressHost)
	// To enable HTTP 2 protocol
	//_ = http2.ConfigureTransport(http.DefaultTransport.(*http.Transport))

	// sending the request from proxy to upstreamServer
	response, err := http.DefaultClient.Do(req)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(rw, err.Error())
		return
	}

	rw.Header().Set("routerarc-proxy-time", time.Since(start).String())

	//copy the Headers
	//for key, values := range response.Header {
	//	for _, value := range values {
	//		fmt.Println("2 loop:", key, value)
	//		rw.Header().Set(key, value)
	//	}
	//}

	for k, v := range response.Header {
		rw.Header().Set(k, strings.Join(v, ";"))
	}

	// flush the output for Stream Channels
	proxyDone := make(chan bool)
	go func() {
		for {
			select {
			case <-time.Tick(10 * time.Millisecond):
				rw.(http.Flusher).Flush()
			case <-proxyDone:
				return
			}
		}
	}()

	//trailers
	trailerKeys := []string{}

	for key := range response.Trailer {
		trailerKeys = append(trailerKeys, key)
	}

	rw.Header().Set("Trailer", strings.Join(trailerKeys, ","))

	//copy the StatusCode
	rw.WriteHeader(response.StatusCode)

	// Copy the Content

	if _, contentCopyError := io.Copy(rw, response.Body); contentCopyError != nil {

		log.Println("[contentCopyError]: ", err)
	}

	defer response.Body.Close()

	//copy the Trailer Values
	for key, values := range response.Trailer {
		for _, value := range values {
			rw.Header().Set(key, value)
		}
	}

	close(proxyDone)
	fmt.Println(time.Since(start).String())
}

// New is used to create new Load balancer with the given algorithm.
func New(algorithm string, endpoints []string) *LB {

	var algo LBAlgorithm

	switch algorithm {

	case ROUND_ROBIN:
		algo = NewRoundRobin(endpoints)
		break
	default:
		algo = NewRoundRobin(endpoints)
	}

	return &LB{
		algorithm:   algorithm,
		endpoint:    endpoints,
		LBAlgorithm: algo,
	}
}
