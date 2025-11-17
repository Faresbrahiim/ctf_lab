package models

type Message struct {
	ID       int    // message ID
	FromUser int    // sender's user ID
	ToUser   int    // recipient's user ID
	Content  string // message content
}
