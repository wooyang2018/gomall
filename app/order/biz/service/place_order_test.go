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
	"fmt"
	"testing"

	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
)

func mockPlaceOrderReq() *order.PlaceOrderReq {
	// 构造 Address 实例
	address := &order.Address{
		StreetAddress: "Main St",
		City:          "Anytown",
		State:         "CA",
		Country:       "USA",
		ZipCode:       12345,
	}

	// 构造 OrderItem 实例
	orderItem := &order.OrderItem{
		Item: &cart.CartItem{
			ProductId: 1,
			Quantity:  2,
		},
		Cost: 19.99,
	}

	// 构造 PlaceOrderReq 实例
	placeOrderReq := &order.PlaceOrderReq{
		UserId:       1,
		UserCurrency: "USD",
		Address:      address,
		Email:        "123456@qq.com",
		OrderItems:   []*order.OrderItem{orderItem},
	}

	return placeOrderReq
}

// GO_ENV=dev go test -run TestPlaceOrder_Run
func TestPlaceOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewPlaceOrderService(ctx)

	// init req and assert value
	req := mockPlaceOrderReq()
	resp, err := s.Run(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	}
	fmt.Println("创建订单成功，", resp.Order)
}
