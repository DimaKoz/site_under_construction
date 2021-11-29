package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"under_construction/app"
	"under_construction/app/apperrors"
	"under_construction/app/etagging"
)

func TestCheckCache(t *testing.T) {
	app.RemoveKey(app.HtmlUnderConstruction)
	_ = app.AddKeyAndPath(app.HtmlUnderConstruction, "./../../html/under_construction.html")
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	rr := httptest.NewRecorder()

	handler := CheckCache(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := app.HeaderValueCacheControl
	if rr.Header().Get(app.HeaderKeyCacheControl) != expected {
		t.Errorf("handler returned unexpected body: got \n %v \nwant\n %v",
			rr.Body.String(), expected)
	}

}

func TestCheckCacheStatic(t *testing.T) {
	app.RemoveKey("assets/css/error_500.css")
	_ = app.AddKeyAndPath("assets/css/error_500.css", "./../../assets/css/error_500.css")
	req, err := http.NewRequest("GET", "/assets/css/error_500.css", nil)
	if err != nil {
		t.Fatal(err)
	}
	dataEtagCalculate, err := app.GetBytes("assets/css/error_500.css")
	if err != nil {
		t.Fatal(err)
	}

	etag := etagging.Generate(string(*dataEtagCalculate), true)
	req.Header.Set(app.HeaderKeyIfNoneMatch, etag)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
	rr := httptest.NewRecorder()

	handler := CheckCache(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotModified {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotModified)
	}

}

func TestCheckCacheFavicon(t *testing.T) {

	req, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	etag := etagging.Generate(app.FaviconData, true)
	req.Header.Set(app.HeaderKeyIfNoneMatch, etag)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
	rr := httptest.NewRecorder()

	handler := CheckCache(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotModified {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotModified)
	}

}

func TestCheckCachePanic(t *testing.T) {

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
			_, ok := err.(*apperrors.NotFoundError)
			if !ok {
				t.Errorf("Unxpected type of error")
			}
		}

	}()

	app.RemoveKey(app.HtmlUnderConstruction)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	rr := httptest.NewRecorder()

	handler := CheckCache(testHandler)
	handler.ServeHTTP(rr, req)
}
