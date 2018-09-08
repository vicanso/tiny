package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc/reflection"

	pb "github.com/vicanso/tiny/proto"

	"github.com/vicanso/tiny/tiny"
	"google.golang.org/grpc"
)

var httpPort = flag.String("httpPort", "3015", "the http listen port")
var grpcPort = flag.String("grpcPort", "3016", "the grpc listen port")

func httpListen() {
	httpServer := tiny.HTTPServer{}
	http.HandleFunc("/ping", httpServer.Ping)
	http.HandleFunc("/optim", httpServer.Optim)
	log.Println("http server listen on:" + *httpPort)
	http.ListenAndServe(":"+*httpPort, nil)
}

func grpcListen() {

	listen, err := net.Listen("tcp", ":"+*grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCompressServer(s, &tiny.GRPCServer{})
	reflection.Register(s)

	log.Println("grp server listen on:" + *grpcPort)
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
	url := "http://127.0.0.1:" + *httpPort + "/ping"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("http check fail: %v", err)
		os.Exit(1)
		return
	}
	statusCode := resp.StatusCode
	if statusCode < 200 || statusCode >= 400 {
		log.Fatalf("http check fail, status:%d", statusCode)
		os.Exit(1)
		return
	}
}

func checkGRPC() {
	conn, err := grpc.Dial("127.0.0.1:"+*grpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc check fail: %v", err)
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

func init() {
	flag.Parse()
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
