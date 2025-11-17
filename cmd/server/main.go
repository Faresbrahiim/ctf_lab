package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"CTF/internal/config"
	"CTF/internal/models"
	"CTF/internal/storage"
	"CTF/internal/web/handlers"
)

func main() {
	// --- Load config ---
	cfg := &config.Config{
		DBPath: "ctf.db",
	}

	// --- Init DB ---
	storage.Init(cfg)
	storage.Migrate()
	seedUsersAndPosts()

	// --- Load templates ---
	tmpl := template.Must(template.ParseGlob("internal/web/templates/*.html"))

	// --- Handlers ---
	dashboardHandler := &handlers.DashboardHandler{Tmpl: tmpl}
	postHandler := &handlers.PostHandler{Tmpl: tmpl}
	profileHandler := &handlers.ProfileHandler{Tmpl: tmpl}
	loginHandler := &handlers.LoginHandler{Tmpl: tmpl}
	messageHandler := &handlers.MessageHandler{Tmpl: tmpl}

	// --- Routes ---
	mux := http.NewServeMux()
	mux.Handle("/login", loginHandler)
	mux.Handle("/dashboard", withAuth(dashboardHandler.ServeHTTP))
	mux.Handle("/post/create", withAuth(postHandler.CreatePost))
	mux.Handle("/post/view", withAuth(postHandler.ViewPost))
	mux.Handle("/profile", withAuth(profileHandler.ViewProfile))
	mux.Handle("/messages", withAuth(messageHandler.Inbox))
	mux.Handle("/message/send", withAuth(messageHandler.Send))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// --- Start server ---
	addr := ":8080"
	log.Printf("Server started at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// --- Simple auth middleware ---
func withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Very simple: user ID stored in query param for demo
		idStr := r.URL.Query().Get("user_id")
		if idStr == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		id, _ := strconv.Atoi(idStr)

		row := storage.DB.QueryRow("SELECT id, username, role, bio FROM users WHERE id = ?", id)
		var u models.User
		if err := row.Scan(&u.ID, &u.Username, &u.Role, &u.Bio); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := contextWithUser(r.Context(), &u)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

// --- Context helpers ---
type contextKey string

func contextWithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, contextKey("user"), user)
}

// --- Seed 6 admins + 6 users ---
func seedUsersAndPosts() {
	users := []models.User{
		{Username: "admin1", Password: "pass", Role: "admin", Bio: "Admin 1 bio"},
		{Username: "admin2", Password: "pass", Role: "admin", Bio: "Admin 2 bio"},
		{Username: "admin3", Password: "pass", Role: "admin", Bio: "Admin 3 bio"},
		{Username: "user1", Password: "pass", Role: "user", Bio: "User 1 bio"},
		{Username: "user2", Password: "pass", Role: "user", Bio: "User 2 bio"},
		{Username: "user3", Password: "pass", Role: "user", Bio: "User 3 bio"},
	}

	for _, u := range users {
		_, err := storage.DB.Exec("INSERT OR IGNORE INTO users (username, password, role, bio) VALUES (?, ?, ?, ?)",
			u.Username, u.Password, u.Role, u.Bio)
		if err != nil {
			log.Println("Seed user error:", err)
		}
	}

	// Optional: seed sample posts
	for i := 1; i <= 6; i++ {
		storage.DB.Exec("INSERT INTO posts (user_id, content) VALUES (?, ?)", i, fmt.Sprintf("Sample post by user %d", i))
	}
}
