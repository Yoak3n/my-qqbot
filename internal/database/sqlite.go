package database

import (
	"fmt"
	"github.com/Yoak3n/gulu/logger"
	"github.com/Yoak3n/gulu/util"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func initSqlite() *gorm.DB {
	_ = util.CreateDirNotExists("./data/db")
	dsn := "./data/db/bot.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("database connected err:%v", err))
	}
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("database connected err:%v", err))
	}
	return db
}
