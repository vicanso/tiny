package main

import (
	"flag"
	"fmt"
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

	showVersion bool

	// Version version of tiny
	Version string
	// BuildAt build at
	BuildAt string
	// GO go version
	GO string
)

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
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&httpAddress, "http", "", "http server listen address, eg: 127.0.0.1:7001")
	flag.StringVar(&grpcAddress, "grpc", "", "grpc server listen address, eg: 127.0.0.1:7002")
	// TODO 由于端口可以通过参数指定，检测需要考虑怎么调整
	if len(os.Args) > 1 && os.Args[1] == "check" {
		healthCheck()
		return
	}
	flag.Parse()
	if showVersion {
		fmt.Printf("version %s\nbuild at %s\n%s\n", Version, BuildAt, GO)
		return
	}

	// 如果两个服务都未指定地址，则使用默认地址
	if httpAddress == "" && grpcAddress == "" {
		httpAddress = defaultHTTPAddress
		grpcAddress = defaultGRPCAddress
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
