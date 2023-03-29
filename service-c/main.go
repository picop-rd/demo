package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hiroyaonoe/bcop-go/contrib/github.com/go-sql-driver/mysql/bcopmysql"
	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	"github.com/hiroyaonoe/bcop-go/propagation"
	bcopnet "github.com/hiroyaonoe/bcop-go/protocol/net"
)

var (
	port         = flag.String("port", "80", "listen port")
	mysqlService = flag.String("mysql-service", "", "mysql service dsn")
	message      = flag.String("message", "", "response message")
)

var db *sql.DB

func main() {
	var err error
	flag.Parse()

	bcopmysql.RegisterDialContext("tcp", propagation.EnvID{})
	db, err = sql.Open("mysql", *mysqlService)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxIdleConns(0)
	for {
		fmt.Println("connecting db")
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println("connected db")

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
	var header int
	var data string
	// mysql serviceにリクエスト
	switch r.Method {
	case http.MethodGet:
		header, data = get(r)
	case http.MethodPost:
		header, data = post(r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(header)
	if header == http.StatusOK {
		// messageを返す
		w.Write([]byte(*message + "\n"))
		w.Write([]byte(data))
	}
}

func get(r *http.Request) (header int, data string) {
	rows, err := db.QueryContext(r.Context(), "SELECT id, content FROM data")
	if err != nil {
		return http.StatusInternalServerError, ""
	}
	for rows.Next() {
		id := 0
		content := ""
		if err = rows.Scan(&id, &content); err != nil {
			return http.StatusInternalServerError, ""
		}
		data += fmt.Sprintf("data{ id: %d, content: %s }\n", id, content)
	}
	return http.StatusOK, data
}

func post(r *http.Request) (header int, data string) {
	content := make([]byte, 256)
	_, err := r.Body.Read(content)
	if err != nil && err != io.EOF {
		return http.StatusBadRequest, ""
	}
	_, err = db.ExecContext(r.Context(), "INSERT INTO data(content) VALUES (?)", content)
	if err != nil {
		return http.StatusInternalServerError, ""
	}
	return http.StatusOK, ""
}
