package dao

import "MetaFarmBackend/component/db"

func InitTable() {
	//让 GORM 根据结构体的定义，自动在数据库中创建对应的表结构，如果表已存在，会尝试更新表结构以匹配结构体定义
	db.DB.AutoMigrate(&User{})
	db.DB.AutoMigrate(&UserAccount{})
	db.DB.AutoMigrate(&UserAssetsSummary{})
	db.DB.AutoMigrate(&UserItems{})
	db.DB.AutoMigrate(&UserWallet{})
	db.DB.AutoMigrate(&LoginSession{})
	db.DB.AutoMigrate(&WalletLoginLog{})
	db.DB.AutoMigrate(&LandInfo{})
	db.DB.AutoMigrate(&LandActivity{})
	db.DB.AutoMigrate(&LandLayout{})
	db.DB.AutoMigrate(&LandMarket{})
	db.DB.AutoMigrate(&LandRental{})
	db.DB.AutoMigrate(&LandUpgrade{})
	db.DB.AutoMigrate(&MarketListings{})
	db.DB.AutoMigrate(&PlotPlanting{})
	db.DB.AutoMigrate(&TransactionRecords{})

}
