package server

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func UnaryInterceptorFactory(c *GrpcController) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()

		h, err := handler(ctx, req)

		//logging
		c.Logger.Println(ctx, time.Since(start))

		return h, err
	}
}
