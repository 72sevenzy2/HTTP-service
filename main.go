package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// json helpers
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string) {
	resp := map[string]string{
		"error": message,
	}
	JSON(w, status, resp)
}

// config structs
type GreetResponse struct {
	Name string `json:"name"`

	Count int `json:"count"`
}

type Request struct {
	Name string
}

type Greeter interface {
	greet(name string) (string, error)
}

type GreetCounter struct {
	count int
	mu    sync.Mutex
}

func (s *GreetCounter) greet(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name cannot be empty")
	}

	s.mu.Lock()
	s.count = s.count + 1
	count := s.count
	defer s.mu.Unlock()

	return fmt.Sprintf("welcome back %s, this is the %d request", name, count), nil
}

func greetHandler(g Greeter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")

		msg, err := g.greet(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, msg)
	}
}

func healthChecker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "fully working API")
}

func main() {
	service := &GreetCounter{}

	http.HandleFunc("/greet", greetHandler(service))
	http.HandleFunc("/health", healthChecker)

	fmt.Println("API working on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
