package handlers

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"under_construction/app"
	err2 "under_construction/app/apperrors"
)

func TestFaviconHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeFavicon)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(rr.Body.String())))
	base64.StdEncoding.Encode(base64Text, []byte(rr.Body.String()))

	expected := app.FaviconData
	if string(base64Text) != expected {
		t.Errorf("handler returned unexpected body: got \n %v \nwant\n %v",
			rr.Body.String(), expected)
	}
}

func TestFaviconHandlerError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
			var err error
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("Unknown error")
			}
			_, ok := err.(*err2.NotFoundError)
			if !ok {
				t.Errorf("Unxpected type of error")
			}
		}

	}()

	req, err := http.NewRequest("GET", "/assets/css/error_500.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := NewRecorderTest()
	handler := http.HandlerFunc(ServeFavicon)

	handler.ServeHTTP(rr, req)

}
