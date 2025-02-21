package agent

import (
	"context"
	agent "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/agent"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func QueryOrder(ctx context.Context, req *agent.QueryOrderReq, callOptions ...callopt.Option) (resp *agent.QueryOrderResp, err error) {
	resp, err = defaultClient.QueryOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "QueryOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
