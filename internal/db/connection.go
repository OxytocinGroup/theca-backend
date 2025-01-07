package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	config "github.com/OxytocinGroup/theca-backend/internal/config"
	domain "github.com/OxytocinGroup/theca-backend/internal/domain"
)

type Database interface {
	AutoMigrate(dst ...interface{}) error
	GetDB() *gorm.DB
}

type GormDatabase struct {
	Conn *gorm.DB
}

func (g *GormDatabase) AutoMigrate(dst ...interface{}) error {
	return g.Conn.AutoMigrate(dst...)
}

func (g *GormDatabase) GetDB() *gorm.DB {
	return g.Conn
}

func ConnectDatabase(cfg config.Config) Database {
    psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s sslmode=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword, "disable")
    log.Printf("Connecting to database with: %s", psqlInfo)
    
    conn, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
        SkipDefaultTransaction: true,
    })
    if dbErr != nil {
        log.Fatalf("Failed to connect to database: %v", dbErr)
    }

    db := &GormDatabase{Conn: conn}
    if err := db.AutoMigrate(&domain.User{}, &domain.Session{}, &domain.Bookmark{}); err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }
    return db
}