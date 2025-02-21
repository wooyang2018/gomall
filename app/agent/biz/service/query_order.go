package service

import (
	"context"
	agent "github.com/cloudwego/biz-demo/gomall/app/agent/kitex_gen/agent"
)

type QueryOrderService struct {
	ctx context.Context
} // NewQueryOrderService new QueryOrderService
func NewQueryOrderService(ctx context.Context) *QueryOrderService {
	return &QueryOrderService{ctx: ctx}
}

// Run create note info
func (s *QueryOrderService) Run(req *agent.QueryOrderReq) (resp *agent.QueryOrderResp, err error) {
	// Finish your business logic.

	return
}
