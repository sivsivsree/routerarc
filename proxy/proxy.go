package proxy

import (
	"context"
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
	upstreams []*url.URL
	current   *url.URL
	count     int
}

func (upstreamServers *UpstreamServers) setUpstreamServer(proxy data.Proxy) {
	upstreamServers.current, _ = url.Parse(proxy.To[0])

}

func (upstreamServers UpstreamServers) proxyFunction(rw http.ResponseWriter, req *http.Request) {

	start := time.Now()

	req.Host = upstreamServers.current.Host
	req.URL.Host = upstreamServers.current.Host
	req.URL.Scheme = upstreamServers.current.Scheme
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

	defer response.Body.Close()

	//copy the Trailer Values
	for key, values := range response.Trailer {
		for _, value := range values {
			rw.Header().Set(key, value)
		}
	}

	close(proxyDone)

}

// ReverseProxyServer contains all the
// ActiveServers used for Reverse proxying.
type ReverseProxyServers struct {
	ActiveServers []struct {
		Name   string
		Server *http.Server
	}
}

// InitReverseProxy will be used to return the ReverseProxyServers instance.
func InitReverseProxy() ReverseProxyServers {
	return ReverseProxyServers{}
}

// startHttpServer the method is used to call UpstreamServers to make Proxy Servers
func (rpServers *ReverseProxyServers) startHttpServer(proxy data.Proxy, handler http.HandlerFunc) *http.Server {

	server := &http.Server{Addr: ":" + proxy.Port, Handler: handler}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("[ReverseProxyServerUp]", proxy.Name, "failed,", err)
		}
	}()
	return server
}

// SpinProxyServers is used to create servers based on the ports specified for
// reverse proxying, will be passing the proxy[] from the configuration file.
func (rpServers *ReverseProxyServers) SpinProxyServers(proxies []data.Proxy) {

	for _, proxyValue := range proxies {

		go func(proxy data.Proxy) {

			upServer := &UpstreamServers{}
			upServer.setUpstreamServer(proxy)
			server := rpServers.startHttpServer(proxy, upServer.proxyFunction)

			rpServers.ActiveServers = append(rpServers.ActiveServers, struct {
				Name   string
				Server *http.Server
			}{Name: proxy.Name, Server: server})

			log.Println("[ReverseProxyServerUp]", proxy.Name, "reverse proxy serving on port", server.Addr)

		}(proxyValue)

	}

}

// ShutdownProxyServers is used to gracefully shutdown all the currently
// running services.
func (rpServers *ReverseProxyServers) ShutdownProxyServers() {

	for _, rp := range rpServers.ActiveServers {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			// extra handling here
			cancel()
		}()
		if err := rp.Server.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}

		log.Println("[ShutdownProxyServers]", rp.Server.Addr)

	}

}
