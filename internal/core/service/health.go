package service

import (
	"context"
	"github.com/richmondgoh8/go-casbin/graph/model"
	"github.com/richmondgoh8/go-casbin/internal/core/ports"
)

type HealthSvcImpl struct {
	healthRepo ports.HealthRepository
}

type HealthSvc interface {
	GetHeartbeat(ctx context.Context) *model.HealthResponse
}

// NewHealthSvc Takes in Interface, Return Struct
func NewHealthSvc(healthRepo ports.HealthRepository) *HealthSvcImpl {
	return &HealthSvcImpl{
		healthRepo: healthRepo,
	}
}

func (srv *HealthSvcImpl) GetHeartbeat(ctx context.Context) *model.HealthResponse {
	return srv.healthRepo.GetHealth()
}
