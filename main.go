package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Greeter interface {
	greet(name string) (string, error)
}

type GreetCounter struct {
	count int
	mu    sync.Mutex
}

func (s *GreetCounter) greet(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name cannnot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.count = s.count + 1

	fmt.Println(s.count)
	return fmt.Sprintf("welcome back %s, greet number %d", name, s.count), nil
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed);
		return;
	}

	w.WriteHeader(http.StatusOK);
	fmt.Fprintln(w, "healthy API")
}

func main() {
	service := &GreetCounter{}
	http.HandleFunc("/greet", greetHandler(service))
	http.HandleFunc("/health", healthHandler)

	http.ListenAndServe(":8080", nil);
	fmt.Println("server running on port 8080");
}