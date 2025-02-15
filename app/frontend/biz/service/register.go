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

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/biz/utils"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/auth"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/common"
	comutils "github.com/cloudwego/biz-demo/gomall/common/utils"
	rpcuser "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
)

type RegisterService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRegisterService(Context context.Context, RequestContext *app.RequestContext) *RegisterService {
	return &RegisterService{RequestContext: RequestContext, Context: Context}
}

func (h *RegisterService) Run(req *auth.RegisterReq) (resp *common.Empty, err error) {
	res, err := utils.UserClient.Register(h.Context, &rpcuser.RegisterReq{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.Password,
	})
	if err != nil {
		return nil, err
	}

	session := sessions.Default(h.RequestContext)
	session.Set(comutils.UserIdKey, res.UserId)
	// Save 方法会将会话数据持久化到存储后端（如内存、文件、Redis等）。
	err = session.Save()
	if err != nil {
		return nil, err
	}
	return
}
