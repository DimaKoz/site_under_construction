package main

import (
	"errors"
	"fmt"
	"github.com/google/logger"
	"html/template"
	"net/http"
	"runtime"
)

func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
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
				ferr, ok := err.(*notFoundError)
				//errState := http.StatusInternalServerError
				if ok {
					fmt.Println("notFoundError", ferr)

					t, err := template.ParseFiles(html404)
					if err != nil {
						http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
						return
					}
					//errState = http.StatusNotFound
					w.WriteHeader(http.StatusNotFound)
					err = t.Execute(w, nil)
					if err != nil {
						http.Error(w, "Something went wrong :(", http.StatusInternalServerError)
						return
					}
					loggingErr(err)
					return
				} else {
					fmt.Println("unknown type of error")
					fmt.Println(err)

					t, err := template.ParseFiles(html500)
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

				}
				loggingErr(err)
				//TODO sendMeMail(err)
				//http.Error(w, err.Error(), errState)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func loggingErr(err error) {
	logger.Error(err.Error())
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)
	logger.Error(string(buf[0:stackSize]))
}