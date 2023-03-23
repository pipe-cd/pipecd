// Copyright 2023 The PipeCD Authors.
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

package model

import (
	"math/rand"
	"time"
	"unsafe"
)

var randomSrc = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // Number of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// GenerateRandomString makes a random string with the given length.
// This implementation was referenced from the link below.
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func GenerateRandomString(n int) string {
	b := make([]byte, n)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters.
	for i, cache, remain := n-1, randomSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
