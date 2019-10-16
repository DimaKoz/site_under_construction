package handlers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"under_construction/app"
	err2 "under_construction/app/app_errors"

)

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	//TODO etag , for example https://github.com/bouk/staticfiles/blob/master/files/files.go and https://github.com/dc0d/cache-control

	path := r.URL.Path[1:]
	files, err := ioutil.ReadDir("./")
	println(files)

	data, err := app.GetBytes(path)
	if err == nil {
		contentType := getContentType(path)
		w.Header().Add("Content-Type", contentType)

		if _, err := w.Write(*data); err != nil {
			panic(err)
		}
	} else {
		panic(err2.NewNotFoundError())
	}
}


func getContentType(path string) string {
	var contentType string

	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(path, ".woff2") {
		contentType = "font/woff2"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(path, ".jpg") {
		contentType = "image/jpg"
	} else if strings.HasSuffix(path, ".svg") {
		contentType = "image/svg+xml"
	} else {
		contentType = "text/plain"
	}

	return contentType
}