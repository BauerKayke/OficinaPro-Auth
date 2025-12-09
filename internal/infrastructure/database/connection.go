package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Driver PostgreSQL
)

// Config representa a configuração do banco de dados
type Config struct {
	Host            string
	Port            string
	Database        string
	User            string
	Password        string
	SSLMode         string
	MaxConnections  int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewConnection cria uma nova conexão com o banco de dados
func NewConnection(config Config) (*sql.DB, error) {
	// Montar connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.Host,
		config.Port,
		config.Database,
		config.User,
		config.Password,
		config.SSLMode,
	)

	log.Printf("[Database] Conectando ao PostgreSQL: %s:%s/%s", config.Host, config.Port, config.Database)

	// Abrir conexão
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(config.MaxConnections)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	log.Println("[Database] Conexão estabelecida com sucesso")
	return db, nil
}

// Close fecha a conexão com o banco de dados
func Close(db *sql.DB) error {
	if db != nil {
		log.Println("[Database] Fechando conexão")
		return db.Close()
	}
	return nil
}
