package orm

import (
	"database/sql"
	"log"
	"store-manager/internal/infrastructure/persistence/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(sqlDB *sql.DB) *gorm.DB {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao conectar com o banco via GORM: %v", err)
	}

	if err := gormDB.AutoMigrate(
		&models.ProductModel{},
		&models.RawMaterialModel{},
		&models.ProductRawMaterialModel{},
	); err != nil {
		log.Fatalf("Falha ao realizar a migration: %v", err)
	}

	return gormDB
}
