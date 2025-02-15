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

package mtl

import (
	"context"

	"github.com/cloudwego/hertz/pkg/route"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/conf"
)

var TracerProvider *tracesdk.TracerProvider

func InitTracing() route.CtxCallback {
	exporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		panic(err)
	}
	// 这里的 exporter 是一个追踪数据的导出器，NewBatchSpanProcessor 会把收集到的跨度批量发送给这个导出器
	processor := tracesdk.NewBatchSpanProcessor(exporter)
	res, err := resource.New(context.Background(), resource.WithAttributes(semconv.ServiceNameKey.String(conf.ServiceName)))
	if err != nil {
		res = resource.Default()
	}

	//管理 Span 处理：TracerProvider 可以配置一个或多个 SpanProcessor，用于处理创建的 Span。例如，将 Span 批量发送到导出器。
	//资源关联：可以将 Span 与特定的资源（如服务名称、版本等）关联起来，这些资源信息会随 Span 一起被导出，方便在追踪系统中进行识别和分析。
	TracerProvider = tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(processor), tracesdk.WithResource(res))
	otel.SetTracerProvider(TracerProvider)

	return func(ctx context.Context) {
		exporter.Shutdown(ctx) //nolint:errcheck
	}
}

// Trace 代表了一个完整的业务操作或者请求在分布式系统中的执行路径。它就像一条线，把系统中各个服务参与处理该请求的过程串联起来，形成一个完整的调用链
// Span 是 Trace 中的一个基本工作单元，代表了在某个服务或者组件中执行的一个特定操作。每个 Span 都有自己的开始时间和结束时间，用于记录该操作的执行时长。
// 组成：一个 Span 通常包含以下信息：
// 	操作名称：描述该 Span 所代表的操作，如数据库查询、HTTP 请求等。
// 	开始时间和结束时间：用于计算操作的执行时长。
// 	父 Span 信息：在分布式系统中，一个 Span 可能是由另一个 Span 调用产生的，因此需要记录其父 Span 的信息，以构建调用关系。
// 	标签（Tags）：用于记录一些与操作相关的元数据，如数据库表名、HTTP 状态码等。
// 	日志（Logs）：用于记录操作执行过程中的关键信息，如错误信息、调试信息等。
// 上下文传播：在分布式系统中，为了将 Trace 和 Span 的信息在不同服务之间传递，需要进行上下文传播。常见的方式是通过 HTTP 头、消息队列等方式将 Trace ID、Span ID 等信息传递给下游服务。
