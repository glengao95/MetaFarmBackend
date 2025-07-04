package blockchain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthClient 以太坊客户端实现
type EthClient struct {
	client     *ethclient.Client
	rpcURL     string
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

// NewEthClient 创建新的以太坊客户端
func NewEthClient(rpcURL, privateKeyHex string) (*EthClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("无法连接到以太坊节点: %w", err)
	}

	var privateKey *ecdsa.PrivateKey
	var address common.Address

	if privateKeyHex != "" {
		privateKey, err = crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			return nil, fmt.Errorf("无效的私钥: %w", err)
		}
		address = crypto.PubkeyToAddress(privateKey.PublicKey)
	}

	return &EthClient{
		client:     client,
		rpcURL:     rpcURL,
		privateKey: privateKey,
		address:    address,
	}, nil
}

// Connect 连接到以太坊节点
func (e *EthClient) Connect(ctx context.Context) error {
	if e.client == nil {
		client, err := ethclient.Dial(e.rpcURL)
		if err != nil {
			return fmt.Errorf("无法连接到以太坊节点: %w", err)
		}
		e.client = client
	}
	return nil
}

// Close 关闭客户端连接
func (e *EthClient) Close() {
	if e.client != nil {
		e.client.Close()
	}
}

// ChainID 获取链ID
func (e *EthClient) ChainID(ctx context.Context) (*big.Int, error) {
	return e.client.ChainID(ctx)
}

// BlockNumber 获取最新区块号
func (e *EthClient) BlockNumber(ctx context.Context) (*big.Int, error) {
	header, err := e.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

// GetBlockByNumber 根据区块号获取区块
func (e *EthClient) GetBlockByNumber(ctx context.Context, number *big.Int) (interface{}, error) {
	ethBlock, err := e.client.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	return ethBlock, nil
}

// BalanceAt 获取账户余额
func (e *EthClient) BalanceAt(ctx context.Context, address string) (*big.Int, error) {
	addr := common.HexToAddress(address)
	return e.client.BalanceAt(ctx, addr, nil)
}

// NonceAt 获取账户Nonce
func (e *EthClient) NonceAt(ctx context.Context, address string) (uint64, error) {
	addr := common.HexToAddress(address)
	return e.client.NonceAt(ctx, addr, nil)
}

// SendTransaction 发送交易
func (e *EthClient) SendTransaction(ctx context.Context, opts TxOptions) (string, error) {
	if e.privateKey == nil {
		return "", errors.New("客户端未初始化私钥")
	}

	from := common.HexToAddress(opts.From)
	if from == (common.Address{}) {
		from = e.address
	}

	to := common.HexToAddress(opts.To)

	// 如果未指定GasPrice，则获取当前建议的GasPrice
	gasPrice := opts.GasPrice
	if gasPrice == nil {
		var err error
		gasPrice, err = e.client.SuggestGasPrice(ctx)
		if err != nil {
			return "", fmt.Errorf("获取GasPrice失败: %w", err)
		}
	}

	// 如果未指定GasLimit，则估算Gas
	gasLimit := opts.GasLimit
	if gasLimit == 0 {
		var err error
		msg := ethereum.CallMsg{
			From:  from,
			To:    &to,
			Value: opts.Value,
			Data:  opts.Data,
		}
		gasLimit, err = e.client.EstimateGas(ctx, msg)
		if err != nil {
			return "", fmt.Errorf("估算Gas失败: %w", err)
		}
	}

	// 获取Nonce
	nonce, err := e.client.NonceAt(ctx, from, nil)
	if err != nil {
		return "", fmt.Errorf("获取Nonce失败: %w", err)
	}

	// 创建交易
	chainID, err := e.client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("获取ChainID失败: %w", err)
	}

	tx := types.NewTransaction(nonce, to, opts.Value, gasLimit, gasPrice, opts.Data)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), e.privateKey)
	if err != nil {
		return "", fmt.Errorf("签名交易失败: %w", err)
	}

	// 发送交易
	if err := e.client.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("发送交易失败: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

// CallContract 调用合约
func (e *EthClient) CallContract(ctx context.Context, contractAddr string, data []byte) ([]byte, error) {
	addr := common.HexToAddress(contractAddr)
	msg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}
	return e.client.CallContract(ctx, msg, nil)
}

// SignMessage 签名消息
func (e *EthClient) SignMessage(privateKeyHex string, message []byte) ([]byte, error) {
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
func (e *EthClient) VerifySignature(address string, message []byte, signature []byte) (bool, error) {
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

// GetClient 获取底层ethclient.Client
func (e *EthClient) GetClient() *ethclient.Client {
	return e.client
}
