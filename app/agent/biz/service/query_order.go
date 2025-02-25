package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/biz-demo/gomall/app/agent/biz/rpc"
	comutils "github.com/cloudwego/biz-demo/gomall/common/utils"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/agent"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
)

type QueryOrderService struct {
	ctx context.Context
}

// NewQueryOrderService new QueryOrderService
func NewQueryOrderService(ctx context.Context) *QueryOrderService {
	return &QueryOrderService{ctx: ctx}
}

func (s *QueryOrderService) Run(req *agent.QueryOrderReq) (resp *agent.QueryOrderResp, err error) {
	userId := comutils.GetUserIdFromCtx(s.ctx)
	listOrderResp, err := rpc.OrderClient.ListOrder(s.ctx, &order.ListOrderReq{UserId: userId})
	if err != nil {
		return nil, err
	}
	// If there is no order, return directly
	if listOrderResp == nil || len(listOrderResp.Orders) == 0 {
		return nil, fmt.Errorf("no matched order")
	}
	return &agent.QueryOrderResp{
		Orders:   listOrderResp.Orders,
		Response: "测试响应",
	}, nil
}
