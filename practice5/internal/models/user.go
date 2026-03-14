package models

import (
"time"

"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}

type FilterParams struct {
	ID        *uuid.UUID
	Name      *string
	Email     *string
	Gender    *string
	BirthDate *time.Time
	OrderBy   string // column name
	OrderDir  string // ASC or DESC
	Page      int
	PageSize  int
}

type CommonFriendsRequest struct {
	UserID1 uuid.UUID `json:"user_id_1"`
	UserID2 uuid.UUID `json:"user_id_2"`
}
