package middleware

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"under_construction/app"
)

func TestNoGzip(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	rr := httptest.NewRecorder()

	handler := GzipWrapper(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get(app.HeaderKeyContentEncoding) != "" {
		t.Fatalf(`expected Content-Encoding: "" got %s`, rr.Header().Get(app.HeaderKeyContentEncoding))
	}

	if rr.Body.String() != "test" {
		t.Fatalf(`expected "test" go "%s"`, rr.Body.String())
	}

	if testing.Verbose() {
		b, _ := httputil.DumpResponse(rr.Result(), true)
		t.Log("\n" + string(b))
	}
}

func TestGzip(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(app.HeaderKeyAcceptEncoding, "gzip, deflate")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(app.HeaderKeyContentLength, "4")
		w.Header().Set(app.HeaderKeyContentType, "text/test")
		w.Write([]byte("test"))
	})

	rr := httptest.NewRecorder()

	handler := GzipWrapper(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Header().Get(app.HeaderKeyContentEncoding) != "gzip" {
		t.Fatalf("expected Content-Encoding: gzip got %s", rr.Header().Get(app.HeaderKeyContentEncoding))
	}
	if rr.Header().Get(app.HeaderKeyContentLength) != "" {
		t.Fatalf(`expected Content-Length: "" got %s`, rr.Header().Get(app.HeaderKeyContentLength))
	}
	if rr.Header().Get(app.HeaderKeyContentType) != "text/test" {
		t.Fatalf(`expected Content-Type: "text/test" got %s`, rr.Header().Get(app.HeaderKeyContentType))
	}

	r, err := gzip.NewReader(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "test" {
		t.Fatalf(`expected "test" go "%s"`, string(body))
	}

	if testing.Verbose() {
		b, _ := httputil.DumpResponse(rr.Result(), true)
		t.Log("\n" + string(b))
	}
}

func TestNoBody(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(app.HeaderKeyAcceptEncoding, "gzip, deflate")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	rr := httptest.NewRecorder()

	handler := GzipWrapper(testHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if rr.Header().Get(app.HeaderKeyContentEncoding) != "" {
		t.Fatalf(`expected Content-Encoding: "" got %s`, rr.Header().Get(app.HeaderKeyContentEncoding))
	}

	if rr.Body.Len() > 0 {
		t.Logf("%q", rr.Body.String())
		t.Fatalf("no body expected for %d bytes", rr.Body.Len())
	}

	if testing.Verbose() {
		b, _ := httputil.DumpResponse(rr.Result(), true)
		t.Log("\n" + string(b))
	}
}

func BenchmarkGzip(b *testing.B) {
	body := []byte("ttttttttttesttesttesttesttesttesttesttesttesttesttesttesttest")

	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Set(app.HeaderKeyAcceptEncoding, "gzip, deflate")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write(body)
			})

			rr := httptest.NewRecorder()

			handler := GzipWrapper(testHandler)
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				b.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			if rr.Body.Len() != 50 {
				b.Fatalf("expected 50 bytes, got %d bytes", rr.Body.Len())
			}
		}
	})
}
