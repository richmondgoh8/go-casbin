package ports

import (
	"github.com/richmondgoh8/go-casbin/graph/model"
)

type HealthRepository interface {
	GetHealth() *model.HealthResponse
}
