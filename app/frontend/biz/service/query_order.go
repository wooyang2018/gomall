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
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"

	bizutils "github.com/cloudwego/biz-demo/gomall/app/frontend/biz/utils"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/product"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/types"
	comutils "github.com/cloudwego/biz-demo/gomall/common/utils"
	rpcagent "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/agent"
	rpcproduct "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
)

type QueryOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewQueryOrderService(Context context.Context, RequestContext *app.RequestContext) *QueryOrderService {
	return &QueryOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *QueryOrderService) Run(req *product.SearchProductsReq) (resp map[string]any, err error) {
	userId := comutils.GetUserIdFromCtx(h.Context)
	var orders []*types.Order
	queryOrderResp, err := bizutils.AgentClient.QueryOrder(h.Context, &rpcagent.QueryOrderReq{UserId: userId, Question: req.Q})
	if err != nil {
		return nil, err
	}
	// If there is no order, return directly
	if queryOrderResp == nil || len(queryOrderResp.Orders) == 0 {
		return utils.H{
			"title":    "Order",
			"orders":   orders,
			"response": "未查询到相关订单",
		}, nil
	}

	for _, v := range queryOrderResp.Orders {
		var items []types.OrderItem
		var total float32
		for _, vv := range v.OrderItems {
			total += vv.Cost
			i := vv.Item
			productResp, err := bizutils.ProductClient.GetProduct(h.Context, &rpcproduct.GetProductReq{Id: i.ProductId})
			if err != nil {
				return nil, err
			}
			if productResp.Product == nil {
				continue
			}
			p := productResp.Product
			items = append(items, types.OrderItem{
				ProductId:   i.ProductId,        // 商品ID
				Qty:         uint32(i.Quantity), // 商品数量
				ProductName: p.Name,             // 商品名称
				Picture:     p.Picture,          // 商品图片
				Cost:        vv.Cost,            // 商品成本
			})
		}
		// 将订单创建时间戳转换为时间对象，第一个参数是秒，第二个参数是纳秒
		timeObj := time.Unix(int64(v.CreatedAt), 0)
		orders = append(orders, &types.Order{
			Cost:  total,
			Items: items,
			// 订单创建日期，格式化为 "2006-01-02 15:04:05"
			CreatedDate: timeObj.Format("2006-01-02 15:04:05"),
			OrderId:     v.OrderId,
			Consignee:   types.Consignee{Email: v.Email},
		})
	}

	return utils.H{
		"title":    "Order",
		"orders":   orders,
		"response": queryOrderResp.Response,
	}, nil
}
