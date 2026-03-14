package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"practice5/internal/models"
	"practice5/internal/repository"
)

type UserHandler struct {
	repo *repository.Repository
}

func NewUserHandler(repo *repository.Repository) *UserHandler {
	return &UserHandler{repo: repo}
}

// writeJSON is a small helper that sends a JSON response.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError sends a JSON error message.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /users
// Query params: page, page_size, order_by, order_dir,
//
//	id, name, email, gender, birth_date (YYYY-MM-DD)
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Pagination
	page := parseIntDefault(q.Get("page"), 1)
	pageSize := parseIntDefault(q.Get("page_size"), 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Build filter params
	fp := models.FilterParams{
		OrderBy:  q.Get("order_by"),
		OrderDir: q.Get("order_dir"),
		Page:     page,
		PageSize: pageSize,
	}

	if idStr := q.Get("id"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			fp.ID = &id
		}
	}

	if name := q.Get("name"); name != "" {
		fp.Name = &name
	}

	if email := q.Get("email"); email != "" {
		fp.Email = &email
	}

	if gender := q.Get("gender"); gender != "" {
		fp.Gender = &gender
	}

	if bdStr := q.Get("birth_date"); bdStr != "" {
		if t, err := time.Parse("2006-01-02", bdStr); err == nil {
			fp.BirthDate = &t
		}
	}

	result, err := h.repo.GetPaginatedUsers(fp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GET /users/common-friends?user_id_1=XYZ&user_id_2=ABC
func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	id1, err1 := uuid.Parse(q.Get("user_id_1"))
	id2, err2 := uuid.Parse(q.Get("user_id_2"))

	if err1 != nil || err2 != nil || id1 == id2 {
		writeError(w, http.StatusBadRequest, "provide two different valid UUIDs for user_id_1 and user_id_2")
		return
	}

	friends, err := h.repo.GetCommonFriends(id1, id2)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id_1":      id1,
		"user_id_2":      id2,
		"common_friends": friends,
		"count":          len(friends),
	})
}

func parseIntDefault(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}
