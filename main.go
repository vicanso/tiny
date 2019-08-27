package main

import (
	"github.com/vicanso/tiny/server"
)

func main() {
	go func() {
		err := server.NewHTTPServer(":3015")
		if err != nil {
			panic(err)
		}
	}()

	err := server.NewGRPCServer(":3016")
	if err != nil {
		panic(err)
	}
}
