package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
)

type ProtectedService struct {
	ctx context.Context
}

// NewProtectedService new ProtectedService
func NewProtectedService(ctx context.Context) *ProtectedService {
	return &ProtectedService{ctx: ctx}
}

// Run create note info
func (s *ProtectedService) Run(req *user.ProtectedReq) (resp *user.ProtectedResp, err error) {
	// 从上下文中获取用户信息，并记录日志
	userID := s.ctx.Value("user_id").(string)
	role := s.ctx.Value("role").(string)
	klog.Infof("Protected method called by user %s with role %s", userID, role)

	return &user.ProtectedResp{
		Result: "Hello World",
	}, nil
}
