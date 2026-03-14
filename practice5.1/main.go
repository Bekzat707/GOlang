package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 54320
	user     = "postgres"
	password = "password" 
	dbname   = "practice51"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	fmt.Println("Successfully connected to database!")

	err = seedDatabase(db)
	if err != nil {
		log.Fatal("Could not seed database:", err)
	}

	repo := NewRepository(db)
	router := NewRouter(repo)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func seedDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			gender VARCHAR(50),
			birth_date DATE
		);
		CREATE TABLE IF NOT EXISTS user_friends (
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			friend_id UUID REFERENCES users(id) ON DELETE CASCADE,
			PRIMARY KEY (user_id, friend_id),
			CHECK (user_id != friend_id)
		);
	`)
	if err != nil {
		return fmt.Errorf("creating tables: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("counting users: %v", err)
	}

	if count >= 20 {
		fmt.Println("Database already seeded.")
		return nil
	}

	fmt.Println("Seeding database...")

	_, _ = db.Exec("TRUNCATE TABLE user_friends CASCADE")
	_, _ = db.Exec("TRUNCATE TABLE users CASCADE")

	var userIDs []string
	for i := 1; i <= 20; i++ {
		var id string
		err := db.QueryRow(`
			INSERT INTO users (name, email, gender, birth_date)
			VALUES ($1, $2, $3, $4) RETURNING id
		`, fmt.Sprintf("User%d", i), fmt.Sprintf("user%d@example.com", i),
			func() string {
				if i%2 == 0 {
					return "Male"
				}
				return "Female"
			}(),
			time.Now().AddDate(-20, 0, -i).Format("2006-01-02")).Scan(&id)
		if err != nil {
			return err
		}
		userIDs = append(userIDs, id)
	}

	friendsPairs := [][]int{
		{0, 2}, {0, 3}, {0, 4},
		{1, 2}, {1, 3}, {1, 4},
		{2, 0}, {3, 0}, {4, 0},
		{2, 1}, {3, 1}, {4, 1},
	}

	for _, p := range friendsPairs {
		_, err := db.Exec(`
			INSERT INTO user_friends (user_id, friend_id)
			VALUES ($1, $2) ON CONFLICT DO NOTHING
		`, userIDs[p[0]], userIDs[p[1]])
		if err != nil {
			return err
		}
	}

	fmt.Println("Database successfully seeded.")
	return nil
}
