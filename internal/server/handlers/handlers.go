package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/k0st1a/metrics/internal/logger"
)

func Counter(res http.ResponseWriter, req *http.Request) {

	logger.LogHttpRequest(req)

	re := regexp.MustCompile(`^/update/counter/(?P<Name>[^/]*)(?:/)?(?P<Value>[^/]*)?$`)
	nameIndex := re.SubexpIndex("Name")
	valueIndex := re.SubexpIndex("Value")
	matches := re.FindStringSubmatch(req.URL.Path)

	logger.Println("matches:%+v", matches)

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

	logger.LogHttpRequest(req)

	re := regexp.MustCompile(`^/update/gauge/(?P<Name>[^/]*)(?:/)?(?P<Value>[^/]*)$`)
	nameIndex := re.SubexpIndex("Name")
	valueIndex := re.SubexpIndex("Value")
	matches := re.FindStringSubmatch(req.URL.Path)

	logger.Println("matches:%+v", matches)

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

func Stub(res http.ResponseWriter, req *http.Request) {

	logger.LogHttpRequest(req)
	logger.Println("stub")

	res.WriteHeader(http.StatusBadRequest)
}
