package service

import (
	"context"
	"testing"

	user "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
)

func TestProtected_Run(t *testing.T) {
	ctx := context.Background()
	s := NewProtectedService(ctx)
	// init req and assert value

	req := &user.ProtectedReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
