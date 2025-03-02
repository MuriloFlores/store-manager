package connection

import (
	"database/sql"
	"fmt"
	"log"
	"store-manager/internal/infrastructure/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectPostgresDB() *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.EnvConfigs.DbUsername,
		config.EnvConfigs.DbPassword,
		config.EnvConfigs.DbHost,
		config.EnvConfigs.DbPort,
		config.EnvConfigs.DbDatabase,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Falha ao conectar com o banco (database/sql): %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao fazer ping no banco de dados: %v", err)
	}

	return db
}
