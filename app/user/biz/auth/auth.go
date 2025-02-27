package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key") // JWT 签名密钥

// 生成 JWT
func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 校验 JWT
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

var TokenString string // 我们会优先解决框架层面的JWT Token传递问题，恢复使用标准的gRPC上下文传递机制

// Kitex 中间件：认证和鉴权
func AuthCasbinMiddleware() endpoint.Middleware {
	enforcer, err := casbin.NewEnforcer("./conf/model.conf", "./conf/policy.csv")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Casbin Enforcer: %v", err))
	}
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			// 从上下文中获取 JWT
			if TokenString == "" {
				return next(ctx, req, resp)
			}

			// 校验 JWT
			claims, err := ValidateToken(TokenString)
			if err != nil {
				return fmt.Errorf("unauthorized: %v", err)
			}

			// 从 JWT 中获取用户 ID 或角色
			userID := claims["user_id"].(string)
			role := "user" // 假设角色为 "user"，实际可以从数据库或 JWT 中获取

			// 获取请求的路径和方法
			path := ctx.Value("path").(string)
			method := ctx.Value("method").(string)

			// 使用 Casbin 进行权限校验
			if ok, err := enforcer.Enforce(role, path, method); !ok || err != nil {
				return fmt.Errorf("forbidden")
			}

			// 将用户信息传递给后续处理函数
			ctx = context.WithValue(ctx, "user_id", userID)
			ctx = context.WithValue(ctx, "role", role)

			// 调用下一个中间件或处理函数
			return next(ctx, req, resp)
		}
	}
}
