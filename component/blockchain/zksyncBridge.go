package blockchain

import (
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
)

// ZkSyncBridge zkSync L1-L2桥接服务
type ZkSyncBridge struct {
	l1Client   *EthClient
	l2Client   *ZkSync2Client
	bridgeAddr common.Address // L1桥接合约地址
}

// NewZkSyncBridge 创建新的桥接服务
func NewZkSyncBridge(l1Client *EthClient, l2Client *ZkSync2Client, bridgeAddr string) (*ZkSyncBridge, error) {
	if l1Client == nil || l2Client == nil {
		return nil, errors.Wrap(nil, "L1和L2客户端不能为空")
	}

	addr := common.HexToAddress(bridgeAddr)
	if addr == (common.Address{}) {
		return nil, errors.Errorf("无效的桥接合约地址")
	}

	return &ZkSyncBridge{
		l1Client:   l1Client,
		l2Client:   l2Client,
		bridgeAddr: addr,
	}, nil
}

// // DepositETH 从L1存款ETH到L2
// func (b *ZkSyncBridge) DepositETH(ctx context.Context, amount *big.Int, gasLimit uint64) (string, error) {
// 	// 获取L1桥接合约实例
// 	bridgeContract, err := bridge.NewBridge(b.bridgeAddr, b.l1Client.GetClient())
// 	if err != nil {
// 		return "", fmt.Errorf("创建桥接合约实例失败: %w", err)
// 	}

// 	// 获取存款者地址
// 	depositor := b.l2Client.address
// 	if depositor == (common.Address{}) {
// 		return "", fmt.Errorf("存款者地址未初始化")
// 	}

// 	// 构建存款交易
// 	txOpts, err := b.l1Client.GetClient().GetTransactor(ctx, b.l1Client.privateKey)
// 	if err != nil {
// 		return "", fmt.Errorf("创建交易选项失败: %w", err)
// 	}
// 	txOpts.Value = amount
// 	txOpts.GasLimit = gasLimit

// 	// 执行存款
// 	tx, err := bridgeContract.DepositETH(txOpts, depositor, [32]byte{})
// 	if err != nil {
// 		return "", fmt.Errorf("存款交易执行失败: %w", err)
// 	}

// 	return tx.Hash().Hex(), nil
// }

// // Withdraw 从L2提款到L1
// func (b *ZkSyncBridge) Withdraw(ctx context.Context, token common.Address, amount *big.Int) (string, error) {
// 	if b.l2Client.wallet == nil {
// 		return "", fmt.Errorf("L2钱包未初始化")
// 	}

// 	// 构建提款交易
// 	withdrawTx := accounts.WithdrawTransaction{
// 		Token:  token,
// 		Amount: amount,
// 		To:     b.l2Client.address,
// 	}

// 	// 发送提款交易
// 	tx, err := b.l2Client.wallet.Withdraw(ctx, withdrawTx)
// 	if err != nil {
// 		return "", fmt.Errorf("提款交易发送失败: %w", err)
// 	}

// 	return tx.Hash().Hex(), nil
// }

// // FinalizeWithdrawal 完成L1提款最终确认
// func (b *ZkSyncBridge) FinalizeWithdrawal(ctx context.Context, withdrawalHash common.Hash) (string, error) {
// 	// 获取L2交易详情
// 	l2Tx, err := b.l2Client.GetL2Client().TransactionByHash(ctx, withdrawalHash)
// 	if err != nil {
// 		return "", fmt.Errorf("获取L2交易失败: %w", err)
// 	}

// 	// 解析提款数据
// 	withdrawal, err := types.ParseWithdrawal(l2Tx.Data())
// 	if err != nil {
// 		return "", fmt.Errorf("解析提款数据失败: %w", err)
// 	}

// 	// 获取证明
// 	proof, err := b.l2Client.GetL2Client().GetWithdrawalProof(ctx, withdrawalHash)
// 	if err != nil {
// 		return "", fmt.Errorf("获取提款证明失败: %w", err)
// 	}

// 	// 获取L1桥接合约
// 	bridgeContract, err := bridge.NewBridge(b.bridgeAddr, b.l1Client.GetClient())
// 	if err != nil {
// 		return "", fmt.Errorf("创建桥接合约实例失败: %w", err)
// 	}

// 	// 构建交易选项
// 	txOpts, err := b.l1Client.GetClient().GetTransactor(ctx, b.l1Client.privateKey)
// 	if err != nil {
// 		return "", fmt.Errorf("创建交易选项失败: %w", err)
// 	}

// 	// 执行最终确认
// 	tx, err := bridgeContract.FinalizeWithdrawal(txOpts, *withdrawal, proof)
// 	if err != nil {
// 		return "", fmt.Errorf("最终确认提款失败: %w", err)
// 	}

// 	return tx.Hash().Hex(), nil
// }

// // GetWithdrawalStatus 查询提款状态
// func (b *ZkSyncBridge) GetWithdrawalStatus(ctx context.Context, withdrawalHash common.Hash) (types.WithdrawalStatus, error) {
// 	return b.l2Client.GetL2Client().GetWithdrawalStatus(ctx, withdrawalHash)
// }

// // IsWithdrawalFinalized 检查提款是否已在L1上最终确认
// func (b *ZkSyncBridge) IsWithdrawalFinalized(ctx context.Context, withdrawalHash common.Hash) (bool, error) {
// 	status, err := b.GetWithdrawalStatus(ctx, withdrawalHash)
// 	if err != nil {
// 		return false, err
// 	}
// 	return status == types.WithdrawalStatusFinalized, nil
// }
