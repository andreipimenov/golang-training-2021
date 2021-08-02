package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	switch t := v.(type) {
	case error:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", t)))
	default:
		b, err := json.Marshal(v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Internal server error"}`))
			return
		}
		w.WriteHeader(code)
		w.Write([]byte(b))
	}
}
