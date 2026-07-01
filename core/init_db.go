package core

import (
	"clover_server/global"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func InitDB() *gorm.DB {
	dbconfig := global.Config.DB
	dbconfig1 := global.Config.DB1

	// 数据库日志配置
	gormLogger := NewSlogGormLogger()
	if dbconfig.Debug {
		gormLogger = gormLogger.LogMode(logger.Info)
	}

	// 主库连接
	db, err := gorm.Open(mysql.Open(dbconfig.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormLogger,
	})
	if err != nil {
		slog.Error("数据库主库连接失败", slog.Any("err", err))
		os.Exit(1)
	}

	// 配置读写分离，DB 作为主库 (Source)，DB1 作为从库 (Replica)
	resolver := dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(dbconfig.DSN())},
		Replicas: []gorm.Dialector{mysql.Open(dbconfig1.DSN())},
		Policy:   dbresolver.RandomPolicy{},
	}).
		SetMaxIdleConns(10).
		SetMaxOpenConns(100).
		SetConnMaxLifetime(10 * time.Second)

	err = db.Use(resolver)
	if err != nil {
		slog.Error("配置数据库读写分离失败", slog.Any("err", err))
		os.Exit(1)
	}

	// 配置主库基础连接池（gorm 默认连接池，针对未命中 dbresolver 的情况）
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(10 * time.Second)
	}
	slog.Info("数据库连接成功")
	return db
}
