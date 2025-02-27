package auth

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
	"github.com/cloudwego/kitex/server"
)

// 示例 Kitex 服务
type ExampleService struct{}

func (s *ExampleService) ProtectedMethod(ctx context.Context, req interface{}) (resp interface{}, err error) {
	userID := ctx.Value("user_id").(string)
	role := ctx.Value("role").(string)
	klog.Infof("User %s with role %s accessed the protected method", userID, role)
	return "Access granted", nil
}

func TestExampleService(t *testing.T) {
	// 初始化 Casbin Enforcer
	enforcer, err := NewEnforcer()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Casbin Enforcer: %v", err))
	}

	addr, err := net.ResolveTCPAddr("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	// 创建 Kitex 服务
	svr := server.NewServer(
		server.WithServiceAddr(addr),
		server.WithMiddleware(AuthCasbinMiddleware(enforcer)),
	)
	// 注册服务
	exampleService := &ExampleService{}
	err = svr.RegisterService(newServiceInfo(), exampleService)
	if err != nil {
		panic(fmt.Sprintf("Failed to register service: %v", err))
	}
	// 启动服务
	if err := svr.Run(); err != nil {
		klog.Fatal(err)
	}
}

// 创建服务元信息
func newServiceInfo() *serviceinfo.ServiceInfo {
	return &serviceinfo.ServiceInfo{
		ServiceName: "example",
	}
}

// 定义测试方法
func TestProtectedMethod(t *testing.T) {
	// 创建服务实例
	exampleService := &ExampleService{}

	// 模拟请求上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Authorization", "<your-token>")
	ctx = context.WithValue(ctx, "path", "/protected")
	ctx = context.WithValue(ctx, "method", "GET")

	// 调用服务方法
	resp, err := exampleService.ProtectedMethod(ctx, nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	} else {
		fmt.Println("Response:", resp)
	}
}
