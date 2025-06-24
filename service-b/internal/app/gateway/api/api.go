package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/config"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/domain/usecase"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/gateway/api/handler"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/gateway/api/middleware"
)

type API struct {
	Handler http.Handler
	cfg     config.Config
	useCase usecase.UseCaseInterface
}

func BasicHandler() http.Handler {
	router := chi.NewMux()
	handler.RegisterHealthCheckRoute(router)

	return router
}

func New(cfg config.Config, useCase usecase.UseCaseInterface) *API {
	api := &API{
		cfg:     cfg,
		useCase: useCase,
	}

	api.setupRouter()

	return api
}

func (api *API) setupRouter() {
	router := chi.NewRouter()

	api.registerRoutes(router)

	api.Handler = router
}

func (api *API) registerRoutes(router *chi.Mux) {
	handler.RegisterHealthCheckRoute(router)

	router.Route("/", func(r chi.Router) {
		r.Use(middleware.Logger, middleware.RealIP, middleware.ContentTypeJSON)
		handler.RegisterRoutes(
			r,
			api.cfg,
			api.useCase,
		)
	})
}
