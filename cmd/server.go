package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/richmondgoh8/go-casbin/graph"
	"github.com/richmondgoh8/go-casbin/graph/model"
	"github.com/richmondgoh8/go-casbin/internal/core/service"
	"github.com/richmondgoh8/go-casbin/internal/platform/config"
	"github.com/richmondgoh8/go-casbin/internal/repositories"
	"github.com/richmondgoh8/go-casbin/pkg/client/postgres"
	custommiddleware "github.com/richmondgoh8/go-casbin/pkg/middleware"
	adapter "github.com/richmondgoh8/go-casbin/pkg/middleware/casbin"
	"github.com/richmondgoh8/go-casbin/pkg/middleware/logger"
	apperror "github.com/richmondgoh8/go-casbin/static"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	config.InitReader()
	logger.Init([]string{"tracking_id", "endpoint"})

	appConfig := config.Init()
	port := os.Getenv("PORT")
	if port == "" {
		panic(apperror.EmptyPort)
	}

	r := chi.NewRouter()
	r.Use(cors.Default().Handler)
	r.Use(custommiddleware.Auth)

	// Start of Third Party Setup

	dbClient, err := postgres.Init(appConfig.DB)
	if err != nil {
		panic(err)
	}
	err = dbClient.Ping()
	if err != nil {
		panic(err)
	}

	if err = adapter.InitializeCasbin(dbClient); err != nil {
		panic(err)
	}
	// End of Third Party Setup

	// Start of Dependency Injection
	healthPort := repositories.NewHealthPort()
	healthSvc := service.NewHealthSvc(healthPort)

	//usersPort := repositories.NewUserPort(dbClient, appConfig.DB)
	//usersSvc := service.NewUserSvc(usersPort)
	// End of Dependency Injection

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			HealthSvc: healthSvc,
		},
		Directives: graph.DirectiveRoot{
			Auth: custommiddleware.AuthDirective,
			HasRole: func(ctx context.Context, obj interface{}, next graphql.Resolver, requestObj model.RequestObj, requestAction model.RequestAction) (res interface{}, err error) {
				payload, err := custommiddleware.GetClaimsFromJWTTokenCtx(ctx)
				if err != nil {
					return nil, err
				}

				if isAllowed, _ := adapter.Policy().Enforce(strings.ToLower(payload.Role), strings.ToLower(requestObj.String()), strings.ToLower(requestAction.String())); !isAllowed {
					// block calling the next resolver
					return nil, errors.New(fmt.Sprintf("no such authority for %s", payload.Role))
				}
				return next(ctx)
			},
		},
	}))

	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		gqlErr := graphql.DefaultErrorPresenter(ctx, e)
		return gqlErr
	})

	srv.SetRecoverFunc(func(ctx context.Context, e interface{}) (userMessage error) {
		return gqlerror.Errorf("internal server error")
	})

	// Disable introspection if DISABLE_INTROSPECTION is set to true
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		if strings.EqualFold(os.Getenv("DISABLE_INTROSPECTION"), "true") {
			graphql.GetOperationContext(ctx).DisableIntrospection = true
		}
		return next(ctx)
	})

	r.Handle("/query", srv)
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
