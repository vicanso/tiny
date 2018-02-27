package main

import (
	"io/ioutil"
	"log"

	pb "github.com/vicanso/tiny/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	buf, err := ioutil.ReadFile("../assets/lodash.min.js")
	if err != nil {
		log.Fatalf("get file data: %v", err)
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCompressClient(conn)

	res, err := c.Do(context.Background(), &pb.CompressRequest{
		Type: pb.Type_GZIP,
		Data: buf,
	})
	if err != nil {
		log.Fatalf("could not do gzip: %v", err)
	}
	log.Printf("gzip success, reduce %d percent", 100-(len(res.Data)*100)/len(buf))

	res, err = c.Do(context.Background(), &pb.CompressRequest{
		Type: pb.Type_BROTLI,
		Data: buf,
	})
	if err != nil {
		log.Fatalf("could not do brotli: %v", err)
	}
	log.Printf("brotli success, reduce %d percent", 100-(len(res.Data)*100)/len(buf))

	buf, err = ioutil.ReadFile("../assets/banner.png")
	if err != nil {
		log.Fatalf("get file data: %v", err)
	}
	res, err = c.Do(context.Background(), &pb.CompressRequest{
		Type:      pb.Type_WEBP,
		Data:      buf,
		ImageType: "png",
		Quality:   75,
	})

	if err != nil {
		log.Fatalf("could not do webp: %v", err)
	}
	log.Printf("webp success, reduce %d percent", 100-(len(res.Data)*100)/len(buf))

	res, err = c.Do(context.Background(), &pb.CompressRequest{
		Type:      pb.Type_JPEG,
		Data:      buf,
		ImageType: "png",
		Quality:   75,
	})

	if err != nil {
		log.Fatalf("could not do jpeg: %v", err)
	}
	log.Printf("jpeg success, reduce %d percent", 100-(len(res.Data)*100)/len(buf))
}
