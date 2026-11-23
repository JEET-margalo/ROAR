package spa

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetectSPAFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><div id=app></div></html>"))
	}))
	defer srv.Close()
	p := Detect(srv.URL + "/#/login")
	if !p.SPAFallback || p.Responses200 != 3 {
		t.Fatalf("profile=%+v; want SPA fallback with three successful probes", p)
	}
}

func TestDetectNonFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/roar-nonexistent-a" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("baseline"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	p := Detect(srv.URL)
	if p.SPAFallback {
		t.Fatalf("profile=%+v; did not expect fallback", p)
	}
}
