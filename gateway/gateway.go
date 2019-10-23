package gateway

import (
	"context"
	"github.com/sivsivsree/routerarc/data"
	"net/http"
	"time"
)

type ApiGatewayServers struct {
	ActiveServers []struct {
		Name   string
		Server *http.Server
	}
}

func InitApiGatewayServer() ApiGatewayServers {

	return ApiGatewayServers{}
}

func (gws ApiGatewayServers) startGatewayServer(router data.Router) *http.Server {
	//server := &http.Server{Addr: ":" + gws.Port, Handler: handler}
	//
	//go func() {
	//	if err := server.ListenAndServe(); err != nil {
	//		log.Println("[startGatewayServerSetup]", "failed,", err)
	//	}
	//}()
	//return server
	return nil

}

func (gws ApiGatewayServers) SpinGatewayServer(routes []data.Router) {

	for _, server := range routes {

		gws.startGatewayServer(server)

	}

}

// ShutdownProxyServers is used to gracefully shutdown all the currently
// running services.
func (gws *ApiGatewayServers) ShutdownGatewayServer() {

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	//if err := rp.Server.Shutdown(ctx); err != nil {
	//	log.Fatalf("Server Shutdown Failed:%+v", err)
	//}

	//log.Println("[ShutdownGatewayServer]", rp.Server.Addr)

}
