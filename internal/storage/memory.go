package storage

import (
	"database/sql"
	"log"

	"CTF/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(cfg *config.Config) {
	var err error

	DB, err = sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to open SQLite DB:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to ping SQLite DB:", err)
	}

	log.Println("[DB] Connected to SQLite:", cfg.DBPath)
}

func Migrate() {
	schema := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			password TEXT,
			role TEXT,
			bio TEXT

		);`,
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			content TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_user INTEGER,
			to_user INTEGER,
			content TEXT,
			FOREIGN KEY(from_user) REFERENCES users(id),
			FOREIGN KEY(to_user) REFERENCES users(id)
		);`,
	}

	for _, stmt := range schema {
		if _, err := DB.Exec(stmt); err != nil {
			log.Fatal("Migration error:", err)
		}
	}

	log.Println("[DB] Migrations OK")
}
