package main

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
	"under_construction/app"
	"under_construction/app/app_errors"
	"under_construction/app/handlers"
)


var log *logger.Logger = nil

func main() {

	lf, errLog := os.OpenFile(app.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
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
	router.NotFoundHandler = app.RecoverWrap(http.HandlerFunc(requestPanic))
	router.Handle(app.PathPatternUnknownError, app.RecoverWrap(http.HandlerFunc(requestUnknownError)))
	router.Handle(app.PathPatternRoot, app.RecoverWrap(app.CheckCache(http.HandlerFunc(handlers.RootHandler), false)))
	router.Handle(app.PathPatternNotFound, app.RecoverWrap(http.HandlerFunc(requestPanic)))
	router.Handle(app.PathPatternFavicon, app.RecoverWrap(app.CheckCache(http.HandlerFunc(handlers.ServeFavicon), false)))
	router.PathPrefix(app.PathPatternWoff2).Handler(app.RecoverWrap(app.CheckCache(http.HandlerFunc(handlers.ServeStatic), true)))
	router.PathPrefix(app.PathPatternCss).Handler(app.RecoverWrap(app.CheckCache(http.HandlerFunc(handlers.ServeStatic), true)))
	router.PathPrefix(app.PathPatternJs).Handler(app.RecoverWrap(app.CheckCache(http.HandlerFunc(handlers.ServeStatic), true)))
	router.PathPrefix(app.PathPatternImage).Handler(app.RecoverWrap(app.CheckCache(http.HandlerFunc(handlers.ServeStatic), true)))
	return router
}

func requestUnknownError(_ http.ResponseWriter, _ *http.Request) {
	panic("oops")
}

func requestPanic(_ http.ResponseWriter, _ *http.Request) {
	panic(app_errors.NewNotFoundError())
}

