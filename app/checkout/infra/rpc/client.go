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

package rpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/conf"
	"github.com/cloudwego/biz-demo/gomall/common/clientsuite"
	"github.com/cloudwego/biz-demo/gomall/common/utils"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment/paymentservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
)

var (
	CartClient    cartservice.Client
	ProductClient productcatalogservice.Client
	PaymentClient paymentservice.Client
	OrderClient   orderservice.Client
	once          sync.Once
	err           error
	registryAddr  string
	serviceName   string
	commonSuite   client.Option
)

func InitClient() {
	once.Do(func() {
		registryAddr = conf.GetConf().Registry.RegistryAddress[0]
		serviceName = conf.GetConf().Kitex.Service
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			CurrentServiceName: serviceName,
			RegistryAddr:       registryAddr,
		})
		initCartClient()
		initProductClient()
		initPaymentClient()
		initOrderClient()
	})
}

func initProductClient() {
	ProductClient, err = productcatalogservice.NewClient("product", commonSuite)
	utils.MustHandleError(err)
}

func initCartClient() {
	CartClient, err = cartservice.NewClient("cart", commonSuite)
	utils.MustHandleError(err)
}

func initPaymentClient() {
	PaymentClient, err = paymentservice.NewClient("payment", commonSuite)
	utils.MustHandleError(err)
}

func initOrderClient() {
	OrderClient, err = orderservice.NewClient("order", commonSuite)
	utils.MustHandleError(err)
}
