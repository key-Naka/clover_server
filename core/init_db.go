package core

import (
	"clover_server/global"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dbconfig := global.Config.DB
	db, err := gorm.Open(mysql.Open(dbconfig.DSN()), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		slog.Error("数据库连接失败", slog.Any("err", err))
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)
	return db
}
