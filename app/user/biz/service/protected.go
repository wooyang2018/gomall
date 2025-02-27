package service

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
)

type ProtectedService struct {
	ctx context.Context
} // NewProtectedService new ProtectedService
func NewProtectedService(ctx context.Context) *ProtectedService {
	return &ProtectedService{ctx: ctx}
}

// Run create note info
func (s *ProtectedService) Run(req *user.ProtectedReq) (resp *user.ProtectedResp, err error) {
	// Finish your business logic.

	return
}
