package main

import (
	"crypto/tls"
	"flag"
	"github.com/sivsivsree/routerarc/config"
	"github.com/sivsivsree/routerarc/gateway"
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
*
*	So basically what we are going to build is
*	a api gateway and proxy server with load balancing
*	without complex configurations. its a start with the
*	single things we need to take into account.
*
 */
func main() {

	conf := flag.String("config", "rules.yaml", "This flag is used for specifying the configuration file.")
	json := flag.Bool("json", false, "If the json flag is true routerarc will read json file for configuration,\nalso you need to to specify the path of the json file using -config flag")
	flag.Parse()

	log.Println("[configuration] Configuration from", *conf)

	// for gracefull shutdown of service.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Parse the serviceConfig
	serviceConfig, err := config.GetConfig(*conf, *json)
	if err != nil {
		log.Fatal("[configuration]", err)
	}

	// Initialize the APIRouter
	gwApi := gateway.InitApiGatewayServer()
	// Run only if there is Router configurations in the configuration.
	if serviceConfig.RouterServiceCount() > 0 {
		gwApi.SpinGatewayServer(serviceConfig.Router)
	}

	// Initialize the ReverseProxy Server
	rp := proxy.InitReverseProxy()
	// Run only if there is proxy configurations in the configuration.
	if serviceConfig.ProxyServiceCount() > 0 {
		rp.SpinProxyServers(serviceConfig.Proxy)
	}

	<-done
	rp.ShutdownProxyServers()

}
