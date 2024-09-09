package models

import "time"

type UserProfile struct {
	ID        uint64    `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	FirsName  string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Bio       string    `json:"bio"`
	Interests []string  `json:"interests"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID                 uint64    `json:"id"`
	// UserName           string  `json:"username"`
	Email              string    `json:"email"`
	Encrypted_Password string    `json:"password"`
	CreatedAt          time.Time `json:"created_at"`
	Verified           bool      `json:"verified"`
}

type TempUser struct {
	ID                 uint64 `json:"id"`
	// UserName           string  `json:"username"`
	Email              string    `json:"email"`
	Encrypted_Password string    `json:"-"`
	Token              string    `json:"-"`
	ExpiresAt          time.Time `jaon:"-"`
}
