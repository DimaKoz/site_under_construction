package handlers

import (
	"encoding/base64"
	"net/http"
	"under_construction/app"
)


func ServeFavicon(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "image/x-icon")
	p, err := base64.StdEncoding.DecodeString(app.FaviconData)
	_ = err
	if _, err := w.Write(p); err != nil {
		panic(err)
	}
}
