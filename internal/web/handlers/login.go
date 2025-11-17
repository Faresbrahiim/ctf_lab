package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"CTF/internal/models"
	"CTF/internal/storage"
)

type LoginHandler struct {
	Tmpl *template.Template
}

// Show login page and handle login
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.Tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	row := storage.DB.QueryRow("SELECT id, username, role, bio FROM users WHERE username=? AND password=?", username, password)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Role, &u.Bio)
	if err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}
	fmt.Print("**********************")
	// Save user in context/session middleware (simplified)
	http.Redirect(w, r, fmt.Sprintf("/dashboard?user_id=%d", u.ID), http.StatusSeeOther)
}
