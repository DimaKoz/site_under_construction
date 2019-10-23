package middleware

import (
	"net/http"
	"strings"
	"under_construction/app"
	"under_construction/app/apperrors"
	"under_construction/app/etagging"
)

// CheckCache is a middleware function that checks 'If-None-Match' header
// and sets http.StatusNotModified when it possible
func CheckCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data *[]byte
		var strData string
		var err error
		if r.URL.Path == "/" {
			data, err = app.GetBytes(app.HtmlUnderConstruction)
			if err == nil {
				strData = string(*data)
			}
		} else if r.URL.Path == "/favicon.ico" {
			strData = app.FaviconData
		} else {
			path := r.URL.Path[1:]
			data, err = app.GetBytes(path)
			if err == nil {
				strData = string(*data)
			}
		}

		if err == nil {
			etagValue := etagging.Generate(strData, true)
			if match := r.Header.Get(app.HeaderKeyIfNoneMatch); match != "" {
				if strings.Contains(match, etagValue) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
			w.Header().Set(app.HeaderKeyCacheControl, app.HeaderValueCacheControl)
			w.Header().Set(app.HeaderKeyEtag, etagValue)
			h.ServeHTTP(w, r)
		} else {
			panic(apperrors.NewNotFoundError())
		}
	})
}
