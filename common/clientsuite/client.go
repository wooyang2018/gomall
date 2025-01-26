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

package clientsuite

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

type CommonGrpcClientSuite struct {
	CurrentServiceName string
	RegistryAddr       string
}

// Options 返回一个包含客户端选项的切片，用于配置Kitex客户端
func (s CommonGrpcClientSuite) Options() []client.Option {
	// 创建一个新的Consul解析器，用于服务发现
	r, err := consul.NewConsulResolver(s.RegistryAddr)
	if err != nil {
		panic(err)
	}

	opts := []client.Option{
		// 使用Consul解析器进行服务发现
		client.WithResolver(r),
		// 使用HTTP/2元数据处理器
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		// 使用gRPC作为传输协议
		client.WithTransportProtocol(transport.GRPC),
		// 设置客户端基本信息，包括服务名称
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		// 使用链路追踪套件
		client.WithSuite(tracing.NewClientSuite()),
	}
	return opts
}
