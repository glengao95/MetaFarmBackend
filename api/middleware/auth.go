package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWT 配置
const (
	jwtSecret = "your_jwt_secret_key"
	jwtIssuer = "MetaFarmBackend"
)

// JWTClaims JWT 声明结构
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			Fail(c, http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		// 2. 解析并验证JWT
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			Fail(c, http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			// 3. 将用户信息存入上下文
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Next()
		} else {
			Fail(c, http.StatusUnauthorized, nil)
			c.Abort()
			return
		}
	}
}

// GenerateJWTToken 生成JWT token
func GenerateJWTToken(userID int, username string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    jwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// SessionMiddleware Session管理中间件
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从cookie获取session ID
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			// 创建新session
			sessionID = generateSessionID()
			c.SetCookie("session_id", sessionID, 3600, "/", "", false, true)
		}

		// 2. 从存储中获取session数据（这里简化为直接存入context）
		c.Set("sessionID", sessionID)
		c.Next()
	}
}

func generateSessionID() string {
	// 实际项目中应使用更安全的随机字符串生成方式
	return "session_" + time.Now().Format("20060102150405")
}
