package custommiddleware

import (
	"context"
	"github.com/richmondgoh8/go-casbin/pkg/middleware/logger"
	"github.com/richmondgoh8/go-casbin/pkg/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

type contextKey struct {
	Name string
}

var (
	// TrackingCtxKey is the context.Context key to store the tracking id for a request.
	TrackingCtxKey = &contextKey{"tracking_id"}
)

func GenTrackingID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// set to context map
		ctx = context.WithValue(ctx, TrackingCtxKey.Name, utils.GenerateTrackingID())
		r = r.WithContext(ctx)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t1 := time.Now()
		defer func() {
			logger.Info("route travel", ctx, map[string]interface{}{
				"latency": time.Since(t1),
				"method":  r.Method,
				"status":  ww.Status(),
			})
		}()

		next.ServeHTTP(ww, r)
	})
}
