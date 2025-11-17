package main

import (
	"log"

	"CTF/internal/config"
	"CTF/internal/storage"
	web "CTF/internal/web/server"
)

func main() {
	log.Println("CTF Server booting...")

	cfg := config.Load()

	storage.Init(cfg)
	storage.Migrate()
	storage.Seed()

	log.Println("Database initialized ✓")
	log.Println("Starting HTTP server next...")

	// start web server
	server := web.NewServer()
	server.Start()
}
