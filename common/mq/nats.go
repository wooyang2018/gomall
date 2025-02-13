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

package mq

import (
	"github.com/nats-io/nats.go"
)

var (
	Nc  *nats.Conn
	err error
)

// NATS是一个轻量级、高性能的开源消息传递系统，用于构建分布式和微服务架构中的实时通信。
// 它提供了简单的发布-订阅和请求-回复模式，使得不同的服务之间可以高效地交换消息。
func Init() {
	Nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
}
