package main

import (
	"encoding/json"
	"flag"
	"github.com/DemoLiang/wss/golib"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	InitialGameMap GameMap
	addr           = flag.String("addr", ":8888", "http service address")
)

func main() {
	cfgFile := flag.String("conf", "/Users/arts/workspace/go/src/github.com/DemoLiang/wss/etc/map.json", "config file path")
	flag.Parse()
	if cfgFile == nil || *cfgFile == "" {
		flag.Usage()
		return
	}

	cfg, err := os.Open(*cfgFile)
	if err != nil {
		golib.Log("err:%v", err)
		return
	}
	defer cfg.Close()

	bs, err := ioutil.ReadAll(cfg)
	if err != nil {
		golib.Log("err:%v", err)
		return
	}
	err = json.Unmarshal(bs, &InitialGameMap)
	if err != nil {
		golib.Log("err:%v", err)
		return
	}
	for idx, data := range InitialGameMap.Map {
		golib.Log("%v %v\n", idx, data)
	}
	//golib.Log("InitialGameMap:%v\n", InitialGameMap)

	//启动注册函数
	go h.run()

	//启动WS处理函数
	http.HandleFunc("/ws", WsHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
