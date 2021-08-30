package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type hash struct {
	Password string `json:"password"`
	SHA256   string `json:"sha256"`
}

func main() {
	http.HandleFunc("/password", handler)

	log.Println("Server is listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	x := r.FormValue("sha256")

	f, err := os.Open("password-sha256.txt")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		h := hash{}
		err := json.Unmarshal(s.Bytes(), &h)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if h.SHA256 == x {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"password": "%s"}`, h.Password)
			return
		}
	}

	if s.Err() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
