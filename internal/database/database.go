package database

import (
	"database/sql"
	"github.com/Yoak3n/gulu/logger"
	"gorm.io/gorm"
	"time"
)

var conn *sql.DB
var mdb *gorm.DB

func init() {
	mdb = initSqlite()
	conn, _ = mdb.DB()

	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(time.Hour)
}

func migrateTables(tables ...interface{}) {
	err := mdb.AutoMigrate(tables...)
	if err != nil {
		logger.Logger.Panic("migrate tables failed")
	}
}

func GetDB() *gorm.DB {
	return mdb
}
func GetConn() *sql.DB {
	return conn
}
