package main

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
	"under_construction/app"
	"under_construction/app/apperrors"
	"under_construction/app/handlers"
	"under_construction/app/middleware"
)

var log *logger.Logger

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
	router.NotFoundHandler = middleware.RecoverWrap(http.HandlerFunc(requestPanic))
	router.Handle(app.PathPatternUnknownError, middleware.RecoverWrap(http.HandlerFunc(requestUnknownError)))
	router.Handle(app.PathPatternRoot, middleware.RecoverWrap(middleware.CheckCache(http.HandlerFunc(handlers.RootHandler), false)))
	router.Handle(app.PathPatternNotFound, middleware.RecoverWrap(http.HandlerFunc(requestPanic)))
	router.Handle(app.PathPatternFavicon, middleware.RecoverWrap(middleware.CheckCache(http.HandlerFunc(handlers.ServeFavicon), false)))
	router.PathPrefix(app.PathPatternWoff2).Handler(middleware.RecoverWrap(middleware.CheckCache(http.HandlerFunc(handlers.ServeStatic), true)))
	router.PathPrefix(app.PathPatternCss).Handler(middleware.RecoverWrap(middleware.CheckCache(http.HandlerFunc(handlers.ServeStatic), true)))
	router.PathPrefix(app.PathPatternJs).Handler(middleware.RecoverWrap(middleware.CheckCache(http.HandlerFunc(handlers.ServeStatic), true)))
	router.PathPrefix(app.PathPatternImage).Handler(middleware.RecoverWrap(middleware.CheckCache(http.HandlerFunc(handlers.ServeStatic), true)))
	return router
}

func requestUnknownError(_ http.ResponseWriter, _ *http.Request) {
	panic("oops")
}

func requestPanic(_ http.ResponseWriter, _ *http.Request) {
	panic(apperrors.NewNotFoundError())
}
