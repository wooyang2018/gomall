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

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hertz-contrib/sessions"

	"github.com/cloudwego/biz-demo/gomall/common/utils"
)

// GlobalAuth 是一个全局认证中间件，用于检查用户是否已登录。
// 如果用户已登录，它将用户ID存储在上下文中，以便后续处理程序可以访问。
// 如果用户未登录，它将继续执行登录处理程序。
func GlobalAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取当前请求的会话
		session := sessions.Default(c)
		// 从会话中获取用户ID
		userId := session.Get(utils.UserIdKey)
		// 如果用户ID为空，说明用户未登录
		if userId == nil {
			// 继续执行下一个处理程序
			c.Next(ctx)
			return
		}
		// 如果用户已登录，将用户ID存储在上下文中
		ctx = context.WithValue(ctx, utils.UserIdKey, userId)
		// 继续执行下一个处理程序
		c.Next(ctx)
	}
}

func Auth() app.HandlerFunc {
	enforcer, err := casbin.NewEnforcer("./conf/model.conf", "./conf/policy.csv")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Casbin Enforcer: %v", err))
	}

	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		userId := session.Get(utils.UserIdKey)
		token := session.Get(utils.UserToken)
		if userId == nil || token == nil {
			// 获取请求头中的 Referer 字段，该字段表示用户从哪个页面跳转过来
			ref := string(c.GetHeader("Referer"))
			next := "/sign-in"
			if ref != "" && utils.ValidateNext(ref) {
				// 构建一个带有 next 查询参数的登录页面 URL，该参数记录了用户跳转过来的页面。
				next = fmt.Sprintf("%s?next=%s", next, ref)
			}
			// 将用户重定向到登录页面。
			c.Redirect(302, []byte(next))
			c.Abort() // 终止当前请求的处理链
			return
		}

		// 校验 JWT
		claims, err := ValidateToken(token.(string))
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized, Validate Token Failed")
			c.Abort()
			return
		}

		// 从 JWT 中获取用户 ID 或角色
		userID := claims[utils.UserIdKey].(float64)
		role := claims[utils.UserRole].(string)
		if int(userID) != int(userId.(float64)) {
			c.String(http.StatusUnauthorized, "Unauthorized, User ID Mismatch")
			c.Abort()
			return
		}

		// 获取请求的 Method 和 Path，使用 Casbin 进行权限校验
		method := string(c.Request.Method())
		path := string(c.Request.URI().Path())
		if ok, err := enforcer.Enforce(role, path, method); !ok || err != nil {
			c.String(http.StatusForbidden, "Permission Denied, Request Path: %s, Method: %s", path, method)
			c.Abort()
			return
		}

		// 将用户信息传递给后续处理函数
		ctx = context.WithValue(ctx, utils.UserIdKey, userId)
		c.Next(ctx)
	}
}

// 校验 JWT
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
