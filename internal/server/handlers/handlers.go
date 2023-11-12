package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/k0st1a/metrics/internal/logger"
)

func Counter(res http.ResponseWriter, req *http.Request) {

	logger.LogHTTPRequest(req)

	re := regexp.MustCompile(`^/update/counter/(?P<Name>[^/]*)(?:/)?(?P<Value>[^/]*)?$`)
	nameIndex := re.SubexpIndex("Name")
	valueIndex := re.SubexpIndex("Value")
	matches := re.FindStringSubmatch(req.URL.Path)

	logger.Println("matches:", matches)

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

func Gauge(res http.ResponseWriter, req *http.Request) {

	logger.LogHTTPRequest(req)

	re := regexp.MustCompile(`^/update/gauge/(?P<Name>[^/]*)(?:/)?(?P<Value>[^/]*)$`)
	nameIndex := re.SubexpIndex("Name")
	valueIndex := re.SubexpIndex("Value")
	matches := re.FindStringSubmatch(req.URL.Path)

	logger.Println("URL:", req.URL.Path)
	logger.Println("matches:", matches)

	if matches == nil {
		http.Error(res, "bad gauge request", http.StatusBadRequest)
		return
	}

	if matches[nameIndex] == "" {
		http.Error(res, "gauge name is empty", http.StatusNotFound)
		return
	}

	if matches[valueIndex] == "" {
		http.Error(res, "gauge value is empty", http.StatusBadRequest)
		return
	}
	_, err := strconv.ParseFloat(matches[valueIndex], 64)
	if err != nil {
		http.Error(res, "gauge value is bad", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
}

func Stub(res http.ResponseWriter, req *http.Request) {

	logger.LogHTTPRequest(req)
	logger.Println("stub")

	res.WriteHeader(http.StatusBadRequest)
}
