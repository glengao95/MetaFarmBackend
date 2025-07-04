package router

import (
	"MetaFarmBackend/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// WalletAuthController 钱包认证控制器
type WalletAuthController struct {
	walletAuthService service.WalletAuthService
}

// 构造函数
func NewWalletAuthController(walletAuthService service.WalletAuthService) *WalletAuthController {
	return &WalletAuthController{
		walletAuthService: walletAuthService,
	}
}

// GenerateLoginMessage 生成登录消息和随机数
func (c *WalletAuthController) GenerateLoginMessage(ctx *gin.Context) {
	var request struct {
		WalletAddress string `json:"wallet_address" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成登录消息和随机数
	message, nonce, err := c.walletAuthService.GenerateLoginMessage(ctx, request.WalletAddress)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
		"nonce":   nonce,
	})
}

// VerifySignatureAndLogin 验证签名并登录
func (c *WalletAuthController) VerifySignatureAndLogin(ctx *gin.Context) {
	var request struct {
		WalletAddress string `json:"wallet_address" binding:"required"`
		Signature     string `json:"signature" binding:"required"`
		Nonce         string `json:"nonce" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取客户端信息
	ipAddress := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	// 验证签名并登录
	result, err := c.walletAuthService.VerifySignatureAndLogin(ctx,
		request.WalletAddress, request.Signature, request.Nonce, ipAddress, userAgent)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 设置会话Cookie（可选）
	c.setSessionCookie(ctx, result.SessionToken)

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":        result.UserID,
		"wallet_address": result.WalletAddress,
		"session_token":  result.SessionToken,
		"expires_at":     result.ExpiresAt,
	})
}

// Logout 注销
func (c *WalletAuthController) Logout(ctx *gin.Context) {
	// 从请求头或Cookie获取会话令牌
	token := c.getSessionToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "缺少会话令牌"})
		return
	}

	// 注销会话
	if err := c.walletAuthService.RevokeSession(ctx, token); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 清除会话Cookie（如果有）
	ctx.SetCookie("session_token", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "注销成功"})
}

// AuthMiddleware 认证中间件
func (c *WalletAuthController) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头或Cookie获取会话令牌
		token := c.getSessionToken(ctx)
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		// 验证会话令牌
		session, err := c.walletAuthService.VerifySessionToken(ctx, token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		// 将用户ID和钱包地址存入上下文
		ctx.Set("user_id", session.UserID)
		ctx.Set("wallet_address", session.WalletAddress)

		ctx.Next()
	}
}

// 获取会话令牌
func (c *WalletAuthController) getSessionToken(ctx *gin.Context) string {
	// 优先从请求头获取
	if token := ctx.GetHeader("Authorization"); token != "" {
		if strings.HasPrefix(token, "Bearer ") {
			return token[7:]
		}
		return token
	}

	// 从Cookie获取
	if token, err := ctx.Cookie("session_token"); err == nil {
		return token
	}

	return ""
}

// 设置会话Cookie
func (c *WalletAuthController) setSessionCookie(ctx *gin.Context, token string) {
	ctx.SetCookie(
		"session_token",
		token,
		86400, // 1天有效期
		"/",   // 路径
		"",    // 域名（默认当前域名）
		false, // 是否仅HTTPS
		true,  // 是否HTTPOnly
	)
}

// 注册路由
func (c *WalletAuthController) RegisterRoutes(r *gin.Engine) {
	r.POST("/login/message", c.GenerateLoginMessage)
	r.POST("/login", c.VerifySignatureAndLogin)
	r.POST("/logout", c.Logout)
}
