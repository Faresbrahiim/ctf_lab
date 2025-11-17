package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"CTF/internal/models"
	"CTF/internal/storage"
)

type DashboardHandler struct {
	Tmpl *template.Template
}

// Dashboard shows all posts, links to chat, profile, and logout
func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	fmt.Print("sksskksks")
	rows, err := storage.DB.Query(`
		SELECT id, user_id, content, rowid
		FROM posts
		ORDER BY rowid DESC
	`)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.ID, &p.Content, &p.CreatedAt)
		posts = append(posts, p)
	}

	h.Tmpl.ExecuteTemplate(w, "dashboard.html", map[string]any{
		"User":  user,
		"Posts": posts,
	})
}
