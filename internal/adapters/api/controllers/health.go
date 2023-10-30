package controllers

import (
	"net/http"

	"github.com/zeusito/go-blueprint/internal/adapters/api"
)

type HealthRoutes struct{}

func NewHealthRoutes(server *api.HTTPServer) *HealthRoutes {
	c := &HealthRoutes{}

	// declare routes
	server.Router.Get("/health/liveness", c.handleLivenessCheck)
	server.Router.Get("/health/readiness", c.handleReadinessCheck)

	return c
}

func (c *HealthRoutes) handleLivenessCheck(w http.ResponseWriter, r *http.Request) {
	api.RenderJSON(r.Context(), w, http.StatusOK, map[string]string{"status": "ok"})
}

func (c *HealthRoutes) handleReadinessCheck(w http.ResponseWriter, r *http.Request) {
	api.RenderJSON(r.Context(), w, http.StatusOK, map[string]string{"status": "ok"})
}
