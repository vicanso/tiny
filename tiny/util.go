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
// limitations under the License

package tiny

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
)

// Commander commander
type Commander func(string, string) []string

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const digitBytes = "0123456789"

// randomString create a random string
func randomString(baseLetters string, n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(baseLetters) {
			b[i] = baseLetters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func doCommandConvert(data []byte, fn Commander, writer *bytes.Buffer) (err error) {
	filename := randomString(letterBytes, 10)
	tmpfile, err := ioutil.TempFile("", filename)
	if err != nil {
		return
	}
	originalFile := tmpfile.Name()
	defer tmpfile.Close()
	defer os.Remove(originalFile)

	_, err = tmpfile.Write(data)
	if err != nil {
		return
	}
	targetFile := originalFile + "-new"
	args := fn(originalFile, targetFile)
	cmd := exec.Command(args[0], args[1:]...)
	err = cmd.Run()
	if err != nil {
		return
	}
	// 删除临时文件
	defer os.Remove(targetFile)
	target, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return
	}
	writer.Write(target)
	return
}
