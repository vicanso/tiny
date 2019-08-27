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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/tiny/pb"
)

func TestDoOptim(t *testing.T) {
	gs := &GRPCServer{}
	t.Run("output type is invalid", func(t *testing.T) {
		assert := assert.New(t)
		req := &pb.OptimRequest{}
		ctx := context.Background()
		reply, err := gs.DoOptim(ctx, req)
		assert.Equal(errOutputTypeIsInvalid, err)
		assert.Nil(reply)
	})

	t.Run("data is nil", func(t *testing.T) {
		assert := assert.New(t)
		req := &pb.OptimRequest{
			Output: pb.Type_GZIP,
		}
		ctx := context.Background()
		reply, err := gs.DoOptim(ctx, req)
		assert.Equal(errDataIsNil, err)
		assert.Nil(reply)
	})

	t.Run("output type is invalid", func(t *testing.T) {
		assert := assert.New(t)
		req := &pb.OptimRequest{
			Source:  pb.Type_JPEG,
			Output:  pb.Type_GZIP,
			Data:    jpegData,
			Quality: 10,
		}
		ctx := context.Background()
		reply, err := gs.DoOptim(ctx, req)
		assert.Equal(errOutputTypeIsInvalid, err)
		assert.Nil(reply)
	})

	t.Run("optim to jpeg", func(t *testing.T) {
		assert := assert.New(t)
		req := &pb.OptimRequest{
			Source:  pb.Type_JPEG,
			Output:  pb.Type_JPEG,
			Data:    jpegData,
			Quality: 10,
		}
		ctx := context.Background()
		reply, err := gs.DoOptim(ctx, req)
		assert.Nil(err)
		assert.Equal(pb.Type_JPEG, reply.Output)
		assert.NotNil(reply.Data)
	})

	t.Run("optim to png", func(t *testing.T) {
		assert := assert.New(t)
		req := &pb.OptimRequest{
			Source:  pb.Type_JPEG,
			Output:  pb.Type_PNG,
			Data:    jpegData,
			Quality: 10,
		}
		ctx := context.Background()
		reply, err := gs.DoOptim(ctx, req)
		assert.Nil(err)
		assert.Equal(pb.Type_PNG, reply.Output)
		assert.NotNil(reply.Data)
	})

	t.Run("optim to webp", func(t *testing.T) {
		assert := assert.New(t)
		req := &pb.OptimRequest{
			Source:  pb.Type_JPEG,
			Output:  pb.Type_WEBP,
			Data:    jpegData,
			Quality: 10,
		}
		ctx := context.Background()
		reply, err := gs.DoOptim(ctx, req)
		assert.Nil(err)
		assert.Equal(pb.Type_WEBP, reply.Output)
		assert.NotNil(reply.Data)
	})

	t.Run("optim to gzip", func(t *testing.T) {
		assert := assert.New(t)
		req := &pb.OptimRequest{
			Output:  pb.Type_GZIP,
			Data:    []byte("abcd"),
			Quality: 6,
		}
		ctx := context.Background()
		reply, err := gs.DoOptim(ctx, req)
		assert.Nil(err)
		assert.Equal(pb.Type_GZIP, reply.Output)
		assert.NotNil(reply.Data)
	})

	t.Run("optim to br", func(t *testing.T) {
		assert := assert.New(t)
		req := &pb.OptimRequest{
			Output:  pb.Type_BR,
			Data:    []byte("abcd"),
			Quality: 6,
		}
		ctx := context.Background()
		reply, err := gs.DoOptim(ctx, req)
		assert.Nil(err)
		assert.Equal(pb.Type_BR, reply.Output)
		assert.NotNil(reply.Data)
	})
}
