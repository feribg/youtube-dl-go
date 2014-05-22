package main

import (
	"encoding/json"
	"flag"
	log "github.com/cihub/seelog"
	"net/http"
	"os/exec"
	"strings"
	"net/url"
)

var cmd = flag.String("command", "youtube-dl", "The command for the youtube-dl executable including the full path if applicable")

type Src struct {
	Title  string `json: "title"`
	Url    string `json: "url"`
	Status bool   `json:"status"`
}

func main() {
	defer log.Flush()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	rc := make(chan Src)
	go process(url, rc)
	res := <-rc
	if res.Status == false {
		http.Error(w, "There was an error convering the video", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "There was an error convering the video", http.StatusBadRequest)
		return
	}
	w.Write(json)
}

func process(srcUrl string, rc chan Src) {
	//remove double quotes from the url if present
	srcUrl = strings.Replace(srcUrl, "\"", "", -1)
	//wrap the url in quotes
	out, err := exec.Command(*cmd, "--skip-download", "--get-url", "--get-title", srcUrl).CombinedOutput()
	var resp Src
	resp.Status = false
	if err != nil {
		log.Critical(err)
		log.Debugf("URL: %s", srcUrl)
		log.Debug(string(out))
		rc <- resp
	}
	parts := strings.Split(string(out), "\n")
	resp.Title = parts[0]
	resp.Url = url.QueryEscape(parts[1])
	resp.Status = true
	rc <- resp
}
