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

package utils

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
)

// MustHandleError log the error info and then exit
func MustHandleError(err error) {
	if err != nil {
		klog.Fatal(err)
	}
}

const UserIdKey = "user_id"
const UserToken = "user_token"
const UserRole = "role"
const JWTExpire = "expire"
const JWTSecret = "your-secret-key" // JWT 签名密钥

func GetUserIdFromCtx(ctx context.Context) uint32 {
	if ctx.Value(UserIdKey) == nil {
		return 0
	}
	return uint32(ctx.Value(UserIdKey).(float64))
}
