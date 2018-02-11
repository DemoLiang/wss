package main

import (
	"flag"
	"go/build"
	"log"
	"net/http"
	"text/template"
	"fmt"
	"path/filepath"
	"github.com/DemoLiang/wss/golib"
)

var (
	addr      = flag.String("addr", ":7777", "http service address")
	assets    = flag.String("assets", defaultAssetPath(), "path to assets")
	homeTempl *template.Template
)

func defaultAssetPath() string {
	p, err := build.Default.Import("wss", "", build.FindOnly)
	if err != nil {
		return "."
	}
	golib.Log("%v\n",p.Dir)
	return p.Dir
}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func benchmarkHandler(c http.ResponseWriter, req *http.Request) {
	h.Broadcast <- []byte("test message")
}

func main() {
	flag.Parse()
	fmt.Printf("homeTempl:%v\n",filepath.Join(*assets, "index.html"))
	go h.run()
	http.HandleFunc("/ws", WsHandler)
	http.HandleFunc("/benchmark", benchmarkHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
