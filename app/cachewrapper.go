package app

import (
	"net/http"
	"strings"
	"under_construction/app/apperrors"
	"under_construction/app/etagging"
)

func CheckCache(h http.Handler, isStatic bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data *[]byte
		var strData string
		var err error
		if !isStatic && r.URL.Path == "/" {
			data, err = GetBytes(HtmlUnderConstruction)
			if err == nil {
				strData = string(*data)
			}
		} else if !isStatic && r.URL.Path == "/favicon.ico" {
			strData = FaviconData
		} else {
			path := r.URL.Path[1:]
			data, err = GetBytes(path)
			if err == nil {
				strData = string(*data)
			}
		}

		if err == nil {
			etagValue := etagging.Generate(strData, true)
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
			panic(apperrors.NewNotFoundError())
		}
	})
}
