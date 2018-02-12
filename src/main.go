package main

import (
	"flag"
	"go/build"
	"log"
	"net/http"
	"text/template"
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
	golib.Log("homeTempl:%v\n",filepath.Join(*assets, "index.html"))

	//启动注册函数
	go h.run()

	//启动WS处理函数
	http.HandleFunc("/ws", WsHandler)
	http.HandleFunc("/benchmark", benchmarkHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
