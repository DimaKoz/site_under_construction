package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {

	logMessage := fmt.Sprintf("method:[%s], path:[%s]", r.Method, r.URL.Path) //get request method

	if r.URL.Path != "/" {
		log.Warningln(logMessage)
		panic(newNotFoundError())
	}

	t, err := template.ParseFiles(htmlUnderConstruction)
	if err != nil {
		panic(newNotFoundError())
		return
	}
	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, nil)
	if err != nil {
		var str string
		str = fmt.Sprintf("unknown error[%s]", err.Error())
		panic(str)
		return
	}

}
