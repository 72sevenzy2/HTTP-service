package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/72sevenzy2/http-router/router"
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
	Name string `json:"name"`
}

type Greeter interface {
	greet(name string) (string, int, error)
}

type GreetCounter struct {
	count int
	mu    sync.Mutex
}

func (s *GreetCounter) greet(name string) (string, int, error) {
	if name == "" {
		return "", 0, fmt.Errorf("name cannot be empty")
	}

	var count int

	s.mu.Lock()
	s.count = s.count + 1
	count = s.count
	s.mu.Unlock()

	return fmt.Sprintf("welcome back %s", name), count, nil
}

func greetHandler(g Greeter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			Error(w, http.StatusMethodNotAllowed, "invalid request method")
			return
		}

		var q Request
		err := json.NewDecoder(r.Body).Decode(&q)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}

		msg, count, err := g.greet(q.Name)
		if err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}

		resp := GreetResponse{
			Name:  msg,
			Count: count,
		}

		JSON(w, http.StatusOK, resp)
	}
}

func healthChecker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request", http.StatusMethodNotAllowed)
		return
	}

	resp := map[string]string{
		"message": "fully operational API",
	}

	JSON(w, http.StatusOK, resp)
}


func main() {
	service := &GreetCounter{}

	r := router.NewRouter()

	r.Handle(http.MethodPost, "/greet", greetHandler(service), router.Logger()) // using the logger middleware which is included in my http router

	fmt.Println("API working on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
