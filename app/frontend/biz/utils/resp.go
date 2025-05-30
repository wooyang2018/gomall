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

	"github.com/cloudwego/biz-demo/gomall/common/utils"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/hertz/pkg/app"
)

// SendErrResponse  pack error response
func SendErrResponse(ctx context.Context, c *app.RequestContext, code int, err error) {
	// todo edit custom code
	c.String(code, err.Error())
}

func WarpResponse(ctx context.Context, c *app.RequestContext, content map[string]any) map[string]any {
	var cartNum int
	userId := utils.GetUserIdFromCtx(ctx)
	cartResp, _ := CartClient.GetCart(ctx, &cart.GetCartReq{UserId: userId})
	if cartResp != nil && cartResp.Cart != nil {
		cartNum = len(cartResp.Cart.Items)
	}
	content[utils.UserIdKey] = ctx.Value(utils.UserIdKey)
	content["cart_num"] = cartNum
	return content
}
