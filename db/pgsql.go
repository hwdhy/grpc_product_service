package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"product_service"
	"product_service/models"
	"time"
)

var PgsqlDB *gorm.DB

// InitConnectionPgsql 数据库连接
func InitConnectionPgsql() {
	pDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			product_service.PgsqlUsername, product_service.PgsqlPassword, product_service.PgsqlDbname, product_service.PgsqlHost, product_service.PgsqlPort),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("connect pgsql db err: %v", err)
	}
	DB, err := pDB.DB()
	if err != nil {
		logrus.Fatalf("connect pgsql db err: %v", err)
	}
	DB.SetMaxOpenConns(100)
	DB.SetConnMaxIdleTime(10)
	DB.SetConnMaxLifetime(time.Minute)

	PgsqlDB = pDB
	_ = AutoMigrate()
}

// AutoMigrate 自动建表
func AutoMigrate() error {
	return PgsqlDB.AutoMigrate(
		&models.Product{},
		&models.Category{},
	)
}
