package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type UUID string

type User struct {
	ID        UUID      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birthDate"`
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type UserFilter struct {
	ID        string
	Name      string
	Email     string
	Gender    string
	BirthDate string
	OrderBy   string 	
}

func (r *Repository) GetPaginatedUsers(page int, pageSize int, filter UserFilter) (PaginatedResponse, error) {
	var users []User
	offset := (page - 1) * pageSize

	baseQuery := "SELECT id, name, email, gender, birth_date FROM users"
	countBaseQuery := "SELECT COUNT(*) FROM users"
	
	conditions := []string{}
	args := []interface{}{}
	argId := 1

	if filter.ID != "" {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argId))
		args = append(args, filter.ID)
		argId++
	}
	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argId))
		args = append(args, "%"+filter.Name+"%")
		argId++
	}
	if filter.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argId))
		args = append(args, "%"+filter.Email+"%")
		argId++
	}
	if filter.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", argId))
		args = append(args, filter.Gender)
		argId++
	}
	if filter.BirthDate != "" {
		conditions = append(conditions, fmt.Sprintf("birth_date = $%d", argId))
		args = append(args, filter.BirthDate)
		argId++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int
	err := r.db.QueryRow(countBaseQuery+whereClause, args...).Scan(&totalCount)
	if err != nil {
		return PaginatedResponse{}, err
	}
	orderClause := " ORDER BY id" 
	
	validColumns := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"gender":     true,
		"birth_date": true,
	}

	if filter.OrderBy != "" {
		parts := strings.Split(filter.OrderBy, " ")
		col := strings.ToLower(parts[0])
		dir := "ASC"
		if len(parts) > 1 && strings.ToUpper(parts[1]) == "DESC" {
			dir = "DESC"
		}
		
		if validColumns[col] {
			orderClause = fmt.Sprintf(" ORDER BY %s %s", col, dir)
		}
	}

	query := baseQuery + whereClause + orderClause + fmt.Sprintf(" LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return PaginatedResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	if users == nil {
		users = []User{}
	}

	return PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(userID1, userID2 string) ([]User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends uf1 ON u.id = uf1.friend_id
		JOIN user_friends uf2 ON u.id = uf2.friend_id
		WHERE uf1.user_id = $1 AND uf2.user_id = $2
	`
	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}

	if friends == nil {
		friends = []User{}
	}

	return friends, nil
}
