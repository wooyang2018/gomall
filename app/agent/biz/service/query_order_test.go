package service

import (
	"context"
	agent "github.com/cloudwego/biz-demo/gomall/app/agent/kitex_gen/agent"
	"testing"
)

func TestQueryOrder_Run(t *testing.T) {
	ctx := context.Background()
	s := NewQueryOrderService(ctx)
	// init req and assert value

	req := &agent.QueryOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
