package middleware

import (
	"errors"
	"fmt"
	log "github.com/google/logger"
	"html/template"
	"net/http"
	"runtime"
	"under_construction/app"
	"under_construction/app/apperrors"
)

func deferredRecover(w http.ResponseWriter) {
	var err error
	r := recover()
	if r != nil {
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("Unknown error")
		}
		log.Warningln("recover() != nil")
		ferr, ok := err.(*apperrors.NotFoundError)

		if ok {
			fmt.Println("NotFoundError", ferr)
			data, err := app.GetBytes(app.Html404)
			if err != nil {
				http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
				return
			}
			strData := string(*data)
			t, err := template.New("404").Parse(strData)
			if err != nil {
				http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNotFound)
			err = t.Execute(w, nil)
			if err != nil {
				http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
				return
			}
			return
		}

		fmt.Println("unknown type of error")
		fmt.Println(err)
		data, err := app.GetBytes(app.Html500)
		if err != nil {
			http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
			return
		}
		strData := string(*data)
		t, err := template.New("500").Parse(strData)
		//t, err := template.ParseFiles(Html500)
		if err != nil {
			http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
			return
		}

		loggingErr(err)
		//TODO sendMeMail(err)

	}
}

// RecoverWrap is a middleware func which tries to recover panics
func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer deferredRecover(w)
		h.ServeHTTP(w, r)
	})
}

func loggingErr(err error) {
	if err == nil {
		return
	}
	log.Error(err.Error())
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)
	log.Error(string(buf[0:stackSize]))
}
