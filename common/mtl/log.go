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

package mtl

import (
	"io"
	"os"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLog 初始化日志记录器
func InitLog(ioWriter io.Writer) {
	// 定义一个 kitexzap.Option 类型的切片，用于存储日志记录器的配置选项
	var opts []kitexzap.Option
	// 定义一个 zapcore.WriteSyncer 类型的变量，用于存储日志记录器的输出目标
	var output zapcore.WriteSyncer

	// 判断当前环境是否为线上环境
	if os.Getenv("GO_ENV") != "online" {
		// 如果不是线上环境，使用控制台编码器，并将日志输出到 ioWriter
		opts = append(opts, kitexzap.WithCoreEnc(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())))
		output = zapcore.AddSync(ioWriter)
	} else {
		// 如果是线上环境，使用 JSON 编码器，并将日志异步输出到 ioWriter
		opts = append(opts, kitexzap.WithCoreEnc(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())))
		// zapcore.BufferedWriteSyncer 是一个缓冲写入器，它会将日志消息先写入一个缓冲区，
		// 然后按照指定的时间间隔（FlushInterval）将缓冲区中的消息批量写入到实际的输出目标。
		output = &zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(ioWriter),
			FlushInterval: time.Minute,
		}
	}

	// 注册一个关闭钩子，在程序关闭时同步日志输出
	server.RegisterShutdownHook(func() {
		// nolint:errcheck 是一个特殊的注释，用于告诉静态代码分析工具
		// 如 golint、go vet 等，忽略对特定代码行的错误检查。
		output.Sync() //nolint:errcheck
	})

	// 创建一个新的日志记录器，并将其设置为全局日志记录器
	log := kitexzap.NewLogger(opts...)
	klog.SetLogger(log)
	klog.SetOutput(output)
}
