package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"net/http"
	"encoding/json"
	"os"
	_ "modernc.org/sqlite"
)

var db *sql.DB

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

func main() {
	var err error
	db, err = sql.Open("sqlite", "tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	initDB()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", listTasksHandler)
	mux.HandleFunc("POST /tasks", createTaskHandler)
	mux.HandleFunc("DELETE /tasks/{id}", deleteTaskHandler)

	port := os.Getenv("PORT")
if port == "" {
	port = "8080"
}
fmt.Println("Server starting on :" + port + "...")
http.ListenAndServe(":"+port, withCORS(mux))
}

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")

	_, err := db.Exec("DELETE FROM tasks WHERE due_date < ?", today)
	if err != nil {
		http.Error(w, "failed to clean up expired tasks", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, description, due_date FROM tasks ORDER BY due_date ASC")
	if err != nil {
		http.Error(w, "failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Description, &t.DueDate)
		if err != nil {
			http.Error(w, "failed to read task", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"INSERT INTO tasks (description, due_date) VALUES (?, ?)",
		t.Description, t.DueDate,
	)
	if err != nil {
		http.Error(w, "failed to save task", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	t.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func initDB() {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		due_date TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

var allowedOrigins = map[string]bool{
	"http://localhost:5173":              true,
	"https://neon-banoffee-ee24de.netlify.app": true,
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}