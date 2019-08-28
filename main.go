package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vicanso/tiny/server"
	"google.golang.org/grpc"
)

const (
	defaultHTTPAddress = ":7001"
	defaultGRPCAddress = ":7002"
)

var (
	httpAddress string
	grpcAddress string
)

func init() {
	httpAddress = os.Getenv("HTTP")
	grpcAddress = os.Getenv("GRPC")
	// 如果两个服务都未指定地址，则使用默认地址
	if httpAddress == "" && grpcAddress == "" {
		httpAddress = defaultHTTPAddress
		grpcAddress = defaultGRPCAddress
	}
}

func checkHTTP() {
	url := "http://127.0.0.1" + httpAddress + "/ping"
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
	conn, err := grpc.Dial("127.0.0.1"+grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc check fail: %v", err)
		os.Exit(1)
		return
	}
	conn.Close()
}

func healthCheck() {
	if httpAddress != "" {
		checkHTTP()
	}
	if grpcAddress != "" {
		checkGRPC()
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "check" {
		healthCheck()
		return
	}

	if httpAddress != "" && grpcAddress != "" {
		go func() {
			err := server.NewHTTPServer(httpAddress)
			if err != nil {
				panic(err)
			}
		}()

		err := server.NewGRPCServer(grpcAddress)
		if err != nil {
			panic(err)
		}
		return
	}
	if httpAddress != "" {
		err := server.NewHTTPServer(httpAddress)
		if err != nil {
			panic(err)
		}
		return
	}
	err := server.NewGRPCServer(grpcAddress)
	if err != nil {
		panic(err)
	}
}
