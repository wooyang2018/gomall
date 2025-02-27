package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/kitex/client"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/conf"
	"github.com/cloudwego/biz-demo/gomall/common/clientsuite"
	"github.com/cloudwego/biz-demo/gomall/common/utils"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user/userservice"
)

// GO_ENV=dev go test -run TestProtected_Run
func TestProtected_Run(t *testing.T) {
	registryAddr := conf.GetConf().Hertz.RegistryAddr
	commonSuite := client.WithSuite(clientsuite.CommonGrpcClientSuite{
		RegistryAddr:       registryAddr,
		CurrentServiceName: conf.ServiceName,
	})
	UserClient, err := userservice.NewClient("user", commonSuite)
	utils.MustHandleError(err)

	// init req and assert value
	req := &user.LoginReq{
		Email:    "123456@qq.com",
		Password: "123456",
	}
	_, err = UserClient.Login(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	resp, err := UserClient.Protected(context.Background(), &user.ProtectedReq{Query: "Hello"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	fmt.Printf("resp: %+v\n", resp)
}
