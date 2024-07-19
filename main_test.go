package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
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

// generate test for chunkHandler
func TestChunkHandler(t *testing.T) {
	// 创建一个路由，模拟实际的路由行为
	r := mux.NewRouter()
	r.HandleFunc("/chunk/{count}", chunkHandler)

	// 创建一个测试服务器
	server := httptest.NewServer(r)
	defer server.Close()

	expectedCount := 3
	url := fmt.Sprintf("%s/chunk/%d", server.URL, expectedCount)
	// 发送请求到测试服务器
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 检查Transfer-Encoding是否为chunked
	if resp.TransferEncoding == nil || resp.TransferEncoding[0] != "chunked" {
		t.Errorf("Expected chunked transfer encoding, got %v", resp.TransferEncoding)
	}

	// 读取响应体，确保它是分多次接收的
	scanner := bufio.NewScanner(resp.Body)
	lineCount := 0

	last_time := time.Now()
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Time:") && strings.Contains(line, "Msg:") {
			lineCount++
			// 检查每个分块之间的时间间隔是否大于1秒
			if lineCount > 1 {
				current_time := time.Now()
				if current_time.Sub(last_time) < 1*time.Second {
					t.Errorf("Chunk interval is less than 1 second")
				}
				last_time = current_time
			}
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading stream: %v", err)
	}
	// 确保我们接收到了预期数量的分块
	if lineCount != expectedCount {
		t.Errorf("Expected %d chunks, got %d", expectedCount, lineCount)
	}
}
