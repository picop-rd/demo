package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	"github.com/hiroyaonoe/bcop-go/propagation"
	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
)

var (
	port         = flag.String("port", "80", "listen port")
	childService = flag.String("child-service", "http://localhost:8080", "child service url")
	message      = flag.String("message", "", "response message")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:        fmt.Sprintf(":%s", *port),
		Handler:     bcophttp.NewHandler(http.DefaultServeMux, propagation.EnvID{}),
		ConnContext: bcophttp.ConnContext,
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	bln := bcopnet.NewListener(ln)
	fmt.Println("serve http")
	log.Fatal(server.Serve(bln))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// child serviceにリクエスト
	client := &http.Client{
		Transport: bcophttp.NewTransport(nil, propagation.EnvID{}),
	}
	req, err := http.NewRequestWithContext(r.Context(), r.Method, *childService, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	// messageを返す
	w.Write([]byte(*message + "\n"))
	// child serviceのレスポンスを返す
	io.Copy(w, resp.Body)
	return
}
