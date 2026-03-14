package repository

import (
"database/sql"
"fmt"
"strings"

"github.com/google/uuid"

"practice5/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// allowedColumns whitelists sortable/filterable column names to prevent SQL injection.
var allowedColumns = map[string]string{
	"id":         "u.id",
	"name":       "u.name",
	"email":      "u.email",
	"gender":     "u.gender",
	"birth_date": "u.birth_date",
}

// GetPaginatedUsers returns a paginated, filtered, and sorted list of users.
func (r *Repository) GetPaginatedUsers(p models.FilterParams) (models.PaginatedResponse, error) {
	args := []interface{}{}
	argIdx := 1

	conditions := []string{}

	if p.ID != nil && *p.ID != uuid.Nil {
		conditions = append(conditions, fmt.Sprintf("u.id = $%d", argIdx))
		args = append(args, *p.ID)
		argIdx++
	}

	if p.Name != nil && *p.Name != "" {
		conditions = append(conditions, fmt.Sprintf("u.name ILIKE $%d", argIdx))
		args = append(args, "%"+*p.Name+"%")
		argIdx++
	}

	if p.Email != nil && *p.Email != "" {
		conditions = append(conditions, fmt.Sprintf("u.email ILIKE $%d", argIdx))
		args = append(args, "%"+*p.Email+"%")
		argIdx++
	}

	if p.Gender != nil && *p.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("u.gender = $%d", argIdx))
		args = append(args, *p.Gender)
		argIdx++
	}

	if p.BirthDate != nil {
		conditions = append(conditions, fmt.Sprintf("u.birth_date = $%d", argIdx))
		args = append(args, *p.BirthDate)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Validate and resolve ORDER BY column
	orderCol := "u.id"
	if col, ok := allowedColumns[strings.ToLower(p.OrderBy)]; ok {
		orderCol = col
	}

	orderDir := "ASC"
	if strings.ToUpper(p.OrderDir) == "DESC" {
		orderDir = "DESC"
	}

	// Count total matching records
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users u %s`, whereClause)
	var totalCount int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return models.PaginatedResponse{}, fmt.Errorf("count query: %w", err)
	}

	// Fetch paginated records
	offset := (p.Page - 1) * p.PageSize
	query := fmt.Sprintf(
`SELECT u.id, u.name, u.email, u.gender, u.birth_date
		 FROM users u
		 %s
		 ORDER BY %s %s
		 LIMIT $%d OFFSET $%d`,
whereClause, orderCol, orderDir, argIdx, argIdx+1,
)
	args = append(args, p.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return models.PaginatedResponse{}, fmt.Errorf("paginated query: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return models.PaginatedResponse{}, fmt.Errorf("scan: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return models.PaginatedResponse{}, fmt.Errorf("rows iteration: %w", err)
	}

	if users == nil {
		users = []models.User{}
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       p.Page,
		PageSize:   p.PageSize,
	}, nil
}

// GetCommonFriends returns friends shared by user1 and user2 using a single JOIN query (no N+1).
func (r *Repository) GetCommonFriends(userID1, userID2 uuid.UUID) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends uf1 ON uf1.friend_id = u.id AND uf1.user_id = $1
		JOIN user_friends uf2 ON uf2.friend_id = u.id AND uf2.user_id = $2
		ORDER BY u.id
	`

	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("common friends query: %w", err)
	}
	defer rows.Close()

	var friends []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		friends = append(friends, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if friends == nil {
		friends = []models.User{}
	}
	return friends, nil
}
