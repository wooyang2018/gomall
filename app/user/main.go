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

	"github.com/cloudwego/biz-demo/gomall/app/user/biz/auth"
	"github.com/cloudwego/biz-demo/gomall/app/user/biz/dal"
	"github.com/cloudwego/biz-demo/gomall/app/user/conf"
	"github.com/cloudwego/biz-demo/gomall/common/mtl"
	"github.com/cloudwego/biz-demo/gomall/common/serversuite"
	"github.com/cloudwego/biz-demo/gomall/common/utils"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user/userservice"
)

var serviceName = conf.GetConf().Kitex.Service

func main() {
	// 加载环境变量文件
	if err := godotenv.Load(); err != nil {
		klog.Error(err.Error())
	}

	// 初始化日志记录器
	mtl.InitLog(&lumberjack.Logger{
		// 日志文件名
		Filename: conf.GetConf().Kitex.LogFileName,
		// 每个日志文件的最大大小（MB）
		MaxSize: conf.GetConf().Kitex.LogMaxSize,
		// 保留的旧日志文件的最大数量
		MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
		// 保留旧日志文件的最大天数
		MaxAge: conf.GetConf().Kitex.LogMaxAge,
	})
	klog.SetLevel(conf.LogLevel())

	// 初始化跟踪器
	mtl.InitTracing(serviceName)

	// 初始化指标收集器
	mtl.InitMetric(serviceName, conf.GetConf().Kitex.MetricsPort, conf.GetConf().Registry.RegistryAddress[0])

	// 初始化数据访问层
	dal.Init()

	// 创建一个新的Kitex服务实例m，启动Kitex服务
	opts := kitexInit()
	svr := userservice.NewServer(new(UserServiceImpl), opts...)
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
	opts = append(opts, server.WithServiceAddr(addr), server.WithSuite(serversuite.CommonServerSuite{CurrentServiceName: serviceName, RegistryAddr: conf.GetConf().Registry.RegistryAddress[0]}))
	opts = append(opts, server.WithMiddleware(auth.AuthCasbinMiddleware()))
	return
}
