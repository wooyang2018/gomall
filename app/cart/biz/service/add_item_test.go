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

	"github.com/cloudwego/biz-demo/gomall/app/cart/biz/dal"
	"github.com/cloudwego/biz-demo/gomall/app/cart/biz/rpc"
	"github.com/cloudwego/biz-demo/gomall/app/cart/conf"
	"github.com/cloudwego/biz-demo/gomall/common/mtl"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
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

	dal.Init() // 初始化MySQL实例
	rpc.InitClient()
}

// 创建一个 AddItemReq 实例
func mockAddItemReq() *cart.AddItemReq {
	req := &cart.AddItemReq{
		UserId: 1,
		Item: &cart.CartItem{
			ProductId: 1,
			Quantity:  3,
		},
	}
	return req
}

// GO_ENV=dev go test -run TestAddItem_Run
func TestAddItem_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAddItemService(ctx)

	// init req and assert value
	req := mockAddItemReq()
	resp, err := s.Run(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	}
}
