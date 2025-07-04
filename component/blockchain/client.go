package blockchain

import (
	"context"
	"math/big"
)

// BlockchainClient 定义区块链客户端通用接口
type BlockchainClient interface {
	// 连接相关方法
	Connect(ctx context.Context) error
	Close()
	ChainID(ctx context.Context) (*big.Int, error)

	// 区块相关方法
	BlockNumber(ctx context.Context) (*big.Int, error)
	GetBlockByNumber(ctx context.Context, number *big.Int) (interface{}, error)

	// 账户相关方法
	BalanceAt(ctx context.Context, address string) (*big.Int, error)
	NonceAt(ctx context.Context, address string) (uint64, error)

	// 交易相关方法
	SendTransaction(ctx context.Context, opts TxOptions) (string, error)
	CallContract(ctx context.Context, contractAddr string, data []byte) ([]byte, error)

	// 签名验签方法
	SignMessage(privateKeyHex string, message []byte) ([]byte, error)
	VerifySignature(address string, message []byte, signature []byte) (bool, error)
}

// 交易选项
type TxOptions struct {
	From     string          // 发送者地址
	To       string          // 接收者地址
	Value    *big.Int        // 转账金额
	GasLimit uint64          // Gas限制
	GasPrice *big.Int        // Gas价格
	Data     []byte          // 交易数据
}


