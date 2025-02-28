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

	bizutils "github.com/cloudwego/biz-demo/gomall/app/frontend/biz/utils"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/auth"
	"github.com/cloudwego/biz-demo/gomall/common/utils"
	rpcuser "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
)

type LoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginService(Context context.Context, RequestContext *app.RequestContext) *LoginService {
	return &LoginService{RequestContext: RequestContext, Context: Context}
}

func (h *LoginService) Run(req *auth.LoginReq) (resp string, err error) {
	res, err := bizutils.UserClient.Login(h.Context, &rpcuser.LoginReq{Email: req.Email, Password: req.Password})
	if err != nil {
		return
	}

	// 特别注意 RequestContext *app.RequestContext 才是真正的请求上下文，从中可以获取 session
	// 会话数据默认存储在客户端的Cookie中，session.Set 和 session.Save 用于设置会话数据并保存会话。
	session := sessions.Default(h.RequestContext)
	session.Set(utils.UserIdKey, res.UserId)
	session.Set(utils.UserToken, res.Token)
	err = session.Save()
	utils.MustHandleError(err)

	redirect := "/" // 默认重定向到首页
	// 如果有 next 参数，重定向到 next 参数指定的页面
	if utils.ValidateNext(req.Next) {
		redirect = req.Next
	}
	if err != nil {
		return "", err
	}
	return redirect, nil
}
