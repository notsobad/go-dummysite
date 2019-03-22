package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strconv"
	"time"
)

type dynamicResp struct {
	Path      string `json:"path"`
	Query     string `json:"query"`
	URI       string `json:"uri"`
	Body      string `json:"body"`
	Arguments string `json:"arguments"`
	Headers   string `json:"headers"`
}

var urls = []string{
	"/static/abc.js",
	"/static/abc/xyz.css",
	"/static/abc/xyz/uvw.txt",
	"/static/abc.html",
	"/static/abc.jpg",
	"/dynamic/abc.php",
	"/dynamic/abc.asp",
	"/code/200",
	"/code/400",
	"/code/404",
	"/code/502",
	"/size/11k.zip",
	"/size/1k.bin",
	"/slow/3",
	"/slow/4-10",
	"/redirect/301?url=http://www.notsobad.vip",
	"/redirect/302?url=http://www.notsobad.vip",
	"/redirect/js?url=http://www.notsobad.vip",
}

func getNodeID() string {
	hostname, _ := os.Hostname()
	h := md5.New()
	h.Write([]byte(hostname))
	return hex.EncodeToString(h.Sum(nil))[:7]
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	const tpl = `
        <h1>YNM3K Test site</h1>
        <h2>Request header</h2>
        <pre>{{.Headers}}
        </pre>
        <h2>Links</h2>
		<ul>
			{{range .Urls}}
            <li><a href="{{.}}">{{.}}</a></li>
            {{end}}
        </ul>
        <footer>
            <hr/>SERVER-ID: {{.NodeID}}, Powered by YNM3K <a href="https://github.com/notsobad/ynm3k">Fork me</a> on Github
        </footer>
	`
	headers, _ := httputil.DumpRequest(r, true)
	data := struct {
		Urls    []string
		NodeID  string
		Headers string
	}{
		Urls:    urls,
		NodeID:  getNodeID(),
		Headers: string(headers),
	}

	t, _ := template.New("webpage").Parse(tpl)
	_ = t.Execute(w, data)
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

	fmt.Fprintf(w, "hello :-)<pre>%s</pre><hr>%s", respJSON, "happen")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	fmt.Fprintf(w, "Start at: %s<br/>", now.Format(time.RFC3339))
	time.Sleep(10 * time.Second)
	now = time.Now()
	fmt.Fprintf(w, "End at: %s", now.Format(time.RFC3339))
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/static/{filename:.*}", staticHandler)
	r.HandleFunc("/code/{code:[1-5][0-9][0-9]}", codeHandler)
	r.HandleFunc("/dynamic/{filename:.*}", dynamicHandler)
	r.HandleFunc("/slow/{[0-9]+}", slowHandler)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
