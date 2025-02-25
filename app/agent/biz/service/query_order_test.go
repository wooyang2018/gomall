package service

import (
	"context"
	"os"
	"testing"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/joho/godotenv"

	"github.com/cloudwego/biz-demo/gomall/app/agent/biz/rpc"
	"github.com/cloudwego/biz-demo/gomall/app/agent/conf"
	"github.com/cloudwego/biz-demo/gomall/common/mtl"
	"github.com/cloudwego/biz-demo/gomall/common/utils"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/agent"
)

func init() {
	os.Chdir("..") //nolint:errcheck
	os.Chdir("..") //nolint:errcheck
	if _, err := os.Getwd(); err != nil {
		klog.Error(err.Error())
	}

	// 加载环境变量文件
	if err := godotenv.Load(); err != nil {
		klog.Error(err.Error())
	}

	var serviceName = conf.GetConf().Kitex.Service
	mtl.InitTracing(serviceName)
	mtl.InitMetric(serviceName, conf.GetConf().Kitex.MetricsPort, conf.GetConf().Registry.RegistryAddress[0])

	rpc.InitClient()
}

func TestQueryOrder_Run(t *testing.T) {
	userId := int64(1)
	ctx := context.WithValue(context.Background(), utils.UserIdKey, userId)
	s := NewQueryOrderService(ctx)

	// init req and assert value
	req := &agent.QueryOrderReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)
}
