package main

import (
    "encoding/json"
    "log"
    "net/http/httputil"
    "net/http"
    "net/url"
    "regexp"
)


type errorResponse struct {
    Code int `json: code`
    Message string `json: message`
}


func main() {
    var acctPattern = regexp.MustCompile("account")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        var proxy *httputil.ReverseProxy

        switch {
            case acctPattern.MatchString(r.RequestURI):
                proxy = reverseProxy("http://account")
	    default:
		log.Printf("Route not found: %s", r.RequestURI)

		// Set Headers
		w.WriteHeader(http.StatusNotFound)
                w.Header().Set("Content-Type", "application/json")

		// Configure & Write Response
		msg := errorResponse{404, "Route not found"}
                js, _ := json.Marshal(msg)
		w.Write(js)
		return
        }
        proxy.ServeHTTP(w, r)
        log.Printf("%s %s %s%s", r.Header.Get("X-Forwarded-For"), r.Method, r.Host, r.RequestURI)
    })
    log.Fatal(http.ListenAndServe(":80", nil))
}

func reverseProxy(target string) *httputil.ReverseProxy {
    url, _ := url.Parse(target)
    return httputil.NewSingleHostReverseProxy(url)
}
