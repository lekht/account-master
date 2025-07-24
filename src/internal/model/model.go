package model

import "github.com/google/uuid"

// 1. id (uuid, unique)
// 2. email
// 3. username (unique)
// 4. password
// 5. admin (bool)

type Profile struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Admin    bool      `json:"admin"`
}
