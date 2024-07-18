package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday/v2"
)

//go:embed README.md
var readmeFS embed.FS

type dynamicResp struct {
	Path      string `json:"path"`
	Query     string `json:"query"`
	URI       string `json:"uri"`
	Body      string `json:"body"`
	Arguments string `json:"arguments"`
	Headers   string `json:"headers"`
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	now := time.Now()
	cacheTime := 95270
	expired := now.Add(time.Second * time.Duration(cacheTime))

	contentType := mime.TypeByExtension(path.Ext(filename))

	w.Header().Set("Last-Modified", now.Format(time.RFC1123))
	w.Header().Set("Expires", expired.Format(time.RFC1123))
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", cacheTime))

	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	fmt.Fprint(w, filename)
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code, err := strconv.Atoi(vars["code"])
	if err != nil {
		fmt.Fprintf(w, "wrong status code %s", vars["code"])
		return
	}
	now := time.Now()
	w.WriteHeader(code)
	fmt.Fprintf(w, "<h1>Http %s</h1> <hr/>Generated at %s", vars["code"], now.Format(time.RFC3339))
}

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	headers, _ := httputil.DumpRequest(r, true)

	resp := &dynamicResp{
		Path:    r.URL.Path,
		Query:   r.URL.RawQuery,
		URI:     r.RequestURI,
		Body:    "",
		Headers: string(headers),
	}

	respJSON, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, "hello :-)<pre>%s</pre><hr>", respJSON)
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sleepTime, _ := strconv.Atoi(vars["time"])

	now := time.Now()
	fmt.Fprintf(w, "Start at: %s\n", now.Format(time.RFC3339))
	ticker := time.NewTicker(time.Duration(sleepTime) * time.Second).C
	now = <-ticker
	fmt.Fprintf(w, "End at: %s", now.Format(time.RFC3339))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	method := vars["method"]
	url := r.FormValue("url")

	switch method {
	case "301", "302":
		code, _ := strconv.Atoi(method)
		http.Redirect(w, r, url, code)
	case "js":
		fmt.Fprintf(w, "<script>location.href=\"%s\";</script>", url)
	case "meta":
		fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"0; url=%s\" />", url)
	}
}

func sizeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	size, _ := strconv.Atoi(vars["size"])
	switch vars["measure"] {
	case "k":
		size = size * 1024
	case "m":
		size = size * 1024 * 1024
	}
	fmt.Fprint(w, strings.Repeat("o", size))
}

func chunkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	count, _ := strconv.Atoi(vars["count"])

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Transfer-Encoding", "chunked")

	//time.Sleep(time.Duration(count) * time.Second)
	for i := 0; i < count; i++ {
		randomStr := randomString(10)
		currentTime := time.Now().Format(time.RFC3339)
		fmt.Fprintf(w, "Time: %s, Msg: %s\n", currentTime, randomStr)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(1 * time.Second)
	}
}

func traceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\r\n", r.Method, r.URL, r.Proto)

	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "%v: %v\r\n", name, value)
		}
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1048576))
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "\r\n%s", string(body))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//content, err := os.ReadFile("README.md")
	content, err := fs.ReadFile(readmeFS, "README.md")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	html := blackfriday.Run(content)

	w.Header().Set("Content-Type", "text/html")

	w.Write(html)
}

func appRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/static/{filename:.*}", staticHandler)
	r.HandleFunc("/code/{code:[1-5][0-9][0-9]}", codeHandler)
	r.HandleFunc("/dynamic/{filename:.*}", dynamicHandler)
	r.HandleFunc("/slow/{time:[0-9]+}", slowHandler)
	r.HandleFunc("/redirect/{method}", redirectHandler)
	r.HandleFunc("/size/{size:[0-9]+}{measure:[k|m]?}{ext:.*}", sizeHandler)
	r.HandleFunc("/chunk/{count:[0-9]+}", chunkHandler)
	r.HandleFunc("/trace", traceHandler)
	return r
}

func main() {

	ip := flag.String("ip", "0.0.0.0", "IP to use, default 0.0.0.0")
	port := flag.Int("port", 9527, "Port to use, default 9527")
	flag.Parse()

	r := appRouter()
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	address := fmt.Sprintf("%s:%d", *ip, *port)
	fmt.Println("# Listening on", address)
	log.Fatal(http.ListenAndServe(address, loggedRouter))
}
