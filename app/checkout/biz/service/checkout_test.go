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

package service

import (
	"context"
	"os"
	"testing"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/joho/godotenv"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/rpc"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/conf"
	"github.com/cloudwego/biz-demo/gomall/common/mq"
	"github.com/cloudwego/biz-demo/gomall/common/mtl"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/checkout"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment"
)

func init() {
	os.Chdir("..") //nolint:errcheck
	os.Chdir("..") //nolint:errcheck
	if _, err := os.Getwd(); err != nil {
		klog.Error(err.Error())
	}

	// 加载环境变量文件
	if err := godotenv.Load(); err != nil {
		klog.Error(err.Error())
	}

	var serviceName = conf.GetConf().Kitex.Service
	mtl.InitTracing(serviceName)
	mtl.InitMetric(serviceName, conf.GetConf().Kitex.MetricsPort, conf.GetConf().Registry.RegistryAddress[0])

	rpc.InitClient()
	mq.Init()
}

func mockCheckoutReq() *checkout.CheckoutReq {
	// 构造Address对象
	address := &checkout.Address{
		StreetAddress: "Main St",
		City:          "San Francisco",
		State:         "CA",
		Country:       "USA",
		ZipCode:       "94107",
	}

	// 构造CreditCardInfo对象
	creditCard := &payment.CreditCardInfo{
		CreditCardNumber:          "4111111111111111",
		CreditCardCvv:             123,
		CreditCardExpirationYear:  2025,
		CreditCardExpirationMonth: 12,
	}

	// 构造CheckoutReq对象
	checkoutReq := &checkout.CheckoutReq{
		UserId:     1,
		Firstname:  "John",
		Lastname:   "Doe",
		Email:      "john.doe@example.com",
		Address:    address,
		CreditCard: creditCard,
	}

	return checkoutReq
}

// GO_ENV=dev go test -run TestCheckout_Run
func TestCheckout_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCheckoutService(ctx)

	// init req and assert value
	resp, err := s.Run(mockCheckoutReq())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	}
}
