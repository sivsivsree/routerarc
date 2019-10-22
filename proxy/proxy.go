package proxy

import (
	"context"
	"github.com/sivsivsree/routerarc/balancer"
	"github.com/sivsivsree/routerarc/data"
	"log"
	"net/http"
	"net/url"
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

			//balanceServers := setBalanceServers(proxy.To)

			//var balanceHandler *http.Handler

			//switch proxy.Loadbalacer {
			//case "round-robin":
			//	balanceHandler = balancer.LoadBalancer(proxy.To, proxyFunction)
			//	break;
			//default:
			//	balanceHandler = balancer.GetRoundRobinBalancer(proxy.To, proxyFunction)
			//}
			//
			//if balanceHandler == nil {
			//	log.Println("[ReverseProxyServerUp]", proxy.Name, "Failed")
			//}

			lb := balancer.New("round-robin", proxy.To)
			server := rpServers.startHttpServer(proxy, lb)

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

func setBalanceServers(balanceesStrings []string) []url.URL {
	var balancees = []url.URL{}
	for _, u := range balanceesStrings {
		var purl, _ = url.Parse(u)
		balancees = append(balancees, *purl)
	}
	return balancees
}
