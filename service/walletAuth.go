package service

import (
	"MetaFarmBackend/dao"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// WalletAuthService 钱包认证服务接口
type WalletAuthService interface {
	// 生成登录消息和随机数
	GenerateLoginMessage(ctx context.Context, walletAddress string) (string, string, error)

	// 验证签名并登录
	VerifySignatureAndLogin(ctx context.Context, walletAddress, signature, nonce string,
		ipAddress, userAgent string) (*LoginResult, error)

	// 验证会话令牌
	VerifySessionToken(ctx context.Context, token string) (*SessionInfo, error)

	// 注销会话
	RevokeSession(ctx context.Context, token string) error
}

// 实现WalletAuthService接口
type walletAuthServiceImpl struct {
	dao        *dao.Dao
	sessionTTL time.Duration // 会话有效期
}

// 登录结果
type LoginResult struct {
	UserID        uint64    `json:"user_id"`
	WalletAddress string    `json:"wallet_address"`
	SessionToken  string    `json:"session_token"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// 会话信息
type SessionInfo struct {
	UserID        uint64    `json:"user_id"`
	WalletAddress string    `json:"wallet_address"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// 构造函数
func NewWalletAuthService(dao *dao.Dao, sessionTTL time.Duration) WalletAuthService {
	return &walletAuthServiceImpl{
		dao:        dao,
		sessionTTL: sessionTTL,
	}
}

// 生成登录消息和随机数
func (s *walletAuthServiceImpl) GenerateLoginMessage(ctx context.Context, walletAddress string) (string, string, error) {
	// 标准化钱包地址为小写
	walletAddress = strings.ToLower(walletAddress)

	// 生成随机数
	nonce := generateRandomNonce()

	// 构建登录消息
	message := buildLoginMessage(walletAddress, nonce)

	// 保存或更新用户钱包记录
	err := s.saveOrUpdateWallet(ctx, walletAddress, nonce)
	if err != nil {
		return "", "", errors.Wrap(err, "保存钱包信息失败")
	}

	return message, nonce, nil
}

// 验证签名并登录
func (s *walletAuthServiceImpl) VerifySignatureAndLogin(ctx context.Context, walletAddress, signature, nonce string,
	ipAddress, userAgent string) (*LoginResult, error) {

	// 标准化钱包地址为小写
	walletAddress = strings.ToLower(walletAddress)

	// 验证签名
	err := verifySignature(walletAddress, signature, nonce)
	if err != nil {
		// 记录登录失败日志
		s.dao.RecordLoginLog(ctx, walletAddress, ipAddress, userAgent, false, err.Error())
		return nil, errors.Wrap(err, "签名验证失败")
	}

	// 查找用户钱包记录
	wallet, err := s.dao.GetUserWalletByAddressAndNonce(walletAddress, nonce)
	if err != nil {
		return nil, err
	}

	// 创建新会话
	sessionToken, expiresAt, err := s.createSession(ctx, wallet.UserID, walletAddress)
	if err != nil {
		return nil, errors.Wrap(err, "创建会话失败")
	}

	// 更新最后登录时间
	err = s.dao.UpdateLastLoginAt(walletAddress)
	if err != nil {
		return nil, errors.Wrap(err, "更新最后登录时间失败")
	}

	// 记录登录成功日志
	s.dao.RecordLoginLog(ctx, walletAddress, ipAddress, userAgent, true, "")

	return &LoginResult{
		UserID:        wallet.UserID,
		WalletAddress: walletAddress,
		SessionToken:  sessionToken,
		ExpiresAt:     expiresAt,
	}, nil
}

// 验证会话令牌
func (s *walletAuthServiceImpl) VerifySessionToken(ctx context.Context, token string) (*SessionInfo, error) {
	session, err := s.dao.GetValidSessionByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("无效的会话令牌")
		}
		return nil, errors.Wrap(err, "查询会话失败")
	}

	return &SessionInfo{
		UserID:        session.UserID,
		WalletAddress: session.WalletAddress,
		ExpiresAt:     session.ExpiresAt,
	}, nil
}

// 注销会话
func (s *walletAuthServiceImpl) RevokeSession(ctx context.Context, token string) error {
	return s.dao.RevokeSessionByToken(token)
}

// 保存或更新用户钱包记录
func (s *walletAuthServiceImpl) saveOrUpdateWallet(ctx context.Context, walletAddress, nonce string) error {
	_, err := s.dao.GetUserWalletByAddress(walletAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新用户和钱包记录
			user := dao.User{
				Username:  generateUsername(walletAddress),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := s.dao.CreateUser(&user); err != nil {
				return errors.Wrap(err, "创建用户失败")
			}

			wallet := dao.UserWallet{
				UserID:        user.ID,
				WalletAddress: walletAddress,
				Nonce:         nonce,
				LastLoginAt:   time.Now(),
				IsPrimary:     true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			return s.dao.CreateUserWallet(&wallet)
		}
		return errors.Wrap(err, "查询用户钱包失败")
	}

	// 更新随机数
	return s.dao.UpdateNonce(walletAddress, nonce)
}

// 创建会话
func (s *walletAuthServiceImpl) createSession(ctx context.Context, userID uint64, walletAddress string) (string, time.Time, error) {
	// 生成会话ID和令牌
	sessionID := uuid.New().String()
	sessionToken := generateSessionToken()
	expiresAt := time.Now().Add(s.sessionTTL)

	// 保存会话
	session := dao.LoginSession{
		ID:            sessionID,
		UserID:        userID,
		WalletAddress: walletAddress,
		Token:         sessionToken,
		ExpiresAt:     expiresAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return sessionToken, expiresAt, s.dao.CreateLoginSession(&session)
}

// 生成随机用户名
func generateUsername(walletAddress string) string {
	// 取钱包地址后8位
	suffix := walletAddress[len(walletAddress)-8:]
	return fmt.Sprintf("user_%s", suffix)
}

// 生成随机数
func generateRandomNonce() string {
	// 生成16字节随机数
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// 生成会话令牌
func generateSessionToken() string {
	return uuid.New().String()
}

// 构建登录消息
func buildLoginMessage(walletAddress, nonce string) string {
	// EIP-712风格的登录消息
	domain := apitypes.TypedDataDomain{
		Name:    "MetaFarm",
		Version: "1.0.0",
		ChainId: (*math.HexOrDecimal256)(big.NewInt(137)), // Polygon主网
		Salt:    "0x0000000000000000000000000000000000000000000000000000000000000000",
	}

	message := map[string]interface{}{
		"wallet":  walletAddress,
		"nonce":   nonce,
		"expires": time.Now().Add(30 * time.Minute).Unix(),
	}

	typedData := apitypes.TypedData{
		Domain: domain,
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
			},
			"Login": []apitypes.Type{
				{Name: "wallet", Type: "address"},
				{Name: "nonce", Type: "string"},
				{Name: "expires", Type: "uint256"},
			},
		},
		PrimaryType: "Login",
		Message:     message,
	}

	// 序列化消息
	rawData, err := typedData.HashStruct("Login", typedData.Message)
	if err != nil {
		// 使用简单消息格式作为后备
		return fmt.Sprintf("Welcome to MetaFarm!\n\nClick to sign in and accept the MetaFarm Terms of Service:\n\nhttps://metafarm.com/tos\n\nThis request will not trigger a blockchain transaction or cost any gas fees.\n\nWallet address:\n%s\n\nNonce:\n%s", walletAddress, nonce)
	}

	return hex.EncodeToString(rawData)
}

// 验证签名
func verifySignature(walletAddress, signature, nonce string) error {
	// 标准化钱包地址
	walletAddress = strings.ToLower(walletAddress)

	// 移除签名前缀（如果有）
	if strings.HasPrefix(signature, "0x") {
		signature = signature[2:]
	}

	// 解析签名
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return errors.Wrap(err, "解析签名失败")
	}

	// 检查签名长度
	if len(sigBytes) != 65 {
		return errors.New("无效的签名长度")
	}

	// 调整v值（某些钱包返回的v值为27/28，需要转换为0/1）
	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}

	// 构建要签名的消息
	message := buildLoginMessage(walletAddress, nonce)
	signHash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	// 恢复公钥
	recoveredPubKey, err := crypto.SigToPub(signHash.Bytes(), sigBytes)
	if err != nil {
		return errors.Wrap(err, "恢复公钥失败")
	}

	// 从公钥计算地址
	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey)

	// 验证地址匹配
	if strings.ToLower(recoveredAddr.Hex()) != walletAddress {
		return errors.New("签名与钱包地址不匹配")
	}

	return nil
}

// 布尔转整数
