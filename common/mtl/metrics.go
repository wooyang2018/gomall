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
	"net"
	"net/http"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Registry *prometheus.Registry

// InitMetric 初始化 Prometheus 指标收集器并注册到 Consul 注册中心
// 参数:
//
//	serviceName: 服务名称
//	metricsPort: 指标收集器监听的端口
//	registryAddr: Consul 注册中心的地址
func InitMetric(serviceName string, metricsPort string, registryAddr string) {
	// 创建一个新的 Prometheus 注册中心
	Registry = prometheus.NewRegistry()
	// 注册 Go 运行时指标收集器
	Registry.MustRegister(collectors.NewGoCollector())
	// 注册进程指标收集器
	Registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// 创建一个新的 Consul 注册中心实例
	r, _ := consul.NewConsulRegister(registryAddr)
	// 解析 metricsPort 为 TCP 地址
	addr, _ := net.ResolveTCPAddr("tcp", metricsPort)

	// 创建一个新的注册中心信息实例
	registryInfo := &registry.Info{
		ServiceName: "prometheus",
		Addr:        addr,
		Weight:      1,                                         // 设置权重为 1
		Tags:        map[string]string{"service": serviceName}, // 设置标签，包含服务名称
	}

	// 将注册中心信息注册到 Consul 注册中心
	_ = r.Register(registryInfo)
	// 注册一个关闭钩子函数，在服务器关闭时注销注册中心信息
	server.RegisterShutdownHook(func() {
		r.Deregister(registryInfo) //nolint:errcheck
	})

	// 注册一个 HTTP 处理程序，用于暴露 Prometheus 指标
	http.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	// 在一个新的 goroutine 中启动 HTTP 服务器，监听 metricsPort 端口，忽略错误
	go http.ListenAndServe(metricsPort, nil) //nolint:errcheck
}
