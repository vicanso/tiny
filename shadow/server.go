package shadow

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "../compress"
	"google.golang.org/grpc/reflection"
)



func compress(buf []byte, category string) ([]byte, error) {
	// buf, err := ioutil.ReadFile(file)
	// if err != nil {
	// 	return nil, err
	// }
	var newBuf []byte
	var err error
	switch category {
	default:
		newBuf, err = doGzip(buf)
	case "brotli":
		newBuf, err = doBrotli(buf)
	}
	if err != nil {
		return nil, err
	}
	return newBuf, nil
}

// server is used to implement compress.CompressServer.
type server struct{}

func (s *server) Do(ctx context.Context, in *pb.CompressRequest) (*pb.CompressReply, error) {
	return &pb.CompressReply{
		Type: pb.DataType_TEXT,
		Data: []byte("Here is a string...."),
	}, nil
}

func Run() * grpc.Server {
	s := grpc.NewServer()
	pb.RegisterCompressServer(s, &server{})
	reflection.Register(s)
	return s;
}