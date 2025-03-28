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

package redis

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/extra/redisprometheus/v9"
	"github.com/redis/go-redis/v9"

	"github.com/cloudwego/biz-demo/gomall/app/product/conf"
	"github.com/cloudwego/biz-demo/gomall/common/mtl"
)

var RedisClient *redis.Client

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().Redis.Address,
		Username: conf.GetConf().Redis.Username,
		Password: conf.GetConf().Redis.Password,
		DB:       conf.GetConf().Redis.DB,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	// 这行代码使用 redisotel.InstrumentTracing 函数为 Redis 客户端添加分布式追踪功能。
	// redisotel 是一个用于将 Redis 操作与 OpenTelemetry 集成的库，它可以帮助开发者追踪 Redis 请求的执行情况
	if err := redisotel.InstrumentTracing(RedisClient); err != nil {
		klog.Error("redis tracing collect error ", err)
	}
	if err := mtl.Registry.Register(redisprometheus.NewCollector("default", "product", RedisClient)); err != nil {
		klog.Error("redis metric collect error ", err)
	}
}
