package rds

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlConnect() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:HK75*YvR83gNy^U4@tcp(52.220.217.245:3307)/news24h?charset=utf8&parseTime=True&loc=Local", // data source name
		DefaultStringSize:         256,                                                                                            // default size for string fields
		DisableDatetimePrecision:  true,                                                                                           // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                           // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                           // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return db
}
