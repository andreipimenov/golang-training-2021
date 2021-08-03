package handler

import (
	"encoding/json"
	"errors"
	"log"
	"mime"
	"net/http"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type errResponse struct {
	Error string `json:"error"`
}

func writeResponse(w http.ResponseWriter, v interface{}) {
	code := http.StatusOK

	setISE := func(err error) {
		log.Println(err)
		code = http.StatusInternalServerError
		v = errResponse{"Internal server error"}
	}

	if err, ok := v.(error); ok {
		if badRequest := (model.BadRequest{}); errors.As(err, &badRequest) {
			code = http.StatusBadRequest
			v = errResponse{badRequest.Msg}
		} else {
			setISE(err)
		}
	}

	body, err := json.Marshal(v)
	if err != nil {
		setISE(err)
	}
	w.Header().Add("Content-Type", mime.TypeByExtension(".json"))
	w.WriteHeader(code)
	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
}
