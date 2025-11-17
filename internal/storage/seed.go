package storage

import "log"

func Seed() {
	users := []struct {
		username string
		password string
		role     string
	}{
		{"admin1", "pass123", "admin"},
		{"admin2", "pass123", "admin"},
		{"user1", "pass123", "user"},
	}

	for _, u := range users {
		_, err := DB.Exec(
			"INSERT OR IGNORE INTO users(username, password, role) VALUES(?,?,?)",
			u.username, u.password, u.role,
		)
		if err != nil {
			log.Fatal("Error inserting user:", err)
		}
	}

	_, err := DB.Exec(`
		INSERT OR IGNORE INTO posts(user_id, content)
		VALUES 
			(1, 'I think my friend admin2 also knows the flag...'),
			(2, 'Yes, I know it too… but I won''t say it publicly.'),
			(3, 'Hello everyone! Normal user here.')
	`)
	if err != nil {
		log.Fatal("Error inserting posts:", err)
	}

	log.Println("[DB] Seed data inserted")
}
