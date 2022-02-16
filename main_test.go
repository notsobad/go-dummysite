package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestIndexHandler(t *testing.T) {
	// Create handler
	handler := indexHandler

	// Create and record request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Got HTTP status code %d, expect 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "SERVER-ID") {
		t.Errorf("Got HTTP status code %d, expect 200", w.Code)
	}
}
func TestStaticHandler(t *testing.T) {
	// Create and record request
	req := httptest.NewRequest("GET", "/static/1.jpg", nil)
	w := httptest.NewRecorder()
	r := appRouter()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Got HTTP status code %d, expect 200", w.Code)
	}
	fmt.Println("body", w.Body.String())
	if w.Header().Get("Content-Type") != "image/jpeg" {
		t.Errorf("Expect image/jpeg")
	}
}
func TestCodeHandler(t *testing.T) {
	codes := []int{200, 301, 302, 400, 403, 404, 500, 502, 504}
	for _, code := range codes {
		req := httptest.NewRequest("GET", fmt.Sprintf("/code/%d", code), nil)
		w := httptest.NewRecorder()
		r := appRouter()
		r.ServeHTTP(w, req)

		if w.Code != code {
			t.Errorf("Got HTTP status code %d, expect %d", w.Code, code)
		}
	}
}
func TestDynamicHandler(t *testing.T) {
	paths := []string{"a.php", "a.jsp", "a.asp", "a.cgi"}
	for _, path := range paths {
		uri := fmt.Sprintf("/dynamic/%s", path)
		req := httptest.NewRequest("GET", uri, nil)
		w := httptest.NewRecorder()
		r := appRouter()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Got HTTP status code %d, expect %d", w.Code, 200)
		}

		if !strings.Contains(w.Body.String(), uri) {
			t.Errorf("The body should contain %s", uri)
		}
	}
}
func TestRedirectHandler(t *testing.T) {

	for _, code := range []int{301, 302} {
		rurl := "http://www.notsobad.vip/"
		uri := fmt.Sprintf("/redirect/%d?url=%s", code, rurl)
		req := httptest.NewRequest("GET", uri, nil)
		w := httptest.NewRecorder()
		r := appRouter()
		r.ServeHTTP(w, req)

		if w.Code != code {
			t.Errorf("Got HTTP status code %d, expect %d", w.Code, code)
		}

		if w.Header().Get("Location") != rurl {
			t.Errorf("Location not match, expect %s", rurl)
		}
		//fmt.Println(w.Header())
	}
}

func TestRedirectHandler2(t *testing.T) {

	for _, method := range []string{"js", "meta"} {
		rurl := "http://www.notsobad.vip/"
		uri := fmt.Sprintf("/redirect/%s?url=%s", method, rurl)
		req := httptest.NewRequest("GET", uri, nil)
		w := httptest.NewRecorder()
		r := appRouter()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Got HTTP status code %d, expect %d", w.Code, 200)
		}

		if !strings.Contains(w.Body.String(), rurl) {
			t.Errorf("Body not match, expect %s in body", rurl)
		}
	}
}

func TestSlowHandler(t *testing.T) {
	for _, sec := range []int{3, 5} {
		start := time.Now().UnixMilli()
		uri := fmt.Sprintf("/slow/%d", sec)
		req := httptest.NewRequest("GET", uri, nil)
		w := httptest.NewRecorder()
		r := appRouter()
		r.ServeHTTP(w, req)
		end := time.Now().UnixMilli()

		timeSpend := end - start
		fmt.Println("time ", timeSpend)

		if timeSpend-int64(sec)*1000 >= 1000 {
			t.Errorf("Time spend %d, expect %d", timeSpend, sec*1000)
		}
		if w.Code != 200 {
			t.Errorf("Got HTTP status code %d, expect %d", w.Code, 200)
		}

	}
}
