package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
	repo *Repository
}

func NewRouter(repo *Repository) *Router {
	r := &Router{
		Router: mux.NewRouter(),
		repo:   repo,
	}

	r.HandleFunc("/users", r.getUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}/common-friends/{friendId}", r.getCommonFriendsHandler).Methods("GET")

	return r
}

func (r *Router) getUsersHandler(w http.ResponseWriter, req *http.Request) {
	pageStr := req.URL.Query().Get("page")
	pageSizeStr := req.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	filter := UserFilter{
		ID:        req.URL.Query().Get("id"),
		Name:      req.URL.Query().Get("name"),
		Email:     req.URL.Query().Get("email"),
		Gender:    req.URL.Query().Get("gender"),
		BirthDate: req.URL.Query().Get("birthDate"),
		OrderBy:   req.URL.Query().Get("orderBy"),
	}

	resp, err := r.repo.GetPaginatedUsers(page, pageSize, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (r *Router) getCommonFriendsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID1 := vars["id"]
	userID2 := vars["friendId"]

	if userID1 == "" || userID2 == "" {
		http.Error(w, "missing id or friendId", http.StatusBadRequest)
		return
	}

	friends, err := r.repo.GetCommonFriends(userID1, userID2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}
