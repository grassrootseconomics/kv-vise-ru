package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/grassrootseconomics/kv-vise-ru/pkg/store"
	"github.com/kamikazechaser/common/httputil"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

type (
	APIOpts struct {
		EnableMetrics bool
		ListenAddress string
		CORS          []string
		Logg          *slog.Logger
		Store         store.Store
	}

	API struct {
		validator httputil.ValidatorProvider
		router    *bunrouter.Router
		server    *http.Server
		logg      *slog.Logger
		store     store.Store
	}
)

const (
	apiVersion = "/api/v1"
	slaTimeout = 10 * time.Second
)

func New(o APIOpts) *API {
	api := &API{
		validator: httputil.NewValidator(""),
		logg:      o.Logg,
		store:     o.Store,
		router: bunrouter.New(
			bunrouter.WithNotFoundHandler(notFoundHandler),
			bunrouter.WithMethodNotAllowedHandler(methodNotAllowedHandler),
		),
	}

	if o.EnableMetrics {
		api.router.GET("/metrics", metricsHandler)
	}

	api.router.Use(newCorsMiddleware(o.CORS)).WithGroup(apiVersion, func(g *bunrouter.Group) {
		if os.Getenv("DEV") != "" {
			g = g.Use(reqlog.NewMiddleware())
		}

		g.GET("/lookup/address/:phone", api.addressHandler)
	})

	api.server = &http.Server{
		Addr:    o.ListenAddress,
		Handler: api.router,
	}

	return api
}

func (a *API) Start() error {
	a.logg.Info("API server starting", "address", a.server.Addr)
	if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *API) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
