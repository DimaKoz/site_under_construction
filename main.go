package main

import (
	"about/etagging"
	"errors"
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	logPath = "log.txt"

	pathPatternRoot         = "/"
	pathPatternNotFound     = "/404"
	pathPatternUnknownError = "/500"
	pathPatternFavicon      = "/favicon.ico"
	pathPatternWoff2        = "/assets/woff2/"
	pathPatternCss          = "/assets/css/"
	pathPatternJs           = "/assets/script/"
	pathPatternImage        = "/assets/image/"
)

var log *logger.Logger = nil

func main() {

	lf, errLog := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if errLog != nil {
		logger.Fatalf("Failed to open log file: %v", errLog)
	}
	defer lf.Close()

	log = logger.Init("LoggerExample", true, false, lf)
	defer log.Close()
	logger.Warningln("")
	logger.Warningln("================================================================")
	logger.Warningln("")
	logger.Warningln("Logger started")
	logger.Warningln("")
	logger.Warningln("================================================================")

	router := defaultMux()

	address := ":8000" //"127.0.0.1:8000"
	srv := &http.Server{
		Handler: router,
		Addr:    address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func defaultMux() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = RecoverWrap(http.HandlerFunc(requestPanic))
	router.Handle(pathPatternUnknownError, RecoverWrap(http.HandlerFunc(requestUnknownError)))
	router.Handle(pathPatternRoot, RecoverWrap(checkCache(http.HandlerFunc(rootHandler), false)))
	router.Handle(pathPatternNotFound, RecoverWrap(http.HandlerFunc(requestPanic)))
	router.Handle(pathPatternFavicon, RecoverWrap(http.HandlerFunc(favicon)))
	router.PathPrefix(pathPatternWoff2).Handler(RecoverWrap(checkCache(http.HandlerFunc(serveStatic), true)))
	router.PathPrefix(pathPatternCss).Handler(RecoverWrap(checkCache(http.HandlerFunc(serveStatic), true)))
	router.PathPrefix(pathPatternJs).Handler(RecoverWrap(checkCache(http.HandlerFunc(serveStatic), true)))
	router.PathPrefix(pathPatternImage).Handler(RecoverWrap(checkCache(http.HandlerFunc(serveStatic), true)))
	return router
}

func requestUnknownError(w http.ResponseWriter, r *http.Request) {
	panic("oops")
}

func requestPanic(w http.ResponseWriter, r *http.Request) {
	panic(newNotFoundError())
}

func checkCache(h http.Handler, isStatic bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		var err error
		if (!isStatic && r.URL.Path == "/") {
			data, err = ioutil.ReadFile("./html/under_construction.html")
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

					t, err := template.ParseFiles("./html/error_404_page.html")
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
					return
				} else {
					fmt.Println("unknown type of error")
					fmt.Println(err)

					t, err := template.ParseFiles("./html/error_500_page.html")
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
