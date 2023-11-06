package metricks

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var infoLog *log.Logger

func logRequest(req *http.Request) {
	infoLog.Println("Method:", req.Method)
	infoLog.Println("Body:", req.Body)
	infoLog.Println("URL.Path:", req.URL.Path)
	infoLog.Println("URL:", req.URL)
	infoLog.Println("ContentLength:", req.ContentLength)
	infoLog.Println("TransferEncoding:", req.TransferEncoding)
	infoLog.Println("RequestURI:", req.RequestURI)
	infoLog.Printf("req:%+v", req)
}

func other(res http.ResponseWriter, req *http.Request) {

	//logRequest(req)
	infoLog.Println("other")

	res.WriteHeader(http.StatusBadRequest)
}

func gauge(res http.ResponseWriter, req *http.Request) {

	//logRequest(req)

	re := regexp.MustCompile(`^/update/gauge/(?P<Name>[^/]*)(?:/)?(?P<Value>[^/]*)$`)
	nameIndex := re.SubexpIndex("Name")
	valueIndex := re.SubexpIndex("Value")
	matches := re.FindStringSubmatch(req.URL.Path)

	//infoLog.Println("matches:%+v", matches)

	if matches == nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if matches[nameIndex] == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := strconv.ParseFloat(matches[valueIndex], 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func counter(res http.ResponseWriter, req *http.Request) {

	//logRequest(req)

	re := regexp.MustCompile(`^/update/counter/(?P<Name>[^/]*)(?:/)?(?P<Value>[^/]*)?$`)
	nameIndex := re.SubexpIndex("Name")
	valueIndex := re.SubexpIndex("Value")
	matches := re.FindStringSubmatch(req.URL.Path)

	//infoLog.Println("matches:%+v", matches)

	if matches == nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if matches[nameIndex] == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := strconv.Atoi(matches[valueIndex])
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func middleware1(next http.Handler) http.Handler {
	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// здесь пишем логику обработки
		infoLog.Println("middleware1")
		// замыкание: используем ServeHTTP следующего хендлера
		next.ServeHTTP(w, r)
	})
}

func middleware2(next http.Handler) http.Handler {
	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// здесь пишем логику обработки
		infoLog.Println("middleware2")
		// замыкание: используем ServeHTTP следующего хендлера
		next.ServeHTTP(w, r)
	})
}

func main() {
	f, err := os.OpenFile("./server.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog = log.New(f, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	//infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog.Println("Log started")

	mux := http.NewServeMux()
	infoLog.Println("Mux started")

	mux.Handle("/", Conveyor(http.HandlerFunc(other), middleware1, middleware2))
	mux.HandleFunc("/update/gauge/", gauge)
	mux.HandleFunc("/update/counter/", counter)

	err = http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
