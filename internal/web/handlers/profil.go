package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"CTF/internal/models"
	"CTF/internal/storage"
)

type ProfileHandler struct {
	Tmpl *template.Template
}

// ViewProfile shows a user profile (vulnerable to SSTI)
func (h *ProfileHandler) ViewProfile(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	row := storage.DB.QueryRow("SELECT id, username, role, bio FROM users WHERE id = ?", id)

	var u models.User
	row.Scan(&u.ID, &u.Username, &u.Role, &u.Bio)

	current := r.Context().Value("user").(*models.User)

	// SSTI: Bio is executed as template with current admin/user context
	tmpl, _ := template.New("bio").Parse(u.Bio)
	tmpl.Execute(w, map[string]any{
		"User": current, // attacker can access current.User fields
	})
}
