// Package gateway provides the REST API Gateway layer that handles HTTP routing
// and JSON to Protobuf conversion for the IPAM service.
package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bananaops/ipam-bananaops/internal/service"
	pb "github.com/bananaops/ipam-bananaops/proto"
	"github.com/gorilla/mux"
)

// RESTGateway handles HTTP REST requests and converts between JSON and Protobuf
type RESTGateway struct {
	serviceLayer *service.ServiceLayer
	router       *mux.Router
}

// NewRESTGateway creates a new REST gateway instance
func NewRESTGateway(serviceLayer *service.ServiceLayer) *RESTGateway {
	g := &RESTGateway{
		serviceLayer: serviceLayer,
		router:       mux.NewRouter(),
	}
	g.setupRoutes()
	return g
}

// setupRoutes configures all REST API routes
func (g *RESTGateway) setupRoutes() {
	// API v1 routes
	api := g.router.PathPrefix("/api/v1").Subrouter()

	// Subnet endpoints
	api.HandleFunc("/subnets", g.handleCreateSubnet).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/subnets", g.handleListSubnets).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/subnets/{id}", g.handleGetSubnet).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/subnets/{id}", g.handleUpdateSubnet).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/subnets/{id}", g.handleDeleteSubnet).Methods(http.MethodDelete, http.MethodOptions)

	// Health check endpoints
	g.router.HandleFunc("/health", g.handleHealth).Methods(http.MethodGet)
	g.router.HandleFunc("/ready", g.handleReady).Methods(http.MethodGet)
}

// Handler returns the HTTP handler with CORS middleware
func (g *RESTGateway) Handler() http.Handler {
	return g.corsMiddleware(g.router)
}

// corsMiddleware adds CORS headers to responses
func (g *RESTGateway) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleHealth returns the health status of the service
func (g *RESTGateway) handleHealth(w http.ResponseWriter, r *http.Request) {
	g.writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// handleReady returns the readiness status of the service
func (g *RESTGateway) handleReady(w http.ResponseWriter, r *http.Request) {
	g.writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// writeJSON writes a JSON response with the given status code
func (g *RESTGateway) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// writeError writes an error response in JSON format
func (g *RESTGateway) writeError(w http.ResponseWriter, status int, code, message string) {
	errResp := &ErrorResponse{
		Error: &ErrorDetail{
			Code:      code,
			Message:   message,
			Timestamp: time.Now().Unix(),
		},
	}
	g.writeJSON(w, status, errResp)
}

// writeProtobufError writes a Protobuf error as JSON response
func (g *RESTGateway) writeProtobufError(w http.ResponseWriter, pbErr *pb.Error) {
	status := g.errorCodeToHTTPStatus(pbErr.Code)
	errResp := &ErrorResponse{
		Error: &ErrorDetail{
			Code:      pbErr.Code,
			Message:   pbErr.Message,
			Details:   pbErr.Details,
			Timestamp: pbErr.Timestamp,
		},
	}
	g.writeJSON(w, status, errResp)
}

// errorCodeToHTTPStatus maps error codes to HTTP status codes
func (g *RESTGateway) errorCodeToHTTPStatus(code string) int {
	switch code {
	case "INVALID_CIDR", "INVALID_IP", "INVALID_REQUEST", "MISSING_FIELD", "INVALID_MESSAGE_FORMAT":
		return http.StatusBadRequest
	case "SUBNET_NOT_FOUND":
		return http.StatusNotFound
	case "DUPLICATE_SUBNET":
		return http.StatusConflict
	case "DB_ERROR", "DB_CONNECTION_ERROR", "CALCULATION_ERROR":
		return http.StatusInternalServerError
	case "PROVIDER_UNAVAILABLE", "PROVIDER_AUTH_FAILED", "PROVIDER_RATE_LIMITED":
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
