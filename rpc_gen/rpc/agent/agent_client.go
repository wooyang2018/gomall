package agent

import (
	"context"
	agent "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/agent"

	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/agent/agentservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() agentservice.Client
	Service() string
	QueryOrder(ctx context.Context, Req *agent.QueryOrderReq, callOptions ...callopt.Option) (r *agent.QueryOrderResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := agentservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient agentservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() agentservice.Client {
	return c.kitexClient
}

func (c *clientImpl) QueryOrder(ctx context.Context, Req *agent.QueryOrderReq, callOptions ...callopt.Option) (r *agent.QueryOrderResp, err error) {
	return c.kitexClient.QueryOrder(ctx, Req, callOptions...)
}
