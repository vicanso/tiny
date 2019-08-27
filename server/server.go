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

import "github.com/vicanso/hes"

var (
	errOutputIsNil               = hes.New("output can not be nil")
	errOutputTypeIsInvalid       = hes.New("output type is not supported")
	errURLIsNil                  = hes.New("url can not be nil")
	errContentTypeIsNil          = hes.New("can not get content type of resource")
	errContentTypeIsInvalid      = hes.New("content type is invalid")
	errContentTypeIsNotSupported = hes.New("content type is not supported")
	errImageIsNil                = hes.New("image data can not be nil")
	errTextIsNil                 = hes.New("text data can not be nil")
	errDataIsNil                 = hes.New("data can not be nil")
)
