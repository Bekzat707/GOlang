package modules

import "time"

type User struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"password"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
