package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"google.golang.org/grpc"
)

type GrpcController struct {
	Logger        *log.Logger
	NextRequestID func() string
	Healthy       int64
}

func (c *GrpcController) Shutdown(ctx context.Context, server *grpc.Server) context.Context {
	ctx, done := context.WithCancel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer done()

		<-quit
		signal.Stop(quit)
		close(quit)

		atomic.StoreInt64(&c.Healthy, 0)
		c.Logger.Printf("server is shutting down...\n")

		server.GracefulStop()
	}()

	return ctx
}
