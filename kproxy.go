package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
)

var configFile = flag.String("conf", "conf.json", "configuration file")
var config map[string]interface{}

func register(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s%s", r.Method, r.RemoteAddr, r.RequestURI)
		p.ServeHTTP(w, r)
	}
}

func runServer(u *url.URL, port string) {
    proxy := httputil.NewSingleHostReverseProxy(u)
    director := proxy.Director
    proxy.Director = func(req *http.Request) {
        director(req)
        req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
        req.Host = req.URL.Host
    }

    server := http.NewServeMux()
    server.HandleFunc("/", register(proxy))
    http.ListenAndServe(fmt.Sprintf(":%s", port), server)
}

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: %s [options]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

    finish := make(chan bool)

	folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(filepath.Join(folder, *configFile))
	if err != nil {
		log.Fatalln(err)
	}

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatalln(err)
	}

	for port, host := range config["routes"].(map[string]interface{}) {
		log.Printf("%s -> %s", port, host)
        u, err := url.Parse(host.(string))
        if err != nil {
            // skip invalid hosts
            log.Println(err)
            continue
        }
        go runServer (u, port)
	}
	<-finish
}
