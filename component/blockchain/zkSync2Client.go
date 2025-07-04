package blockchain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
)

// ZkSync2Client zkSync2客户端实现
type ZkSync2Client struct {
	client     *clients.Client
	l1Client   *ethclient.Client
	wallet     *accounts.Wallet
	rpcURL     string
	privateKey *ecdsa.PrivateKey
	address    common.Address
	chainID    *big.Int
	mu         sync.Mutex
}

// NewZkSync2Client 创建新的zkSync2客户端
func NewZkSync2Client(rpcURL, privateKeyHex string, l1Client *ethclient.Client) (*ZkSync2Client, error) {
	client, err := clients.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("无法连接到zkSync节点: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取链ID失败: %w", err)
	}

	var privateKey *ecdsa.PrivateKey
	var address common.Address
	var wallet *accounts.Wallet

	if privateKeyHex != "" {
		privateKey, err = crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			return nil, fmt.Errorf("无效的私钥: %w", err)
		}
		address = crypto.PubkeyToAddress(privateKey.PublicKey)

		// 创建zkSync钱包
		wallet, err = accounts.NewWallet(crypto.FromECDSA(privateKey), client, l1Client)
		if err != nil {
			return nil, fmt.Errorf("创建钱包失败: %w", err)
		}
	}

	return &ZkSync2Client{
		client:     client,
		l1Client:   l1Client,
		wallet:     wallet,
		rpcURL:     rpcURL,
		privateKey: privateKey,
		address:    address,
		chainID:    chainID,
	}, nil
}

// Connect 连接到zkSync节点
func (z *ZkSync2Client) Connect(ctx context.Context) error {
	if z.client == nil {
		client, err := clients.Dial(z.rpcURL)
		if err != nil {
			return fmt.Errorf("无法连接到zkSync节点: %w", err)
		}
		z.client = client

		chainID, err := client.ChainID(ctx)
		if err != nil {
			return fmt.Errorf("获取链ID失败: %w", err)
		}
		z.chainID = chainID
	}
	return nil
}

// Close 关闭客户端连接
func (z *ZkSync2Client) Close() {
	if z.client != nil {
		z.client.Close()
	}
}

// ChainID 获取链ID
func (z *ZkSync2Client) ChainID(ctx context.Context) (*big.Int, error) {
	return z.chainID, nil
}

// BlockNumber 获取最新区块号
func (z *ZkSync2Client) BlockNumber(ctx context.Context) (*big.Int, error) {
	header, err := z.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

// GetBlockByNumber 根据区块号获取区块
func (z *ZkSync2Client) GetBlockByNumber(ctx context.Context, number *big.Int) (interface{}, error) {
	zkBlock, err := z.client.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	return zkBlock, nil
}

// BalanceAt 获取账户余额
func (z *ZkSync2Client) BalanceAt(ctx context.Context, address string) (*big.Int, error) {
	addr := common.HexToAddress(address)
	return z.client.BalanceAt(ctx, addr, nil)
}

// NonceAt 获取账户Nonce
func (z *ZkSync2Client) NonceAt(ctx context.Context, address string) (uint64, error) {
	addr := common.HexToAddress(address)
	return z.client.NonceAt(ctx, addr, nil)
}

// SendTransaction 发送交易
func (z *ZkSync2Client) SendTransaction(ctx context.Context, opts TxOptions) (string, error) {
	z.mu.Lock()
	defer z.mu.Unlock()

	if z.wallet == nil {
		return "", errors.New("客户端未初始化钱包")
	}

	return "", nil
}

// CallContract 调用合约
func (z *ZkSync2Client) CallContract(ctx context.Context, contractAddr string, data []byte) ([]byte, error) {
	addr := common.HexToAddress(contractAddr)
	msg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}
	return z.client.CallContract(ctx, msg, nil)
}

// SignMessage 签名消息
func (z *ZkSync2Client) SignMessage(privateKeyHex string, message []byte) ([]byte, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("无效的私钥: %w", err)
	}

	// 以太坊签名前缀
	prefixedMessage := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message))
	hash := crypto.Keccak256Hash(prefixedMessage)
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名失败: %w", err)
	}

	// 调整v值 (EIP-155)
	signature[64] += 27

	return signature, nil
}

// VerifySignature 验证签名
func (z *ZkSync2Client) VerifySignature(address string, message []byte, signature []byte) (bool, error) {
	if len(signature) != 65 {
		return false, errors.New("签名长度必须为65字节")
	}

	// 恢复v值
	signature[64] -= 27

	// 以太坊签名前缀
	prefixedMessage := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message))
	hash := crypto.Keccak256Hash(prefixedMessage)

	// 从签名中恢复公钥
	pubKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return false, fmt.Errorf("恢复公钥失败: %w", err)
	}

	// 验证地址
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return recoveredAddr.Hex() == address, nil
}

// GetL1Client 获取L1客户端
func (z *ZkSync2Client) GetL1Client() *ethclient.Client {
	return z.l1Client
}

// GetL2Client 获取L2客户端
func (z *ZkSync2Client) GetL2Client() *clients.Client {
	return z.client
}

// GetWallet 获取钱包实例
func (z *ZkSync2Client) GetWallet() *accounts.Wallet {
	return z.wallet
}
