package main

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"net/http"
	"os"
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
	html500                 = "./html/error_500_page.html"
	html404                 = "./html/error_404_page.html"
	htmlUnderConstruction   = "./html/under_construction.html"
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
	router.Handle(pathPatternFavicon, RecoverWrap(checkCache(http.HandlerFunc(favicon), false)))
	router.PathPrefix(pathPatternWoff2).Handler(RecoverWrap(checkCache(http.HandlerFunc(serveStatic), true)))
	router.PathPrefix(pathPatternCss).Handler(RecoverWrap(checkCache(http.HandlerFunc(serveStatic), true)))
	router.PathPrefix(pathPatternJs).Handler(RecoverWrap(checkCache(http.HandlerFunc(serveStatic), true)))
	router.PathPrefix(pathPatternImage).Handler(RecoverWrap(checkCache(http.HandlerFunc(serveStatic), true)))
	return router
}

func requestUnknownError(_ http.ResponseWriter, _ *http.Request) {
	panic("oops")
}

func requestPanic(_ http.ResponseWriter, _ *http.Request) {
	panic(newNotFoundError())
}

