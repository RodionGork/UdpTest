package web

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "github.com/rodiongork/udptest/pkg/collect"
)

func RunWebServer(port int, collector collect.EventCollector) error {
    fmt.Println("Starting web-ui at port", port, "user Ctrl-C to terminate")
    http.HandleFunc("/static/", handleStatic);
    http.HandleFunc("/api/latest", func(w http.ResponseWriter, req *http.Request) {handleApiLatest(w, req, collector)});
    http.HandleFunc("/api/count", func(w http.ResponseWriter, req *http.Request) {handleApiCount(w, req, collector)});
    http.HandleFunc("/api/get/", func(w http.ResponseWriter, req *http.Request) {handleApiGet(w, req, collector)});
    http.Handle("/", http.RedirectHandler("/static/main.html", http.StatusFound))
    
    return http.ListenAndServe(":" + strconv.Itoa(port), nil)
}

func handleApiLatest(w http.ResponseWriter, req *http.Request, ec collect.EventCollector) {
    res := ec.Latest(100, req.URL.Query().Get("type"), req.URL.Query().Get("level"))
    bytes, _ := json.Marshal(res)
    w.Header().Set("Content-Type", "application/json")
    w.Write(bytes)
}

func handleApiCount(w http.ResponseWriter, req *http.Request, ec collect.EventCollector) {
    res := map[string]int{"total": ec.TotalCount()}
    bytes, _ := json.Marshal(res)
    w.Header().Set("Content-Type", "application/json")
    w.Write(bytes)
}

func handleApiGet(w http.ResponseWriter, req *http.Request, ec collect.EventCollector) {
    uuid := strings.TrimPrefix(req.URL.Path, "/api/get/")
    var bytes []byte
    event, err := ec.Retrieve(uuid)
    if err != nil {
        bytes = []byte(err.Error())
    } else {
        bytes, _ = json.Marshal(event)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(bytes)
}

func handleStatic(w http.ResponseWriter, req *http.Request) {
    path := strings.TrimPrefix(req.URL.Path, "/static/")
    http.ServeFile(w, req, "static/" + path)
}
