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

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/cloudwego/biz-demo/gomall/common/mq"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
)

type SendService struct {
	ctx context.Context
}

// NewSendService new SendService
func NewSendService(ctx context.Context) *SendService {
	return &SendService{ctx: ctx}
}

func (s *SendService) Run(req *email.EmailReq) (resp *email.EmailResp, err error) {
	data, _ := proto.Marshal(req)
	msg := &nats.Msg{Subject: "email", Data: data, Header: make(nats.Header)}
	if err := mq.Nc.PublishMsg(msg); err != nil {
		return nil, err
	}
	return &email.EmailResp{}, nil
}
