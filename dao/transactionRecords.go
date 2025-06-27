package dao

import (
	"database/sql"
	"time"

	"MetaFarmBackend/component/db"
)

// TransactionRecords 交易记录表结构体
type TransactionRecords struct {
	ID                 int64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TxHash             string          `gorm:"column:tx_hash;size:66;uniqueIndex:idx_tx_hash" json:"tx_hash"`
	UserAddress        string          `gorm:"column:user_address;size:42;not null;index:idx_user_address" json:"user_address"`
	TransactionType    int8            `gorm:"column:transaction_type;not null;index:idx_transaction_type" json:"transaction_type"`
	NFTContractAddress string          `gorm:"column:nft_contract_address;size:42" json:"nft_contract_address"`
	TokenID            int64           `gorm:"column:token_id" json:"token_id"`
	Amount             sql.NullFloat64 `gorm:"column:amount;type:decimal(36,18);default:0.0" json:"amount"`
	Fee                sql.NullFloat64 `gorm:"column:fee;type:decimal(36,18);default:0.0" json:"fee"`
	Status             int8            `gorm:"column:status;not null;default:0" json:"status"`
	CreateTime         time.Time       `gorm:"column:create_time;not null;index:idx_create_time" json:"create_time"`
	UpdateTime         time.Time       `gorm:"column:update_time;not null;autoUpdateTime" json:"update_time"`
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
func GetTransactionByTxHash(txHash string) (*TransactionRecords, error) {
	var transaction TransactionRecords
	err := db.GetDB().Where("tx_hash = ?", txHash).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionRecordsByUser 根据用户地址获取交易记录列表
func GetTransactionRecordsByUser(userAddress string, page, pageSize int) ([]*TransactionRecords, int64, error) {
	var transactions []*TransactionRecords
	var total int64

	// 获取总数
	err := db.GetDB().Model(&TransactionRecords{}).Where("user_address = ?", userAddress).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = db.GetDB().Where("user_address = ?", userAddress).Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// CreateTransactionRecord 创建交易记录
func (t *TransactionRecords) CreateTransactionRecord() error {
	return db.GetDB().Create(t).Error
}

// UpdateTransactionStatus 更新交易状态
func (t *TransactionRecords) UpdateTransactionStatus(status int8) error {
	return db.GetDB().Model(t).Update("status", status).Error
}

// UpdateTransactionHash 更新交易哈希
func (t *TransactionRecords) UpdateTransactionHash(txHash string) error {
	return db.GetDB().Model(t).Update("tx_hash", txHash).Error
}
