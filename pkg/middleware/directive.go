package custommiddleware

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	ErrAccessDenied string = "access denied"
)

func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	tokenData, _ := GetClaimsFromJWTTokenCtx(ctx)
	if tokenData == nil {
		return nil, &gqlerror.Error{
			Message: ErrAccessDenied,
		}
	}

	return next(ctx)
}
