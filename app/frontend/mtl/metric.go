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
	"fmt"
	"net"
	"net/http"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/conf"
	"github.com/cloudwego/biz-demo/gomall/common/utils"

	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hertz-contrib/registry/consul"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Registry *prometheus.Registry

func initMetric() route.CtxCallback {
	// 创建一个新的 Prometheus 注册表，用于管理所有的指标收集器。
	Registry = prometheus.NewRegistry()
	// 注册 Go 运行时的指标收集器，用于收集 Go 程序的运行时信息，如 goroutine 数量、内存使用等。
	Registry.MustRegister(collectors.NewGoCollector())
	// 注册进程级别的指标收集器，用于收集进程的信息，如 CPU 使用、内存使用等。
	Registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	config := consulapi.DefaultConfig()
	config.Address = conf.GetConf().Hertz.RegistryAddr
	consulClient, _ := consulapi.NewClient(config)
	// 创建一个 Consul 服务注册器，并添加标签 service:frontend。
	r := consul.NewConsulRegister(consulClient, consul.WithAdditionInfo(&consul.AdditionInfo{
		Tags: []string{"service:frontend"},
	}))

	localIp := utils.MustGetLocalIPv4()
	ip, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", localIp, conf.GetConf().Hertz.MetricsPort))
	if err != nil {
		hlog.Error(err)
	}
	// 创建一个服务注册信息，包含服务地址、服务名称和权重，并将服务注册到 Consul 注册中心。
	registryInfo := &registry.Info{Addr: ip, ServiceName: "prometheus", Weight: 1}
	err = r.Register(registryInfo)
	if err != nil {
		hlog.Error(err)
	}

	http.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	// 启动一个 HTTP 服务器，监听指定的指标端口，用于提供 Prometheus 指标。
	go http.ListenAndServe(fmt.Sprintf(":%d", conf.GetConf().Hertz.MetricsPort), nil) //nolint:errcheck

	return func(ctx context.Context) {
		r.Deregister(registryInfo) //nolint:errcheck
	}
}
