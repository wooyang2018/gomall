// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package consumer

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"

	"github.com/cloudwego/biz-demo/gomall/common/mq"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
)

func Init() {
	// 创建一个名为 "shop-nats-consumer" 的 OpenTelemetry 跟踪器，用于跟踪消息消费的过程
	tracer := otel.Tracer("shop-nats-consumer")

	// 订阅名为 "email" 的 NATS 主题，并定义一个回调函数来处理接收到的消息
	sub, err := mq.Nc.Subscribe("email", func(m *nats.Msg) {
		var req email.EmailReq
		err := proto.Unmarshal(m.Data, &req)
		if err != nil {
			klog.Error(err)
		}

		noopEmail := NewNoopEmail()
		_ = noopEmail.Send(&req)

		ctx := context.Background()
		// Extract reads cross-cutting concerns from the carrier into a Context.
		// m.Header 是一个 map[string][]string 类型的对象，用于存储 HTTP 请求头的键值对
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(m.Header))
		// 使用 tracer 创建一个新的 span，用于跟踪消息处理过程，并将其与当前的上下文 ctx 关联起来
		// span 是 OpenTelemetry 中的一个概念，用于跟踪和记录分布式系统中的操作。span 代表了一个
		// 操作的开始和结束，以及该操作的一些元数据，如操作名称、开始时间、结束时间、标签等。
		_, span := tracer.Start(ctx, "shop-email-consumer")
		// 在函数结束时结束 span
		defer span.End()
	})

	if err != nil {
		panic(err)
	}

	server.RegisterShutdownHook(func() {
		sub.Unsubscribe() //nolint:errcheck
		mq.Nc.Close()     // 关闭 NATS 连接
	})
}
