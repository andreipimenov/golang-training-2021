package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error":"Internal server error"}`))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
		return
	}
}
