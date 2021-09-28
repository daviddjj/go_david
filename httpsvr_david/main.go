package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	flag.Set("v", "1")
	fmt.Println("Starting http server...")
	glog.Info("Starting http server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	goVersion := os.Getenv("VERSION")
	ipAddr := r.RemoteAddr
	ipAddr_s := strings.Index(ipAddr, ":")
	ipAddr = ipAddr[0:ipAddr_s]
	//将接到的header写入到会包的header
	for k, v := range r.Header {
		var vvv string
		for _, vv := range v {
			vvv += vv
		}
		w.Header().Set(k, vvv)
	}
	//读取本地变量VERSION，写入header中
	w.Header().Set("VERSION", goVersion)
	w.WriteHeader(200)
	io.WriteString(w, "hello, world!\n")
	fmt.Println("Remote IPAddress is", ipAddr)
	glog.Info("Remote IPAddress is", ipAddr)
}
