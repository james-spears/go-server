package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jameswspears/go-server/service"
)

func HttpServer(w *io.Writer, wg *sync.WaitGroup) {

	// handle logging
	logger := log.New(*w, "go_server_http: ", log.LstdFlags)
	logger.Printf("server is starting...")

	// set up server
	httpAddr := ":8080"
	env, ok := os.LookupEnv("PORT")
	if ok {
		httpAddr = fmt.Sprintf(":%s", env)
	}

	c := &service.HttpController{Logger: logger, NextRequestID: func() string { return strconv.FormatInt(time.Now().UnixNano(), 36) }}
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.Index)
	mux.HandleFunc("/healthz", c.Healthz)

	server := &http.Server{
		Addr:         httpAddr,
		Handler:      (service.Interceptors{c.Tracing, c.Logging}).Apply(mux),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	ctx := c.Shutdown(context.Background(), server)

	logger.Printf("server is ready to handle requests at %q", httpAddr)
	atomic.StoreInt64(&c.Healthy, time.Now().UnixNano())

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("could not listen on %q: %s", httpAddr, err)
	}
	<-ctx.Done()
	logger.Printf("server stopped: OK")
	wg.Done()
}

func main() {
	// setup combined logging
	logFile := "go_server.log"

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("%v:", err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	w := io.MultiWriter(os.Stdout, f)
	log.SetOutput(w)

	var wg sync.WaitGroup
	wg.Add(1)

	go HttpServer(&w, &wg)
	wg.Wait()
}
