package handlers

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	err2 "under_construction/app/apperrors"
)

func TestHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "../../../assets/css/error_500.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeStatic)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	data, err := ioutil.ReadFile("./../../assets/css/error_500.css")
	if err != nil {
		t.Fatal(err)
	}
	expected := string(data)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got \n %v \nwant\n %v",
			rr.Body.String(), expected)
	}
}

func TestHandlerError(t *testing.T) {
	defer func() {
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

	}()

	req, err := http.NewRequest("GET", "../../../assets/css/error_500.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := NewRecorderTest()
	handler := http.HandlerFunc(ServeStatic)

	handler.ServeHTTP(rr, req)

}

func TestHandlerNotFoundError(t *testing.T) {
	defer func() {
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

	}()

	req, err := http.NewRequest("GET", "/assets/css/montserrat2.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeStatic)

	handler.ServeHTTP(rr, req)

}

func TestContentType(t *testing.T) {

	expected := []string{
		"text/css",
		"text/html",
		"font/woff2",
		"application/javascript",
		"image/png", "image/jpg",
		"image/svg+xml",
		"text/plain"}

	income := []string{
		"assets/css/error_500.css",
		"assets/html/example.html",
		"assets/woff2/QldKNThLqRwH-OJ1UHjlKGlW5qhExfHwNJU.woff2",
		"assets/js/super.js",
		"assets/image/flame.png",
		"assets/image/under_632.jpg",
		"assets/image/under_s.svg",
		"hdhfy.uye"}

	for i := 0; i < len(income); i++ {
		result := GetContentType(income[i])
		if result != expected[i] {
			t.Errorf("getContentType returned unexpected result: got \n %v \nwant\n %v",
				result, expected[i])
		}
	}

}

// NewRecorder returns an initialized ResponseRecorder.
func NewRecorderTest() *ResponseRecorderTest {
	return &ResponseRecorderTest{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
	}
}

// ResponseRecorder is an implementation of http.ResponseWriter that
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
	return 0, err2.NewNotFoundError()
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
