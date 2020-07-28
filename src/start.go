package main

import (
	"util/logger"
	"net/http"
	"net"
	"plugin/docker"
	"os"
)

func main() {
	err := runHttpServer()
	if err != nil {
		os.Exit(-1)
	}
}

func runHttpServer() (err error) {
	mux := http.NewServeMux()
	mux = RegisterRoutes(mux)
	server := &http.Server{Handler: mux}
	l, err := net.Listen("tcp4", ":8080")
	if err != nil {
		logger.Error.Println(err)
		return
	}
	logger.Info.Printf("starting http server ...")
	err = server.Serve(l)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	return
}

func RegisterRoutes(mux *http.ServeMux) (*http.ServeMux) {
	mux.HandleFunc("/", docker.Index)
	mux.HandleFunc("/container/terminal", docker.Terminal)
	mux.Handle("/static/", http.FileServer(http.Dir("./src")))
	mux.HandleFunc("/nodes/containers/shell/create", docker.CreateContainer)
	mux.HandleFunc("/nodes/containers/shell/ws", docker.ShellContainer)
	mux.HandleFunc("/nodes/containers/shell/resize", docker.ResizeContainer)
	return mux
}