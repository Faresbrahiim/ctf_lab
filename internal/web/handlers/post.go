package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"CTF/internal/models"
	"CTF/internal/storage"
)

type PostHandler struct {
	Tmpl *template.Template
}

// Create a new post
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	user := r.Context().Value("user").(*models.User)
	content := r.FormValue("content")
	if content == "" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	_, err := storage.DB.Exec("INSERT INTO posts (user_id, content) VALUES (?, ?)", user.ID, content)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// View a single post
func (h *PostHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	row := storage.DB.QueryRow("SELECT id, user_id, content, rowid FROM posts WHERE id = ?", id)

	var p models.Post
	if err := row.Scan(&p.ID, &p.ID, &p.Content, &p.CreatedAt); err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	user := r.Context().Value("user").(*models.User)

	h.Tmpl.ExecuteTemplate(w, "post.html", map[string]any{
		"User": user,
		"Post": p,
	})
}
