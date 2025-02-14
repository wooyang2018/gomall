// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/rpc"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/conf"
	"github.com/cloudwego/biz-demo/gomall/common/mq"
	"github.com/cloudwego/biz-demo/gomall/common/mtl"
	"github.com/cloudwego/biz-demo/gomall/common/serversuite"
	"github.com/cloudwego/biz-demo/gomall/common/utils"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/checkout/checkoutservice"
)

var serviceName = conf.GetConf().Kitex.Service

func main() {
	_ = godotenv.Load()

	// gopkg.in/natefinch/lumberjack.v2 是一个用于日志滚动的 Go 语言库。
	// 当日志文件达到指定的大小（例如 100MB）时，它会自动将当前日志文件重命名为一个
	// 新的文件名（例如 app.log.1），并创建一个新的 app.log 文件来继续记录日志。
	mtl.InitLog(&lumberjack.Logger{
		Filename:   conf.GetConf().Kitex.LogFileName,
		MaxSize:    conf.GetConf().Kitex.LogMaxSize,
		MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
		MaxAge:     conf.GetConf().Kitex.LogMaxAge,
	})
	klog.SetLevel(conf.LogLevel())

	mtl.InitTracing(serviceName)
	mtl.InitMetric(serviceName, conf.GetConf().Kitex.MetricsPort, conf.GetConf().Registry.RegistryAddress[0])

	rpc.InitClient()
	mq.Init()

	opts := kitexInit()
	svr := checkoutservice.NewServer(new(CheckoutServiceImpl), opts...)
	if err := svr.Run(); err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	address := conf.GetConf().Kitex.Address
	if strings.HasPrefix(address, ":") {
		localIp := utils.MustGetLocalIPv4()
		address = localIp + address
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))
	opts = append(opts, server.WithSuite(serversuite.CommonServerSuite{CurrentServiceName: serviceName, RegistryAddr: conf.GetConf().Registry.RegistryAddress[0]}))
	return
}
