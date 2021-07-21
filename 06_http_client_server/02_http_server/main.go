package main

import (
	"fmt"
	"log"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ping")
}

func main() {
	srv := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/ping", ping)

	log.Println("Server is listening on :8080")
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// http.HandleFunc("/long", bigSyncJob)

	// userHandler := &UserHandler{}
	// http.Handle("/user", userHandler)

	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	// defer signal.Stop(shutdown)

	// go func() {
	// 	log.Println("Server is listening on :8080")
	// 	err := srv.ListenAndServe()
	// 	if err != nil && err != http.ErrServerClosed {
	// 		log.Fatal(err)
	// 	}
	// }()

	// <-shutdown

	// log.Println("Shutdown signal received")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	// defer func() {
	// 	cancel()
	// }()

	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Server stopped gracefully")

}

// type UserHandler struct {
// 	data sync.Map
// }

// type User struct {
// 	Name string
// 	Age  int
// }

// type Error struct {
// 	Error string
// }

// func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodPost:
// 		u := User{}
// 		err := json.NewDecoder(r.Body).Decode(&u)
// 		if err != nil {
// 			writeResponse(w, http.StatusBadRequest, Error{err.Error()})
//          return
// 		}
// 		h.data.Store(u.Name, u)
// 		writeResponse(w, http.StatusOK, u)
// 	case http.MethodGet:
// 		userName := r.FormValue("name")
// 		u, ok := h.data.Load(userName)
// 		if !ok {
// 			writeResponse(w, http.StatusNotFound, Error{"User not found"})
// 			return
// 		}
// 		writeResponse(w, http.StatusOK, u)
// 	default:
// 		writeResponse(w, http.StatusMethodNotAllowed, Error{"Method not allowed"})
// 	}
// }

// func writeResponse(w http.ResponseWriter, code int, v interface{}) {
// 	b, _ := json.Marshal(v)
// 	w.WriteHeader(code)
// 	w.Write([]byte(b))
// }

// func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.RequestURI, r.RemoteAddr, r.UserAgent())
// 		next(w, r)
// 	}
// }

// func bigSyncJob(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	log.Println("Long work started")

// 	select {
// 	case <-time.After(10 * time.Second):
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Work is done"))
// 	case <-ctx.Done():
// 		log.Printf("Interrupted with err: %v", ctx.Err())
// 	}
// }
