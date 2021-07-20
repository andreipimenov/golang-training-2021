package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ping")
}

func main() {
	http.HandleFunc("/ping", ping)

	userHandler := &UserHandler{}
	http.Handle("/user", userHandler)

	log.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		} else {
			log.Println("Server closed gracefully")
		}
	}
}

type UserHandler struct {
	data sync.Map
}

type User struct {
	Name string
	Age  int
}

type Error struct {
	Error string
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		u := User{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		}
		h.data.Store(u.Name, u)
		writeResponse(w, http.StatusOK, u)
	case http.MethodGet:
		userName := r.FormValue("name")
		u, ok := h.data.Load(userName)
		if !ok {
			writeResponse(w, http.StatusNotFound, Error{"User not found"})
			return
		}
		writeResponse(w, http.StatusOK, u)
	default:
		writeResponse(w, http.StatusMethodNotAllowed, Error{"Method not allowed"})
	}
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	w.Write([]byte(b))
}
