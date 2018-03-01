package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/vicanso/tiny/shadow"
	"google.golang.org/grpc"
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

func contains(list []string, s string) bool {
	for _, item := range list {
		if s == item {
			return true
		}
	}
	return false
}

func checkHTTP() {
	url := "http://127.0.0.1" + httpPort + "/ping"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
		return
	}
	statusCode := resp.StatusCode
	if statusCode < 200 || statusCode >= 400 {
		fmt.Print(err)
		os.Exit(1)
		return
	}
}

func checkGRPC() {
	conn, err := grpc.Dial("127.0.0.1:"+port, grpc.WithInsecure())
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
		return
	}
	conn.Close()
}

func check(services string) {
	if len(services) == 0 {
		checkHTTP()
		checkGRPC()
	} else if strings.Compare(services, "HTTP") == 0 {
		checkHTTP()
	} else {
		checkGRPC()
	}
	os.Exit(0)
}

func main() {
	services := strings.ToUpper(os.Getenv("SERVICES"))

	if contains(os.Args[1:], "check") {
		check(services)
		return
	}
	if len(services) == 0 {
		go httpListen()
		grpcListen()
	} else if strings.Compare(services, "HTTP") == 0 {
		httpListen()
	} else {
		grpcListen()
	}
}
