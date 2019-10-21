package proxy

import (
	"fmt"
	"github.com/sivsivsree/routerarc/data"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type UpstreamServers struct {
	upstreamServer *url.URL
}

func (servers *UpstreamServers) setUpstreamServer(currentServer string) {
	servers.upstreamServer, _ = url.Parse(currentServer)
}

func (servers UpstreamServers) proxyFunction(rw http.ResponseWriter, req *http.Request) {

	start := time.Now()

	req.Host = servers.upstreamServer.Host
	req.URL.Host = servers.upstreamServer.Host
	req.URL.Scheme = servers.upstreamServer.Scheme
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

	if _, contentCopyError := io.Copy(rw, response.Body); contentCopyError != nil {

		log.Println("[contentCopyError]: ", err)
	}
	//copy the Trailer Values
	for key, values := range response.Trailer {
		for _, value := range values {
			rw.Header().Set(key, value)
		}
	}

	close(proxyDone)

}

func startHttpServer(port string, handler http.HandlerFunc) *http.Server {

	server := &http.Server{Addr: ":" + port, Handler: handler}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			// handle err
		}
	}()

	return server
}

func SpinProxyServers(proxies []data.Proxy) {

	for _, proxy := range proxies {

		go func(proxy data.Proxy) {
			userver := &UpstreamServers{}
			userver.setUpstreamServer(proxy.To[0])
			serv := startHttpServer(proxy.Port, userver.proxyFunction)

			log.Println(proxy.Name, "proxy running", serv.Addr)
		}(proxy)

	}

}
