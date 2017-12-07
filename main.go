package main

import (
	"log"
	"net"
	shadow "./shadow"
)

const (
	port = ":50051"
)

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := shadow.Run()
	log.Println("grp server listen on" + port)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}