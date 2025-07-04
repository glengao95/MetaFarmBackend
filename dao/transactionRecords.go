package dao

import (
	"context"
	"database/sql"
	"time"
)

// TransactionRecords 交易记录表结构体
type TransactionRecords struct {
	ID                 int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                        // 主键ID
	TxHash             string          `gorm:"column:tx_hash;size:66;uniqueIndex:idx_tx_hash" json:"tx_hash"`                       // 交易哈希
	UserAddress        string          `gorm:"column:user_address;size:42;not null;index:idx_user_address" json:"user_address"`     // 用户钱包地址
	TransactionType    int8            `gorm:"column:transaction_type;not null;index:idx_transaction_type" json:"transaction_type"` // 交易类型(1:购买, 2:出售, 3:转账等)
	NFTContractAddress string          `gorm:"column:nft_contract_address;size:42" json:"nft_contract_address"`                     // NFT合约地址
	TokenID            int64           `gorm:"column:token_id" json:"token_id"`                                                     // NFT TokenID
	Amount             sql.NullFloat64 `gorm:"column:amount;type:decimal(36,18);default:0.0" json:"amount"`                         // 交易金额
	Fee                sql.NullFloat64 `gorm:"column:fee;type:decimal(36,18);default:0.0" json:"fee"`                               // 交易手续费
	Status             int8            `gorm:"column:status;not null;default:0" json:"status"`                                      // 状态(0:处理中, 1:成功, 2:失败)
	CreateTime         time.Time       `gorm:"column:create_time;not null;index:idx_create_time" json:"create_time"`                // 创建时间
	UpdateTime         time.Time       `gorm:"column:update_time;not null;autoUpdateTime" json:"update_time"`                       // 更新时间
}

// TableName 设置表名
func (t *TransactionRecords) TableName() string {
	return "transaction_records"
}

// NewTransactionRecord 创建新的交易记录实例
func NewTransactionRecord(userAddress string, transactionType int8) *TransactionRecords {
	now := time.Now()
	return &TransactionRecords{
		UserAddress:     userAddress,
		TransactionType: transactionType,
		Status:          0, // 默认状态为处理中
		CreateTime:      now,
		UpdateTime:      now,
	}
}

// GetTransactionByTxHash 根据交易哈希获取交易记录
func (dao *Dao) GetTransactionByTxHash(ctx context.Context, txHash string) (*TransactionRecords, error) {
	var transaction TransactionRecords
	err := dao.DB.WithContext(ctx).Where("tx_hash = ?", txHash).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionRecordsByUser 根据用户地址获取交易记录列表
func (dao *Dao) GetTransactionRecordsByUser(ctx context.Context, userAddress string, page, pageSize int) ([]*TransactionRecords, int64, error) {
	var transactions []*TransactionRecords
	var total int64

	// 获取总数
	err := dao.DB.WithContext(ctx).Model(&TransactionRecords{}).Where("user_address = ?", userAddress).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = dao.DB.WithContext(ctx).Where("user_address = ?", userAddress).Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// CreateTransactionRecord 创建交易记录
func (dao *Dao) CreateTransactionRecord(ctx context.Context, t *TransactionRecords) error {
	return dao.DB.WithContext(ctx).Create(t).Error
}

// UpdateTransactionStatus 更新交易状态
func (dao *Dao) UpdateTransactionStatus(ctx context.Context, t *TransactionRecords, status int8) error {
	return dao.DB.WithContext(ctx).Model(t).Update("status", status).Error
}

// UpdateTransactionHash 更新交易哈希
func (dao *Dao) UpdateTransactionHash(ctx context.Context, t *TransactionRecords, txHash string) error {
	return dao.DB.WithContext(ctx).Model(t).Update("tx_hash", txHash).Error
}
