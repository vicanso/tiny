package main

import (
	"flag"
	"fmt"

	"github.com/vicanso/tiny/server"
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

func main() {
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&httpAddress, "http", "", "http server listen address, eg: 127.0.0.1:7001")
	flag.StringVar(&grpcAddress, "grpc", "", "grpc server listen address, eg: 127.0.0.1:7002")

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
