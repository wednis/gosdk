package gosdk

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 进行默认初始化
func initDefaultDB(db *gorm.DB) {
	sqldb, _ := db.DB()
	sqldb.SetMaxIdleConns(25)                  // 最大空闲连接数
	sqldb.SetMaxOpenConns(100)                 // 最大打开连接数
	sqldb.SetConnMaxIdleTime(5 * time.Minute)  // 空闲连接超时 5 分钟
	sqldb.SetConnMaxLifetime(30 * time.Minute) // 连接最大生命周期 30 分钟
}

func NewMysqlGorm(username string, password string, dbname string, config *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbname)), config)
	if err != nil {
		return nil, err
	}
	initDefaultDB(db)
	return db, nil
}

func NewSqlite3Gorm(path string, config *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), config)
	if err != nil {
		return nil, err
	}
	initDefaultDB(db)
	return db, nil
}
