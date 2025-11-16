package main

import (
	"log"

	"CTF/internal/config"
	"CTF/internal/storage"
)

func main() {
	log.Println("CTF Server booting...")

	cfg := config.Load()

	storage.Init(cfg)
	storage.Migrate()
	storage.Seed()

	log.Println("Database initialized ✓")
	log.Println("Starting HTTP server next...")
}
