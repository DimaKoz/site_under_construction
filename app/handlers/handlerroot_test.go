package handlers

import (
	"errors"
	log "github.com/google/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"under_construction/app"
	err2 "under_construction/app/apperrors"
)

var log1 *log.Logger

var lf *os.File

func initLogger() {
	var errLog error
	lf, errLog = os.OpenFile(app.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if errLog != nil {
		log.Fatalf("Failed to open log file: %v", errLog)
	}
	_ = log.Init("LoggerExample", true, false, lf)
}

func TestRootHandler(t *testing.T) {

	initLogger()
	defer lf.Close()

	//mocking GetBytes
	_ = app.AddKeyAndPath(app.HtmlUnderConstruction, "./../../html/under_construction.html")

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RootHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	data, err := app.GetBytes(app.HtmlUnderConstruction)
	if err != nil {
		t.Fatal(err)
	}

	expected := string(*data)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got \n %v \nwant\n %v",
			rr.Body.String(), expected)
	}

}

func recoveringExpectNotFoundPanic(t *testing.T) {
	_ = lf.Close()
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	} else {
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

}

func TestRootHandlerNotFound2(t *testing.T) {
	initLogger()

	defer func() {
		_ = lf.Close()
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_ = app.AddKeyAndPath(app.HtmlUnderConstruction, "./../../html/under_construction.html")
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := NewRecorderTest()
	handler := http.HandlerFunc(RootHandler)

	handler.ServeHTTP(rr, req)

}
