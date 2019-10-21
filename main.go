package main

import (
	"crypto/tls"
	"flag"
	"github.com/sivsivsree/routerarc/config"
	"github.com/sivsivsree/routerarc/proxy"
	"log"
	"net/http"
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
	//parse the serviceConfig
	serviceConfig, err := config.GetConfig(*conf)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan bool)
	proxy.SpinProxyServers(serviceConfig.Proxy)

	<-ch
}
