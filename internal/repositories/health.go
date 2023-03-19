package repositories

import (
	"github.com/richmondgoh8/go-casbin/graph/model"
)

type HealthPortImpl struct{}

func NewHealthPort() *HealthPortImpl {
	return &HealthPortImpl{}
}

func (h *HealthPortImpl) GetHealth() *model.HealthResponse {
	return &model.HealthResponse{
		Message: "OK",
	}
}
