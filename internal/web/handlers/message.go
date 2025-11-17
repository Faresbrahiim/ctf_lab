package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"CTF/internal/models"
	"CTF/internal/storage"
)

type MessageHandler struct {
	Tmpl *template.Template
}

// Inbox shows messages to the logged-in user
func (h *MessageHandler) Inbox(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	rows, err := storage.DB.Query("SELECT id, from_user, to_user, content FROM messages WHERE to_user = ?", user.ID)
	if err != nil {
		http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		rows.Scan(&m.ID, &m.FromUser, &m.ToUser, &m.Content)
		msgs = append(msgs, m)
	}

	h.Tmpl.ExecuteTemplate(w, "message_inbox.html", map[string]any{
		"User":     user,
		"Messages": msgs,
	})
}

// Send allows sending a message (vulnerable IDOR)
func (h *MessageHandler) Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/messages", http.StatusSeeOther)
		return
	}

	fromStr := r.FormValue("from") // ATTACKER CAN IMPERSONATE
	toStr := r.FormValue("to")
	content := r.FormValue("content")

	fromID, _ := strconv.Atoi(fromStr)
	toID, _ := strconv.Atoi(toStr)

	_, err := storage.DB.Exec("INSERT INTO messages (from_user, to_user, content) VALUES (?, ?, ?)", fromID, toID, content)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/messages", http.StatusSeeOther)
}
