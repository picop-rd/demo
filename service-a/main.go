package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/picop-rd/picop-go/contrib/net/http/picophttp"
	"github.com/picop-rd/picop-go/propagation"
	picopnet "github.com/picop-rd/picop-go/protocol/net"
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
		Handler:     picophttp.NewHandler(http.DefaultServeMux, propagation.EnvID{}),
		ConnContext: picophttp.ConnContext,
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	bln := picopnet.NewListener(ln)
	fmt.Println("serve http")
	log.Fatal(server.Serve(bln))
}

func handler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Transport: picophttp.NewTransport(nil, propagation.EnvID{}),
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
	w.Write([]byte(*message + "\n"))
	io.Copy(w, resp.Body)
	return
}
