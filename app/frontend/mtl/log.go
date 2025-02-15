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
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/server"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/conf"
)

func initLog() {
	var opts []hertzzap.Option
	var output zapcore.WriteSyncer
	if os.Getenv("GO_ENV") != "online" {
		// 配置日志编码器，使其输出适合开发环境的易读格式。
		opts = append(opts, hertzzap.WithCoreEnc(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())))
		output = os.Stdout
	} else {
		// 配置日志编码器，使其输出 JSON 格式的日志，便于机器解析。
		opts = append(opts, hertzzap.WithCoreEnc(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())))
		// 实现异步日志写入，每隔一分钟刷新一次缓冲区。
		output = &zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(os.Stdout),
			FlushInterval: time.Minute,
		}
		server.RegisterShutdownHook(func() {
			output.Sync() //nolint:errcheck
		})
	}
	log := hertzzap.NewLogger(opts...)
	hlog.SetLogger(log)
	hlog.SetLevel(conf.LogLevel())
	hlog.SetOutput(output)
}
