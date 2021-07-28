package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	data sync.Map
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "userID")

	u, ok := h.data.Load(userName)
	if !ok {
		writeResponse(w, http.StatusNotFound, Error{"User not found"})
		return
	}
	writeResponse(w, http.StatusOK, u)
}

func (h *UserHandler) Post(w http.ResponseWriter, r *http.Request) {
	u := User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}
	h.data.Store(u.Name, u)
	writeResponse(w, http.StatusOK, u)
}

func bigSyncJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("Long work started")

	select {
	case <-time.After(10 * time.Second):
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Work is done"))
	case <-ctx.Done():
		log.Printf("Interrupted with err: %v", ctx.Err())
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ping")
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	w.Write([]byte(b))
}
