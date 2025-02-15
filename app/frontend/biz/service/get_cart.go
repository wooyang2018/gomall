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
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"

	bizutils "github.com/cloudwego/biz-demo/gomall/app/frontend/biz/utils"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/common"
	comutils "github.com/cloudwego/biz-demo/gomall/common/utils"
	rpccart "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	rpcproduct "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
)

type GetCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetCartService(Context context.Context, RequestContext *app.RequestContext) *GetCartService {
	return &GetCartService{RequestContext: RequestContext, Context: Context}
}

func (h *GetCartService) Run(req *common.Empty) (resp map[string]any, err error) {
	var items []map[string]string
	carts, err := bizutils.CartClient.GetCart(h.Context, &rpccart.GetCartReq{
		UserId: comutils.GetUserIdFromCtx(h.Context),
	})
	if err != nil {
		return nil, err
	}
	var total float32
	for _, v := range carts.Cart.Items {
		productResp, err := bizutils.ProductClient.GetProduct(h.Context, &rpcproduct.GetProductReq{Id: v.GetProductId()})
		if err != nil {
			continue
		}
		if productResp.Product == nil {
			continue
		}
		p := productResp.Product
		items = append(items, map[string]string{"Name": p.Name, "Description": p.Description, "Picture": p.Picture, "Price": strconv.FormatFloat(float64(p.Price), 'f', 2, 64), "Qty": strconv.Itoa(int(v.Quantity))})
		total += float32(v.Quantity) * p.Price
	}

	return utils.H{
		"title": "Cart",
		"items": items,
		// FormatFloat 根据格式 fmt 和精度 prec 将浮点数 f 转换为字符串。它假定原始值是从 bitSize 位（float32 为 32 位，float64 为 64 位）的浮点数得到的，并对结果进行四舍五入。
		// 格式 fmt 可以是以下之一：
		// 'b'：表示为 -ddddp±ddd 形式，即二进制指数形式。
		// 'e'：表示为 -d.dddd e±dd 形式，即十进制指数形式，指数部分用小写 e。
		// 'E'：表示为 -d.dddd E±dd 形式，即十进制指数形式，指数部分用大写 E。
		// 'f'：表示为 -ddd.dddd 形式，即无指数形式。
		// 'g'：对于大指数使用 'e' 格式，其他情况使用 'f' 格式。
		// 'G'：对于大指数使用 'E' 格式，其他情况使用 'f' 格式。
		// 精度 prec 控制 'e'、'E'、'f'、'g'、'G'、'x' 和 'X' 格式打印的数字位数（不包括指数部分）。对于 'e'、'E'、'f'、'x' 和 'X' 格式，它是小数点后的数字位数。
		"total": strconv.FormatFloat(float64(total), 'f', 2, 64),
	}, nil
}
