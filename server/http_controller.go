package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type HttpController struct {
	Logger        *log.Logger
	NextRequestID func() string
	Healthy       int64
}

func (c *HttpController) Shutdown(ctx context.Context, server *http.Server) context.Context {
	ctx, done := context.WithCancel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer done()

		<-quit
		signal.Stop(quit)
		close(quit)

		atomic.StoreInt64(&c.Healthy, 0)
		server.ErrorLog.Printf("server is shutting down...\n")

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			server.ErrorLog.Fatalf("could not gracefully shutdown the server: %s", err)
		}
	}()

	return ctx
}

func (c *HttpController) Index(w http.ResponseWriter, req *http.Request) {
	RootHandler(w, req)
}

// Health struct represents a healthy response.
type Health struct {
	Uptime time.Duration `json:"uptime"`
}

func (c *HttpController) Healthz(w http.ResponseWriter, req *http.Request) {
	if h := atomic.LoadInt64(&c.Healthy); h == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(Health{Uptime: time.Since(time.Unix(0, h))})
	}
}

var (
	_ http.Handler = http.HandlerFunc((&HttpController{}).Index)
	_ http.Handler = http.HandlerFunc((&HttpController{}).Healthz)
)
