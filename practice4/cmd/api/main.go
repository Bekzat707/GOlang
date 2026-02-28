package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

var db *sql.DB

func main() {
	var err error

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Connected to the database.")

	mux := http.NewServeMux()
	mux.HandleFunc("/movies", moviesHandler)
	mux.HandleFunc("/movies/", movieByIDHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Println("Starting the Server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen Error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	fmt.Println("\nShutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Fatalf("DB Close Error: %v", err)
	}

	fmt.Println("Server exited properly.")
}

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, title, director, year FROM movies ORDER BY id ASC")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var movies []Movie
		for rows.Next() {
			var m Movie
			if err := rows.Scan(&m.ID, &m.Title, &m.Director, &m.Year); err != nil {
				continue
			}
			movies = append(movies, m)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movies)

	case http.MethodPost:
		var m Movie
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := db.QueryRow(
			"INSERT INTO movies (title, director, year) VALUES ($1, $2, $3) RETURNING id",
			m.Title, m.Director, m.Year,
		).Scan(&m.ID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(m)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func movieByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/movies/"):]
	if idStr == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		var m Movie
		err := db.QueryRow("SELECT id, title, director, year FROM movies WHERE id = $1", id).
			Scan(&m.ID, &m.Title, &m.Director, &m.Year)
		if err == sql.ErrNoRows {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)

	case http.MethodPut, "UPDATE": 
		var m Movie
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := db.Exec(
			"UPDATE movies SET title = $1, director = $2, year = $3 WHERE id = $4",
			m.Title, m.Director, m.Year, id,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, "Movie not found or not updated", http.StatusNotFound)
			return
		}

		m.ID = id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)

	case http.MethodDelete:
		res, err := db.Exec("DELETE FROM movies WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		rowsAffected, err := res.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, "Movie not found or not deleted", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"message\": \"Movie %d deleted successfully\"}", id)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
