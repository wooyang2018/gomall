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

package serversuite

import (
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/config-consul/consul"
	consulServer "github.com/kitex-contrib/config-consul/server"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	registryconsul "github.com/kitex-contrib/registry-consul"

	"github.com/cloudwego/biz-demo/gomall/common/mtl"
)

type CommonServerSuite struct {
	CurrentServiceName string
	RegistryAddr       string
}

func (s CommonServerSuite) Options() []server.Option {
	opts := []server.Option{
		// 设置服务器的元数据处理程序为 transmeta.ServerHTTP2Handler，用于处理 HTTP/2 元数据。
		server.WithMetaHandler(transmeta.ServerHTTP2Handler),
	}

	r, err := registryconsul.NewConsulRegister(s.RegistryAddr)
	if err != nil {
		klog.Fatal(err)
	}

	opts = append(opts, server.WithRegistry(r))

	if os.Getenv("CONFIG_CENTER_ENABLED") == "true" {
		// 尝试从环境变量 CONFIG_CENTER_NODES 中获取 Consul 节点地址。
		consulNodes := os.Getenv("CONFIG_CENTER_NODES")
		if consulNodes != "" {
			consulClient, err := consul.NewClient(consul.Options{})
			if err != nil {
				klog.Error(err)
			} else {
				// 服务注册与发现（server.WithRegistry(r)）：
				// - 确保服务可以被其他服务发现和调用。
				// - 这是微服务架构中服务间通信的基础。
				// 动态配置管理（server.WithSuite(consulServer.NewSuite(...))）：
				// - 确保服务可以从 Consul 中获取动态配置（如超时时间、重试策略等）。
				// - 这是实现配置中心化管理和动态更新的关键。
				// 两者的功能是互补的：
				// - 服务注册与发现解决的是“服务在哪里”的问题。
				// - 动态配置管理解决的是“服务如何运行”的问题。
				// 将consulClient与当前服务名一起传递给 consulServer.NewSuite 函数，生成一个配置中心套件
				opts = append(opts, server.WithSuite(consulServer.NewSuite(s.CurrentServiceName, consulClient)))
			}
		}
	}

	// 创建了一个 OpenTelemetry 提供程序实例，并将其配置为使用 mtl.TracerProvider 作为 SDK 跟踪提供程序，同时禁用指标收集。
	_ = provider.NewOpenTelemetryProvider(
		provider.WithSdkTracerProvider(mtl.TracerProvider),
		provider.WithEnableMetrics(false),
	)

	opts = append(opts,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		// tracing.NewServerSuite() 创建了一个 OpenTelemetry 的服务器套件，用于在服务中启用分布式追踪。
		// 这个套件会自动为服务的 RPC 调用添加追踪信息（如 TraceID、SpanID 等），并将这些信息导出到配置的
		// OpenTelemetry 后端（如 Jaeger、Zipkin 等）。
		server.WithSuite(tracing.NewServerSuite()),
		// prometheus.WithDisableServer(true) 表示禁用 Prometheus 服务器，prometheus.WithRegistry(mtl.Registry) 表示使用 mtl.Registry 作为 Prometheus 注册中心。
		// NewServerTracer provides tracer for server access, addr and path is the scrape_configs for prometheus server.
		server.WithTracer(prometheus.NewServerTracer("", "", prometheus.WithDisableServer(true), prometheus.WithRegistry(mtl.Registry))),
	)

	return opts
}
