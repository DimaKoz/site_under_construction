package main

import (
	"about/etagging"
	"io/ioutil"
	"net/http"
	"strings"
)

func checkCache(h http.Handler, isStatic bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		var err error
		if !isStatic && r.URL.Path == "/" {
			data, err = ioutil.ReadFile(htmlUnderConstruction)
		} else {
			path := r.URL.Path[1:]
			data, err = ioutil.ReadFile(path)
		}

		if err == nil {
			etagValue := etagging.Generate(string(data), true)
			if match := r.Header.Get("If-None-Match"); match != "" {
				if strings.Contains(match, etagValue) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
			w.Header().Set("Cache-Control", "max-age=3600")
			w.Header().Set("Etag", etagValue)
			h.ServeHTTP(w, r)
		} else {
			panic(newNotFoundError())
		}
	})
}
