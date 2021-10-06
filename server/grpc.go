// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"net"

	"github.com/vicanso/tiny/log"
	"github.com/vicanso/tiny/pb"
	"github.com/vicanso/tiny/tiny"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	// GRPCServer grpc server
	GRPCServer struct{}
)

// DoOptim do optim
func (gs *GRPCServer) DoOptim(ctx context.Context, in *pb.OptimRequest) (reply *pb.OptimReply, err error) {
	var encodeType tiny.EncodeType
	switch in.Source {
	case pb.Type_JPEG:
		encodeType = tiny.EncodeTypeJPEG
	case pb.Type_PNG:
		encodeType = tiny.EncodeTypePNG
	case pb.Type_WEBP:
		encodeType = tiny.EncodeTypeWEBP
	default:
		encodeType = tiny.EncodeTypeUnknown
	}

	var outputType tiny.EncodeType
	switch in.Output {
	case pb.Type_JPEG:
		outputType = tiny.EncodeTypeJPEG
	case pb.Type_PNG:
		outputType = tiny.EncodeTypePNG
	case pb.Type_AVIF:
		outputType = tiny.EncodeTypeAVIF
	case pb.Type_WEBP:
		outputType = tiny.EncodeTypeWEBP
	case pb.Type_GZIP:
		outputType = tiny.EncodeTypeGzip
	case pb.Type_BR:
		outputType = tiny.EncodeTypeBr
	case pb.Type_SNAPPY:
		outputType = tiny.EncodeTypeSnappy
	case pb.Type_LZ4:
		outputType = tiny.EncodeTypeLz4
	case pb.Type_ZSTD:
		outputType = tiny.EncodeTypeZstd
	default:
		outputType = tiny.EncodeTypeUnknown
	}

	if outputType == tiny.EncodeTypeUnknown {
		err = errOutputTypeIsInvalid
		return
	}
	if len(in.Data) == 0 {
		err = errDataIsNil
		return
	}

	quality := int(in.Quality)

	if encodeType >= tiny.EncodeTypeJPEG && encodeType <= tiny.EncodeTypeWEBP {
		// 图片只能转换为图片
		if outputType < tiny.EncodeTypeJPEG {
			err = errOutputTypeIsInvalid
			return
		}
		crop := tiny.CropType(in.Crop)
		imgInfo, err := tiny.ImageOptim(ctx, in.Data, encodeType, outputType, crop, quality, int(in.Width), int(in.Height))
		if err != nil {
			return nil, err
		}
		reply = &pb.OptimReply{
			Output: in.Output,
			Data:   imgInfo.Data,
			Width:  uint32(imgInfo.Width),
			Height: uint32(imgInfo.Height),
		}
	} else {
		info, err := tiny.TextOptim(in.Data, outputType, quality)
		if err != nil {
			return nil, err
		}
		reply = &pb.OptimReply{
			Data:   info.Data,
			Output: in.Output,
		}
	}

	return
}

// NewGRPCServer new a grpc server
func NewGRPCServer(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterOptimServer(s, &GRPCServer{})
	reflection.Register(s)
	log.Default().Info().
		Str("address", address).
		Msg("grpc server is listening")
	return s.Serve(ln)
}
