package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/sivsivsree/routerarc/config"
	"github.com/sivsivsree/routerarc/proxy"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func SetUpFlags() {

}

/*
	So basically what we are going to build is
	a api gateway and proxy server with load balancing
	without complex configurations. its a start with the
	single things we need to take into account.
*/
func main() {
	conf := flag.String("config", "rules.json", "Configuration file path")
	flag.Parse()

	fmt.Println("Configuration from", *conf)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Parse the serviceConfig
	serviceConfig, err := config.GetConfig(*conf)
	if err != nil {
		log.Fatal("[configuration]", err)
	}

	rp := proxy.InitReverseProxy()

	// Run only if there is proxy configurations in the config json.
	if serviceConfig.ProxyServiceCount() > 0 {
		rp.SpinProxyServers(serviceConfig.Proxy)
	}

	<-done
	rp.ShutdownProxyServers()

}
