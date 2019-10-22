package handlers

import (
	"fmt"
	log "github.com/google/logger"
	"html/template"
	"net/http"
	"sync"
	"under_construction/app"
	err2 "under_construction/app/apperrors"
)

var mut sync.Mutex
var templ *template.Template

// RootHandler serves for "/" path
func RootHandler(w http.ResponseWriter, r *http.Request) {
	mut.Lock()
	defer mut.Unlock()
	if templ == nil {
		logMessage := fmt.Sprintf("method:[%s], path:[%s]", r.Method, r.URL.Path) //get request method

		if r.URL.Path != "/" {
			log.Warningln(logMessage)
			panic(err2.NewNotFoundError())
		}

		data, err := app.GetBytes(app.HtmlUnderConstruction)
		if err != nil {
			log.Warningln(logMessage)
			panic(err2.NewNotFoundError())
		}
		strData := string(*data)
		templ, err = template.New("root").Parse(strData)
		if err != nil {
			panic(err2.NewNotFoundError())
		}
	}
	w.WriteHeader(http.StatusOK)
	err := templ.Execute(w, nil)
	if err != nil {
		var str string
		str = fmt.Sprintf("unknown error[%s]", err.Error())
		panic(str)
	}

}
