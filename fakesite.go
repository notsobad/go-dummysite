package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"path"
	"strconv"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hi sss %s", r.URL.Path[1:])
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

	fmt.Fprintf(w, filename)
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
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/static/{filename:.*}", staticHandler)
	r.HandleFunc("/code/{code:[1-5][0-9][0-9]}", codeHandler)
	r.HandleFunc("/dynamic/{filename:.*}", dynamicHandler)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
