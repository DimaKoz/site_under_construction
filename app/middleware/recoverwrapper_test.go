package middleware

import (
	"bytes"
	log "github.com/google/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"under_construction/app"
	"under_construction/app/apperrors"
)

var _ *log.Logger

var lf *os.File

func initLogger() {
	var errLog error
	lf, errLog = os.OpenFile(app.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if errLog != nil {
		log.Fatalf("Failed to open log file: %v", errLog)
	}
	_ = log.Init("LoggerExample", true, false, lf)
}

func recoveringExpectPanic(t *testing.T) {
	_ = lf.Close()

	r := recover()
	if r != nil {
		t.Errorf("The code do panic")
	}
}

func TestRecoverWrap404(t *testing.T) {

	initLogger()
	defer recoveringExpectPanic(t)
	_ = app.AddKeyAndPath(app.Html404, "./../../html/error_404_page.html")
	req, err := http.NewRequest("GET", "/api/users" /*there is any path*/, nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		panic(apperrors.NewNotFoundError())
	})

	rr := httptest.NewRecorder()
	// func middleware (h http.Handler) http.Handler
	handler := RecoverWrap(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestRecoverWrap500WhenNotFoundPanic(t *testing.T) {

	initLogger()
	defer recoveringExpectPanic(t)
	app.RemoveKey(app.Html404)
	_ = app.AddKeyAndPath(app.Html404, "")
	req, err := http.NewRequest("GET", "/api/users" /*there is any path*/, nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		panic(apperrors.NewNotFoundError())

	})

	rr := httptest.NewRecorder()
	// func middleware (h http.Handler) http.Handler

	handler := RecoverWrap(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestRecoverWrap500When500(t *testing.T) {

	initLogger()
	defer recoveringExpectPanic(t)
	app.RemoveKey(app.Html500)
	_ = app.AddKeyAndPath(app.Html500, "./../../html/error_500_page.html")
	req, err := http.NewRequest("GET", "/api/users" /*there is any path*/, nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		panic(os.PathError{})

	})

	rr := NewRecorderTest()
	// func middleware (h http.Handler) http.Handler

	handler := RecoverWrap(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestLogging(t *testing.T) {
	initLogger()
	defer recoveringExpectPanic(t)
	app.RemoveKey(app.Html500)
	_ = app.AddKeyAndPath(app.Html500, "./../../html/error_500_page.html")
	req, err := http.NewRequest("GET", "/api/users" /*there is any path*/, nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		panic("os.PathError{}")

	})

	rr := httptest.NewRecorder()

	// func middleware (h http.Handler) http.Handler

	handler := RecoverWrap(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	loggingErr(nil)
}

// NewRecorder returns an initialized ResponseRecorder.
func NewRecorderTest() *ResponseRecorderTest {
	return &ResponseRecorderTest{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
	}
}

// ResponseRecorderTest is an implementation of http.ResponseWriter that
// records its mutations for later inspection in tests.
type ResponseRecorderTest struct {
	// Code is the HTTP response code set by WriteHeader.
	Code int

	// HeaderMap contains the headers explicitly set by the Handler.
	// It is an internal detail.
	HeaderMap http.Header

	// Body is the buffer to which the Handler's Write calls are sent.
	// If nil, the Writes are silently discarded.
	Body *bytes.Buffer

	// Flushed is whether the Handler called Flush.
	Flushed bool

	result      *http.Response // cache of Result's return value
	snapHeader  http.Header    // snapshot of HeaderMap at first Write
	wroteHeader bool
}

func (rw *ResponseRecorderTest) Write(buf []byte) (int, error) {
	return 0, apperrors.NewNotFoundError()
}

func (rw *ResponseRecorderTest) Header() http.Header {
	m := rw.HeaderMap
	if m == nil {
		m = make(http.Header)
		rw.HeaderMap = m
	}
	return m
}

func (rw *ResponseRecorderTest) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.Code = code
	rw.wroteHeader = true
	if rw.HeaderMap == nil {
		rw.HeaderMap = make(http.Header)
	}
	rw.snapHeader = rw.HeaderMap.Clone()
}
