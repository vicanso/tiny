package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	shadow "./shadow"
)

const (
	httpPort = ":3015"
	port     = ":3016"
)

func httpListen() {
	shadow.RunHTTP()
	log.Println("http server listen on" + httpPort)
	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatal(err)
	}
}

func grpcListen() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := shadow.RunGRPC()
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
