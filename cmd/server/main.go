package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func logerHandler(res http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Method: %s\r\n", req.Method)
	body += "Header ===============\r\n"
	for k, v := range req.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "Query parameters ===============\r\n"
	if err := req.ParseForm(); err != nil {
		res.Write([]byte(err.Error()))
		return
	}
	for k, v := range req.Form {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}

	res.Write([]byte(body))
}

func main() {
	fmt.Println("Start")
	mux := http.NewServeMux()
	fmt.Println("`Started mux")
	mux.HandleFunc("/", logerHandler)

	f, err := os.OpenFile("./server.info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

	//server := &http.Server{Addr: `localhost:8080`, Handler: mux, ErrorLog: infoLog}
	server := &http.Server{Addr: `localhost:8080`, Handler: mux}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
