package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMainFunc(t *testing.T) {

	initLogger()
	defer closeLogger()

	srv := startHttpServer(nil, nil)
	workingTime := 1
	logStr := fmt.Sprintf("main: serving for %d second", workingTime)
	log.Info(logStr)

	duration := time.Duration(workingTime) * time.Second
	time.Sleep(time.Duration(duration))

	log.Info("main: stopping HTTP server")
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	log.Info("main: done. exiting")
}
