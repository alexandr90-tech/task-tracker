package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alexandr90-tech/task-tracker/internal/task"
)

func main() {

	store := task.NewStore()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is alive")
	})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tasks := store.ListTasks()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)
		case "POST":
			var input struct {
				Title string `json:"title"`
			}

			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil || input.Title == "" {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			task := store.AddTask(input.Title)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(task)
		case "PUT":
			idStr := r.URL.Query().Get("id")
			if idStr == "" {
				http.Error(w, "Missing id parameter", http.StatusBadRequest)
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid id parameter", http.StatusBadRequest)
				return
			}

			var input struct {
				Done bool `json:"done"`
			}

			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			task, ok := store.UpdateTaskDone(id, input.Done)
			if !ok {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
