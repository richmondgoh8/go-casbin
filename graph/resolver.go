//go:generate go run github.com/99designs/gqlgen
package graph

import (
	"github.com/richmondgoh8/go-casbin/internal/core/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	HealthSvc service.HealthSvc
}
