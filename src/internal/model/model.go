package model

// 1. id (uuid, unique)
// 2. email
// 3. username (unique)
// 4. password
// 5. admin (bool)

type Profile struct {
	Id       int
	Email    string
	Username string
	Password string
	Admin    bool
}
