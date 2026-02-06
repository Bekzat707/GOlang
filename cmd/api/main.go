package main

import (
	"log"
	"net/http"

	"github.com/bekzatsaparbekov/task-api/internal/handlers"
	"github.com/bekzatsaparbekov/task-api/internal/middleware"
	"github.com/bekzatsaparbekov/task-api/internal/storage"
)

func main() {
	taskStorage := storage.NewTaskStorage()
	taskHandler := handlers.NewTaskHandler(taskStorage)
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		case http.MethodPatch:
			taskHandler.UpdateTask(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	handler := middleware.Logger(middleware.APIKeyAuth(mux))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
