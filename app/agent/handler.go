package main

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/agent/biz/service"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/agent"
)

// AgentServiceImpl implements the last service interface defined in the IDL.
type AgentServiceImpl struct{}

// QueryOrder implements the AgentServiceImpl interface.
func (s *AgentServiceImpl) QueryOrder(ctx context.Context, req *agent.QueryOrderReq) (resp *agent.QueryOrderResp, err error) {
	resp, err = service.NewQueryOrderService(ctx).Run(req)

	return resp, err
}
