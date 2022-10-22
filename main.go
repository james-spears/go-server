package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jameswspears/go-server/protos"
	"github.com/jameswspears/go-server/server"
	"google.golang.org/grpc"
)

func GrpcServer(w *io.Writer, wg *sync.WaitGroup) {

	// handle logging
	logger := log.New(*w, "grpc: ", log.LstdFlags)
	logger.Printf("grpc server is starting...")

	// setup server
	port := ":9090"
	env, ok := os.LookupEnv("TOGGLEBUTTON_GRPC_PORT")
	if ok {
		port = fmt.Sprintf(":%s", env)
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	c := &server.GrpcController{Logger: logger, NextRequestID: func() string { return strconv.FormatInt(time.Now().UnixNano(), 36) }}
	s := grpc.NewServer(grpc.UnaryInterceptor(server.UnaryInterceptorFactory(c)))
	ctx := c.Shutdown(context.Background(), s)

	logger.Printf("grpc server is ready to handle requests at %q", port)
	atomic.StoreInt64(&c.Healthy, time.Now().UnixNano())

	protos.RegisterToggleButtonServer(s, &server.ImplementedToggleButtonServer{})
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("grpc server could not listen on %q: %s", port, err)
	}
	<-ctx.Done()
	logger.Printf("grpc server is stopped")
	wg.Done()
}

func HttpServer(w *io.Writer, wg *sync.WaitGroup) {

	// handle logging
	logger := log.New(*w, "http: ", log.LstdFlags)
	logger.Printf("http server is starting...")

	// set up server
	httpAddr := ":8080"
	env, ok := os.LookupEnv("TOGGLEBUTTON_HTTP_PORT")
	if ok {
		httpAddr = fmt.Sprintf(":%s", env)
	}

	c := &server.HttpController{Logger: logger, NextRequestID: func() string { return strconv.FormatInt(time.Now().UnixNano(), 36) }}
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.Index)
	mux.HandleFunc("/healthz", c.Healthz)

	server := &http.Server{
		Addr:         httpAddr,
		Handler:      (server.Interceptors{c.Tracing, c.Logging}).Apply(mux),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	ctx := c.Shutdown(context.Background(), server)

	logger.Printf("http server is ready to handle requests at %q", httpAddr)
	atomic.StoreInt64(&c.Healthy, time.Now().UnixNano())

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("http server could not listen on %q: %s", httpAddr, err)
	}
	<-ctx.Done()
	logger.Printf("http server stopped: OK")
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
	wg.Add(2)

	go HttpServer(&w, &wg)
	go GrpcServer(&w, &wg)
	wg.Wait()
}
