package server

import (
	"context"
	"log"

	"github.com/jameswspears/go-server/protos"
)

// ToggleButtonServer is used to implement protos.ToggleButtonServer.
type ImplementedToggleButtonServer struct {
	protos.UnimplementedToggleButtonServer
}

// HealthCheck implements protos.HealthCheck
func (s *ImplementedToggleButtonServer) HealthCheck(ctx context.Context, in *protos.HealthCheckRequest) (*protos.HealthCheckResponse, error) {
	log.Print("HealthCheck called")
	return &protos.HealthCheckResponse{Alive: true}, nil
}
