package main

import (
	"fmt"
	"net/http"
	"flag"
	"net/url"
)

var cmd = flag.String("command", "youtube-dl", "The command for the youtube-dl executable including the full path if applicable")

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := url.QueryEscape(r.FormValue("url"))
	fmt.Fprintf(w, "Hi there, I love %s!", url)
}
