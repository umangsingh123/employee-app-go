package router

import (
	"net/http"

	"emplopyee-app-go/internal/handler"
	"emplopyee-app-go/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Package router provides the application's HTTP routing and middleware configuration.
//
// This file defines the router that handles API versioning, route grouping, and
// middleware setup for the Employee management service.
//
// The NewRouter function sets up a chi.Router with logging, recovery, timeout, and
// request ID middleware, and mounts the Employee resource handlers at /api/v1/employees.
//
// Routes:
//
//	POST   /api/v1/employees/          - Create new employee
//	GET    /api/v1/employees/          - List all employees
//	GET    /api/v1/employees/{id}/     - Get employee by ID
//	PUT    /api/v1/employees/{id}/     - Update employee by ID
//	DELETE /api/v1/employees/{id}/     - Delete employee by ID
//	GET    /health                     - Health check endpoint
//
// Dependencies:
//   - github.com/go-chi/chi/v5: Routing framework
//   - github.com/go-chi/chi/v5/middleware: Middleware for logging, recovery, timeouts, etc.
func NewRouter(svc service.EmployeeService) http.Handler {
	// Initialize a new chi.Router instance to handle incoming HTTP requests
	r := chi.NewRouter()

	// Register common middleware for all routes:
	//   - RequestID:       Assigns a unique request ID to each HTTP request for tracking/logging.
	//   - Logger:          Logs the start and end of each request, including path, method, duration, and status.
	//   - Recoverer:       Recovers from panics within handlers and returns a 500 error instead of crashing the server.
	//   - Timeout(30s):    Ensures that handlers take no more than 30 seconds to respond, preventing resource exhaustion.
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * 1e9)) // 30s

	h := handler.NewEmployeeHandler(svc)

	r.Route("/api/v1/employees", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Route("/{id:[0-9]+}", func(r chi.Router) {
			r.Get("/", h.Get)
			r.Put("/", h.Update)
			r.Delete("/", h.Delete)
		})
	})

	// health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	return r
}
