package proxy

import (
	"context"
	"fmt"
	"github.com/sivsivsree/routerarc/balancer"
	"github.com/sivsivsree/routerarc/data"
	"github.com/sivsivsree/routerarc/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

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
func (rpServers *ReverseProxyServers) startHttpServer(proxy data.Proxy, handler http.Handler) *http.Server {

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

			if proxy.Static != "" {

				fmt.Println(proxy)
				fileServer := http.FileServer(utils.FileSystem{FS: http.Dir(proxy.Static)})

				//http.Handle("/", http.StripPrefix(strings.TrimRight("/", "/"), fileServer))
				handler := http.NewServeMux()
				handler.Handle("/", http.StripPrefix(strings.TrimRight("/", "/"), fileServer))
				log.Println("using [index.html] as default entrypoint for", proxy.Name)
				staticServer := rpServers.startHttpServer(proxy, handler)
				rpServers.addActiveServers(proxy, staticServer)
			}

			if proxy.Loadbalacer != "" && proxy.Static == "" {
				// make the balance handle pass through here
				lb := balancer.New(proxy.Loadbalacer, proxy.To)
				server := rpServers.startHttpServer(proxy, lb)
				rpServers.addActiveServers(proxy, server)

			}

		}(proxyValue)

	}

}

func (rpServers *ReverseProxyServers) addActiveServers(proxy data.Proxy, server *http.Server) {
	rpServers.ActiveServers = append(rpServers.ActiveServers, struct {
		Name   string
		Server *http.Server
	}{Name: proxy.Name, Server: server})

	log.Println("[ReverseProxyServerUp]", proxy.Name, "reverse proxy serving on port", server.Addr)
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
