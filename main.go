package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

const (
	HTTP_PORT  = ":8080"
	HTTPS_PORT = ":8081"
	CERT_FILE  = "server.cert"
	CERT_KEY   = "server.key"
)

/*
	So basically what we are going to build is
	a api gateway and proxy server with load balancing
	without complex configurations. its a start with the
	single things we need to take into account.
*/
func main() {

	upstreamServer, upstreamServerErr := url.Parse("http://localhost:8000")

	if upstreamServerErr != nil {
		log.Println("[contentCopyError]: ", upstreamServerErr)
	}

	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		start := time.Now()

		req.Host = upstreamServer.Host
		req.URL.Host = upstreamServer.Host
		req.URL.Scheme = upstreamServer.Scheme
		req.RequestURI = ""

		remoteAddressHost, _, _ := net.SplitHostPort(req.RemoteAddr)

		req.Header.Set("X-Forwarded-For", remoteAddressHost)

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
		for key, values := range response.Header {
			for _, value := range values {
				rw.Header().Set(key, value)
			}
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
		_, contentCopyError := io.Copy(rw, response.Body)

		//copy the Trailer Values
		for key, values := range response.Trailer {
			for _, value := range values {
				rw.Header().Set(key, value)
			}
		}

		close(proxyDone)
		if contentCopyError != nil {
			log.Println("[contentCopyError]: ", err)
		}

	})

	go func() {

		log.Println("Proxy Listening on http", HTTP_PORT)
		if err := http.ListenAndServe(HTTP_PORT, proxy); err != nil {
			log.Println("ListenAndServe", err)
		}

	}()

	log.Println("Proxy Listening on https", HTTPS_PORT)
	if err := http.ListenAndServeTLS(HTTPS_PORT, CERT_FILE, CERT_KEY, proxy); err != nil {
		log.Println("ListenAndServe", err)
	}

}
