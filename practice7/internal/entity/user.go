package entity

import "github.com/google/uuid"

type User struct {
    ID       uuid.UUID `json:"ID" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Username string    `json:"username"`
    Email    string    `json:"email"`
    Password string    `json:"password"`
    Role     string    `json:"role"` // user, admin, etc.
    Verified bool      `json:"verified"`
}
