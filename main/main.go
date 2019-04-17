package main

import (
	"flag"
	"github.com/lexkong/log"
	"gosignaler-cluster/handler"
	"gosignaler-cluster/rpcservice"
	"gosignaler-cluster/signalerconst"
	"gosignaler-cluster/util"
	"net"
	"net/http"
	"net/rpc"
)

var addr = flag.String("addr", signalerconst.SERVER_PORT, "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", signalerconst.RPC_FAIL)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", signalerconst.METHOD_NOT_ALLOW)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {

	flag.Parse()
	//init the log config
	util.InitLogCfg()
	hub := handler.NewHub()

	//rpc service
	go rpcsignaler(hub)
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf("URL: %s\n", r.URL.String())
		r.ParseForm()
		id := r.Form.Get("id")
		log.Infof("has connected........%s",id)
		if id != "" {
			handler.ServeWs(hub, w, r, id)
		}
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

//rpc service
func rpcsignaler(hub *handler.Hub) {

	signalerService := new(rpcservice.HandleSignalService)
	signalerService.Hub = hub

	//register service
	rpcservice.RegisterHandleSignalService(signalerService)

	listener, err := net.Listen("tcp", signalerconst.REMOTE_ADDRESS_PORT)
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		go rpc.ServeConn(conn)
	}
}
