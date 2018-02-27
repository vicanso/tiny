package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/vicanso/tiny/shadow"
)

const (
	httpPort = ":3015"
	port     = ":3016"
)

func httpListen() {
	log.Println("http server listen on" + httpPort)
	fasthttp.ListenAndServe(httpPort, shadow.HTTPHandler)
}

func grpcListen() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := shadow.GetGRPCServer()
	log.Println("grp server listen on" + port)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	services := strings.ToUpper(os.Getenv("SERVICES"))
	if len(services) == 0 {
		go httpListen()
		grpcListen()
	} else if strings.Compare(services, "HTTP") == 0 {
		httpListen()
	} else {
		grpcListen()
	}
}
