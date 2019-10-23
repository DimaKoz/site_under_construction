package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
	"under_construction/app"
	"under_construction/app/apperrors"
	"under_construction/app/handlers"
	"under_construction/app/middleware"
)

var wait sync.WaitGroup

// I've used an idea from https://stackoverflow.com/a/42533360/3166697
func startHttpServer(waitGroup *sync.WaitGroup, ch *chan error) *http.Server {

	router := defaultMux()

	address := ":8000" //"127.0.0.1:8000"
	srv := &http.Server{
		Handler: router,
		Addr:    address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		defer func() {
			if waitGroup != nil {
				waitGroup.Done()
			}
		}()

		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			if ch != nil {
				*ch <- err
			}

		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func main() {
	channelError := make(chan error)

	initLogger()
	defer closeLogger()
	wait.Add(1)
	_ = startHttpServer(&wait, &channelError)
	wait.Wait()

	err := <-channelError
	log.Fatalf("ListenAndServe(): %s", err)
}

func defaultMux() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = middleware.RecoverWrap(http.HandlerFunc(requestPanic))
	router.Handle(app.PathPatternUnknownError, middleware.RecoverWrap(http.HandlerFunc(requestUnknownError)))
	router.Handle(app.PathPatternRoot, middleware.RecoverWrap(middleware.GzipWrapper(middleware.CheckCache(http.HandlerFunc(handlers.RootHandler)))))
	router.Handle(app.PathPatternNotFound, middleware.RecoverWrap(http.HandlerFunc(requestPanic)))
	router.Handle(app.PathPatternFavicon, middleware.RecoverWrap(middleware.CheckCache(http.HandlerFunc(handlers.ServeFavicon))))
	router.PathPrefix(app.PathPatternWoff2).Handler(middleware.RecoverWrap(middleware.CheckCache(http.HandlerFunc(handlers.ServeStatic))))
	router.PathPrefix(app.PathPatternCss).Handler(middleware.RecoverWrap(middleware.GzipWrapper(middleware.CheckCache(http.HandlerFunc(handlers.ServeStatic)))))
	router.PathPrefix(app.PathPatternJs).Handler(middleware.RecoverWrap(middleware.GzipWrapper(middleware.CheckCache(http.HandlerFunc(handlers.ServeStatic)))))
	router.PathPrefix(app.PathPatternImage).Handler(middleware.RecoverWrap(middleware.CheckCache(http.HandlerFunc(handlers.ServeStatic))))
	return router
}

func requestUnknownError(_ http.ResponseWriter, _ *http.Request) {
	panic("oops")
}

func requestPanic(_ http.ResponseWriter, _ *http.Request) {
	panic(apperrors.NewNotFoundError())
}
