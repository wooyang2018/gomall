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

	"github.com/cloudwego/kitex/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var TracerProvider *tracesdk.TracerProvider

func InitTracing(serviceName string) {
	// 创建一个新的 gRPC 导出器，用于将跟踪数据发送到 OpenTelemetry Collector
	exporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		panic(err)
	}

	// 注册一个关闭钩子函数，在服务器关闭时关闭导出器
	server.RegisterShutdownHook(func() {
		exporter.Shutdown(context.Background()) //nolint:errcheck
	})

	// 创建一个新的批处理跨度处理器，用于处理和导出跟踪数据
	processor := tracesdk.NewBatchSpanProcessor(exporter)

	// 创建一个新的资源实例，包含服务名称属性
	res, err := resource.New(context.Background(), resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)))
	if err != nil {
		// 如果创建资源失败，使用默认资源
		res = resource.Default()
	}

	// 创建一个新的跟踪提供程序实例，使用批处理跨度处理器和资源
	TracerProvider = tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(processor), tracesdk.WithResource(res))
	otel.SetTracerProvider(TracerProvider)
}
