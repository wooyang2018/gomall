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

	"github.com/cloudwego/biz-demo/gomall/app/payment/biz/dal"
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
	dal.Init() // 初始化MySQL实例
}

func mockChargeReq() *payment.ChargeReq {
	// 构造 CreditCardInfo 实例
	creditCard := &payment.CreditCardInfo{
		CreditCardNumber:          "4111111111111111",
		CreditCardCvv:             123,
		CreditCardExpirationYear:  2025,
		CreditCardExpirationMonth: 12,
	}

	// 构造 ChargeReq 实例
	chargeReq := &payment.ChargeReq{
		Amount:     100.50,
		CreditCard: creditCard,
		OrderId:    "order1",
		UserId:     1,
	}

	return chargeReq
}

// GO_ENV=dev go test -run TestCharge_Run
func TestCharge_Run(t *testing.T) {
	ctx := context.Background()
	s := NewChargeService(ctx)

	// init req and assert value
	req := mockChargeReq()
	resp, err := s.Run(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	}
}
