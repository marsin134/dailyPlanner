package database

import (
	"dailyPlanner/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

type MethodsDB interface {
	Close()
	RunMigrations(migrationFilePath string) error
	HealthCheck() error
	GetDB() *DB
}

type DB struct {
	*sqlx.DB
}

func ConnectDB(cfg *config.Config) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.DbHost, cfg.DB.DbPort, cfg.DB.DbUser, cfg.DB.DbPass, cfg.DB.DbName, cfg.DB.DbSSLMode)

	log.Printf("Connect for DB: host=%s, dbname=%s\n", cfg.DB.DbPort, cfg.DB.DbName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to the database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error checking the connection to the database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)

	dbStruct := DB{db}

	log.Println("Successful connection to PostgreSQL")
	return &dbStruct, nil
}

func (db *DB) Close() {
	db.DB.Close()
}

func (db *DB) RunMigrations(migrationFilePath string) error {
	if _, err := os.Stat(migrationFilePath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", migrationFilePath)
	}

	migration, err := os.ReadFile(migrationFilePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", migrationFilePath, err)
	}

	log.Printf("Attempting to run migration %s\n", migrationFilePath)

	_, err = db.Exec(string(migration))
	if err != nil {
		return fmt.Errorf("error executing migration %s: %w", migrationFilePath, err)
	}

	log.Println("Successfully run migration")
	return nil
}

func (db *DB) HealthCheck() error {
	if db == nil {
		return fmt.Errorf("connection to the database is not initialized")
	}

	return db.Ping()
}

func (db *DB) GetDB() *DB {
	return db
}

// psql -U postgres -h localhost
// CREATE DATABASE dailyPlanner;
// psql -h localhost -U postgres -d dailyPlanner -c "\dt"
// psql -h localhost -U postgres -d dailyPlanner -f migrations/001_create_tables.sql
// DROP DATABASE IF EXISTS dailyPlanner;
// CREATE DATABASE dailyPlanner;
