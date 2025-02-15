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

import (
	"context"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	prometheus "github.com/kitex-contrib/monitor-prometheus"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/conf"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/mtl"
	"github.com/cloudwego/biz-demo/gomall/common/clientsuite"
	"github.com/cloudwego/biz-demo/gomall/common/utils"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/checkout/checkoutservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user/userservice"
)

var (
	ProductClient  productcatalogservice.Client
	UserClient     userservice.Client
	CartClient     cartservice.Client
	CheckoutClient checkoutservice.Client
	OrderClient    orderservice.Client

	once         sync.Once
	err          error
	registryAddr string
	commonSuite  client.Option
)

func InitClient() {
	once.Do(func() {
		registryAddr = conf.GetConf().Hertz.RegistryAddr
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: conf.ServiceName,
		})
		initProductClient()
		initUserClient()
		initCartClient()
		initCheckoutClient()
		initOrderClient()
	})
}

func initProductClient() {
	var opts []client.Option

	// 创建一个新的熔断套件，用于实现熔断机制。通过传入的函数将 rpcinfo.RPCInfo 转换为唯一的键，用于标识不同的 RPC 调用。
	cbs := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string {
		return circuitbreak.RPCInfo2Key(ri)
	})
	// 为特定的 RPC 方法 shop-frontend/product/GetProduct 更新熔断配置。这里设置了熔断机制为启用状态，错误率阈值为 0.5（即 50%），最小样本数为 2。
	cbs.UpdateServiceCBConfig("shop-frontend/product/GetProduct", circuitbreak.CBConfig{Enable: true, ErrRate: 0.5, MinSample: 2})
	// client.WithFallback：添加熔断降级策略。当 RPC 调用出错时，会根据具体的方法名进行不同的处理
	opts = append(opts, commonSuite, client.WithCircuitBreaker(cbs), client.WithFallback(fallback.NewFallbackPolicy(fallback.UnwrapHelper(func(ctx context.Context, req, resp interface{}, err error) (fbResp interface{}, fbErr error) {
		methodName := rpcinfo.GetRPCInfo(ctx).To().Method()
		if err == nil {
			return resp, err
		}
		// 如果调用的方法不是 ListProducts，同样返回原始响应和错误信息。
		if methodName != "ListProducts" {
			return resp, err
		}
		// 如果调用的方法是 ListProducts，则返回一个包含默认商品信息的 ListProductsResp 对象。
		return &product.ListProductsResp{
			Products: []*product.Product{
				{
					Price:       6.6,
					Id:          3,
					Picture:     "/static/image/t-shirt.jpeg",
					Name:        "T-Shirt",
					Description: "CloudWeGo T-Shirt",
				},
			},
		}, nil
	}))))

	// client.WithTracer：添加一个 Prometheus 客户端追踪器，用于收集和监控客户端的性能指标
	opts = append(opts, client.WithTracer(prometheus.NewClientTracer("", "", prometheus.WithDisableServer(true), prometheus.WithRegistry(mtl.Registry))))

	ProductClient, err = productcatalogservice.NewClient("product", opts...)
	utils.MustHandleError(err)
}

func initUserClient() {
	UserClient, err = userservice.NewClient("user", commonSuite)
	utils.MustHandleError(err)
}

func initCartClient() {
	CartClient, err = cartservice.NewClient("cart", commonSuite)
	utils.MustHandleError(err)
}

func initCheckoutClient() {
	CheckoutClient, err = checkoutservice.NewClient("checkout", commonSuite)
	utils.MustHandleError(err)
}

func initOrderClient() {
	OrderClient, err = orderservice.NewClient("order", commonSuite)
	utils.MustHandleError(err)
}
