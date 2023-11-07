package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var MyLog *log.Logger

func LogHTTPRequest(req *http.Request) {
	MyLog.Println("Method:", req.Method)
	MyLog.Println("Body:", req.Body)
	MyLog.Println("URL.Path:", req.URL.Path)
	MyLog.Println("URL:", req.URL)
	MyLog.Println("ContentLength:", req.ContentLength)
	MyLog.Println("TransferEncoding:", req.TransferEncoding)
	MyLog.Println("RequestURI:", req.RequestURI)
	MyLog.Printf("req:%+v\n", req)
}

func Println(v ...any) {
	fmt.Println(v...)
}

func Run() {
	f, err := os.OpenFile("./server.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	MyLog = log.New(f, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	//Log = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	MyLog.Println("Log started")
}
