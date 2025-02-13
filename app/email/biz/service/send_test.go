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
	"os"
	"testing"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/joho/godotenv"

	"github.com/cloudwego/biz-demo/gomall/app/email/biz/consumer"
	"github.com/cloudwego/biz-demo/gomall/common/mq"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
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
	mq.Init()
	consumer.Init()
}

func mockEmailReq() *email.EmailReq {
	req := &email.EmailReq{
		From:        "from@example.com",
		To:          "to@example.com",
		ContentType: "text/plain",
		Subject:     "You just created an order in CloudWeGo shop",
		Content:     fmt.Sprintf("You just created an order in CloudWeGo shop at %s", time.Now().Format(time.RFC822)),
	}
	return req
}

// GO_ENV=dev go test -run TestSend_Run
func TestSend_Run(t *testing.T) {
	ctx := context.Background()
	s := NewSendService(ctx)

	// init req and assert value
	resp, err := s.Run(mockEmailReq())
	time.Sleep(500 * time.Millisecond)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	}
}
