package shadow

import (
	pb "../compress"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func compress(in *pb.CompressRequest) ([]byte, error) {
	var newBuf []byte
	var err error
	alg := in.Type
	buf := in.Data
	switch alg {
	default:
		newBuf, err = doGzip(buf)
	case pb.Type_BROTLI:
		newBuf, err = doBrotli(buf)
	case pb.Type_WEBP:
		newBuf, err = doWebp(buf, in.Width, in.Height, in.Quality, in.ImageType)
	case pb.Type_JPEG:
		newBuf, err = doJPEG(buf, in.Width, in.Height, in.Quality, in.ImageType)
	}
	if err != nil {
		return nil, err
	}
	return newBuf, nil
}

// server is used to implement compress.CompressServer.
type server struct{}

func (s *server) Do(ctx context.Context, in *pb.CompressRequest) (*pb.CompressReply, error) {
	buf, err := compress(in)
	if err != nil {
		return nil, err
	}
	return &pb.CompressReply{
		Data: buf,
	}, nil
}

// Run 启动GRPC服务
func Run() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterCompressServer(s, &server{})
	reflection.Register(s)
	return s
}
