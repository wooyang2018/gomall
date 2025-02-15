// Copyright 2024 CloudWeGo Authors
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

package utils

import "net/url"

var validHost = []string{
	"localhost:8080",
}

// ValidateNext 验证传入的 next 参数是否为合法的 URL，且该 URL 的主机名是否在预定义的合法主机名列表中。
func ValidateNext(next string) bool {
	urlObj, err := url.Parse(next)
	if err != nil {
		return false
	}
	if InArray(urlObj.Host, validHost) {
		return true
	}
	return false
}

// InArray 泛型函数，用于检查某个元素是否存在于一个切片中。
func InArray[T int | int32 | int64 | float32 | float64 | string](needle T, haystack []T) bool {
	for _, k := range haystack {
		if needle == k {
			return true
		}
	}
	return false
}
