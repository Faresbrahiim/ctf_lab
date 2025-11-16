package models

type User struct {
	ID       int
	Username string
	Role     string // "admin" or "user"
	Email    string
	Secret   string // private info (leaked by SSTI)
}
