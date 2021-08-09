package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_writeResponse(t *testing.T) {
	t.Run("FailJsonMarshal", func(t *testing.T) {
		rr := httptest.NewRecorder()
		writeResponse(rr, 0, func() {})
		if rr.Result().StatusCode != http.StatusInternalServerError {
			t.Errorf("Status code is not InternalServerError")
		}
	})
	t.Run("AsExpected", func(t *testing.T) {
		rr := httptest.NewRecorder()
		writeResponse(rr, 0, 0)
		if rr.Result().StatusCode != http.StatusOK {
			t.Errorf("Status code is not OK")
		}
	})
}
