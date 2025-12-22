package gateway

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bananaops/ipam-bananaops/internal/cloudprovider"
	"github.com/bananaops/ipam-bananaops/internal/repository"
	"github.com/bananaops/ipam-bananaops/internal/service"
	pb "github.com/bananaops/ipam-bananaops/proto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Gateway handles HTTP REST requests with cloud provider integration
type Gateway struct {
	serviceLayer *service.ServiceLayer
	cloudManager *cloudprovider.Manager
	router       *mux.Router
}

// NewGateway creates a new gateway instance with cloud provider support
func NewGateway(serviceLayer *service.ServiceLayer, cloudManager *cloudprovider.Manager) *Gateway {
	g := &Gateway{
		serviceLayer: serviceLayer,
		cloudManager: cloudManager,
		router:       mux.NewRouter(),
	}
	g.setupRoutes()
	return g
}

// setupRoutes configures all REST API routes
func (g *Gateway) setupRoutes() {
	// API v1 routes
	api := g.router.PathPrefix("/api/v1").Subrouter()

	// Subnet endpoints
	api.HandleFunc("/subnets", g.handleCreateSubnetRepository).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/subnets", g.handleListSubnetsRepository).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/subnets/{id}", g.handleGetSubnet).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/subnets/{id}", g.handleUpdateSubnet).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/subnets/{id}", g.handleDeleteSubnet).Methods(http.MethodDelete, http.MethodOptions)
	api.HandleFunc("/subnets/{id}/children", g.handleGetSubnetChildren).Methods(http.MethodGet, http.MethodOptions)

	// Cloud provider endpoints
	api.HandleFunc("/cloud/sync", g.HandleCloudSync).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/cloud/status", g.HandleCloudStatus).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/cloud/utilization/update", g.HandleUpdateUtilization).Methods(http.MethodPost, http.MethodOptions)

	// Health check endpoints
	g.router.HandleFunc("/health", g.handleHealth).Methods(http.MethodGet)
	g.router.HandleFunc("/ready", g.handleReady).Methods(http.MethodGet)
}

// Handler returns the HTTP handler with CORS middleware
func (g *Gateway) Handler() http.Handler {
	return g.corsMiddleware(g.router)
}

// corsMiddleware adds CORS headers to responses
func (g *Gateway) corsMiddleware(next http.Handler) http.Handler {
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
func (g *Gateway) handleHealth(w http.ResponseWriter, r *http.Request) {
	g.writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// handleReady returns the readiness status of the service
func (g *Gateway) handleReady(w http.ResponseWriter, r *http.Request) {
	g.writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// writeJSON writes a JSON response with the given status code
func (g *Gateway) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// writeErrorResponse writes an error response in JSON format
func (g *Gateway) writeErrorResponse(w http.ResponseWriter, status int, code, message string, err error) {
	if err != nil {
		log.Printf("Error: %s - %v", message, err)
	}

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
func (g *Gateway) writeProtobufError(w http.ResponseWriter, pbErr *pb.Error) {
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
func (g *Gateway) errorCodeToHTTPStatus(code string) int {
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

// Subnet handlers (copied from handlers.go and adapted)

// handleCreateSubnet handles POST /api/v1/subnets
func (g *Gateway) handleCreateSubnet(w http.ResponseWriter, r *http.Request) {
	log.Println("[CreateSubnet] Received request")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[CreateSubnet] Failed to read body: %v", err)
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to read request body", err)
		return
	}
	defer r.Body.Close()

	log.Printf("[CreateSubnet] Request body: %s", string(body))

	// Validate request body is not empty
	if len(body) == 0 {
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Request body is required", nil)
		return
	}

	// Convert JSON to Protobuf request
	req, err := JSONToCreateSubnetRequest(body)
	if err != nil {
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_MESSAGE_FORMAT", err.Error(), err)
		return
	}

	// Validate required fields
	if req.Cidr == "" {
		g.writeErrorResponse(w, http.StatusBadRequest, "MISSING_FIELD", "CIDR is required", nil)
		return
	}
	if req.Name == "" {
		g.writeErrorResponse(w, http.StatusBadRequest, "MISSING_FIELD", "Name is required", nil)
		return
	}

	log.Printf("[CreateSubnet] Protobuf request: %+v", req)

	// Call service layer
	resp, err := g.serviceLayer.CreateSubnet(r.Context(), req)
	if err != nil {
		log.Printf("[CreateSubnet] Service layer error: %v", err)
		g.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), err)
		return
	}

	// Check for service-level errors
	if resp.Error != nil {
		log.Printf("[CreateSubnet] Service returned error: %+v", resp.Error)
		g.writeProtobufError(w, resp.Error)
		return
	}

	log.Printf("[CreateSubnet] Successfully created subnet: %s", resp.Subnet.Id)

	// Convert response to JSON and send
	jsonSubnet := SubnetToJSON(resp.Subnet)
	g.writeJSON(w, http.StatusCreated, jsonSubnet)
}

// handleListSubnets handles GET /api/v1/subnets
func (g *Gateway) handleListSubnets(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	req := &pb.ListSubnetsRequest{
		LocationFilter:      query.Get("location"),
		CloudProviderFilter: query.Get("cloud_provider"),
		SearchQuery:         query.Get("search"),
		Page:                parseIntParam(query.Get("page"), 0),
		PageSize:            parseIntParam(query.Get("page_size"), 50),
	}

	// Call service layer
	resp, err := g.serviceLayer.ListSubnets(r.Context(), req)
	if err != nil {
		g.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), err)
		return
	}

	// Check for service-level errors
	if resp.Error != nil {
		g.writeProtobufError(w, resp.Error)
		return
	}

	// Convert response to JSON and send
	jsonResp := &ListSubnetsResponseJSON{
		Subnets:    SubnetsToJSON(resp.Subnets),
		TotalCount: resp.TotalCount,
	}
	g.writeJSON(w, http.StatusOK, jsonResp)
}

// handleGetSubnet handles GET /api/v1/subnets/{id}
func (g *Gateway) handleGetSubnet(w http.ResponseWriter, r *http.Request) {
	// Extract subnet ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		g.writeErrorResponse(w, http.StatusBadRequest, "MISSING_FIELD", "Subnet ID is required", nil)
		return
	}

	req := &pb.GetSubnetRequest{
		Id: id,
	}

	// Call service layer
	resp, err := g.serviceLayer.GetSubnet(r.Context(), req)
	if err != nil {
		g.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), err)
		return
	}

	// Check for service-level errors
	if resp.Error != nil {
		g.writeProtobufError(w, resp.Error)
		return
	}

	// Convert response to JSON and send
	jsonSubnet := SubnetToJSON(resp.Subnet)
	g.writeJSON(w, http.StatusOK, jsonSubnet)
}

// handleUpdateSubnet handles PUT /api/v1/subnets/{id}
func (g *Gateway) handleUpdateSubnet(w http.ResponseWriter, r *http.Request) {
	// Extract subnet ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		g.writeErrorResponse(w, http.StatusBadRequest, "MISSING_FIELD", "Subnet ID is required", nil)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to read request body", err)
		return
	}
	defer r.Body.Close()

	// Validate request body is not empty
	if len(body) == 0 {
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Request body is required", nil)
		return
	}

	// Convert JSON to Protobuf request
	req, err := JSONToUpdateSubnetRequest(id, body)
	if err != nil {
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_MESSAGE_FORMAT", err.Error(), err)
		return
	}

	// Call service layer
	resp, err := g.serviceLayer.UpdateSubnet(r.Context(), req)
	if err != nil {
		g.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), err)
		return
	}

	// Check for service-level errors
	if resp.Error != nil {
		g.writeProtobufError(w, resp.Error)
		return
	}

	// Convert response to JSON and send
	jsonSubnet := SubnetToJSON(resp.Subnet)
	g.writeJSON(w, http.StatusOK, jsonSubnet)
}

// handleDeleteSubnet handles DELETE /api/v1/subnets/{id}
func (g *Gateway) handleDeleteSubnet(w http.ResponseWriter, r *http.Request) {
	// Extract subnet ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		g.writeErrorResponse(w, http.StatusBadRequest, "MISSING_FIELD", "Subnet ID is required", nil)
		return
	}

	req := &pb.DeleteSubnetRequest{
		Id: id,
	}

	// Call service layer
	resp, err := g.serviceLayer.DeleteSubnet(r.Context(), req)
	if err != nil {
		g.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), err)
		return
	}

	// Check for service-level errors
	if resp.Error != nil {
		g.writeProtobufError(w, resp.Error)
		return
	}

	// Return success response
	g.writeJSON(w, http.StatusOK, &DeleteResponseJSON{Success: resp.Success})
}

// handleGetSubnetChildren handles GET /api/v1/subnets/{id}/children
func (g *Gateway) handleGetSubnetChildren(w http.ResponseWriter, r *http.Request) {
	// Extract subnet ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		g.writeErrorResponse(w, http.StatusBadRequest, "MISSING_FIELD", "Subnet ID is required", nil)
		return
	}

	ctx := r.Context()
	children, err := g.serviceLayer.GetSubnetChildren(ctx, id)
	if err != nil {
		g.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), err)
		return
	}

	// Convert repository models to JSON
	jsonChildren := RepositorySubnetsToJSON(children)

	g.writeJSON(w, http.StatusOK, map[string]interface{}{
		"children": jsonChildren,
		"count":    len(jsonChildren),
	})
}

// handleListSubnetsRepository handles GET /api/v1/subnets using repository models
func (g *Gateway) handleListSubnetsRepository(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	filters := repository.SubnetFilters{
		LocationFilter:      query.Get("location"),
		CloudProviderFilter: query.Get("cloud_provider"),
		SearchQuery:         query.Get("search"),
		Page:                parseIntParam(query.Get("page"), 0),
		PageSize:            parseIntParam(query.Get("page_size"), 50),
	}

	ctx := r.Context()

	// Use repository directly to get enhanced data
	result, err := g.serviceLayer.ListSubnetsRepository(ctx, filters)
	if err != nil {
		g.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), err)
		return
	}

	// Convert repository models to JSON
	jsonSubnets := RepositorySubnetsToJSON(result.Subnets)

	jsonResp := &ListSubnetsResponseJSON{
		Subnets:    jsonSubnets,
		TotalCount: result.TotalCount,
	}
	g.writeJSON(w, http.StatusOK, jsonResp)
}

// handleCreateSubnetRepository handles POST /api/v1/subnets using repository models
func (g *Gateway) handleCreateSubnetRepository(w http.ResponseWriter, r *http.Request) {
	log.Println("[CreateSubnetRepository] Received request")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[CreateSubnetRepository] Failed to read body: %v", err)
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to read request body", err)
		return
	}
	defer r.Body.Close()

	log.Printf("[CreateSubnetRepository] Request body: %s", string(body))

	// Validate request body is not empty
	if len(body) == 0 {
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Request body is required", nil)
		return
	}

	// Parse JSON directly to repository model
	var subnetData struct {
		CIDR         string         `json:"cidr"`
		Name         string         `json:"name"`
		Description  string         `json:"description,omitempty"`
		Location     string         `json:"location,omitempty"`
		LocationType string         `json:"location_type,omitempty"`
		CloudInfo    *CloudInfoJSON `json:"cloud_info,omitempty"`
		ParentID     string         `json:"parent_id,omitempty"`
	}

	if err := json.Unmarshal(body, &subnetData); err != nil {
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_MESSAGE_FORMAT", err.Error(), err)
		return
	}

	// Validate required fields
	if subnetData.CIDR == "" {
		g.writeErrorResponse(w, http.StatusBadRequest, "MISSING_FIELD", "CIDR is required", nil)
		return
	}
	if subnetData.Name == "" {
		g.writeErrorResponse(w, http.StatusBadRequest, "MISSING_FIELD", "Name is required", nil)
		return
	}

	// Create repository subnet model
	subnet := &repository.Subnet{
		ID:           uuid.New().String(),
		Name:         subnetData.Name,
		CIDR:         subnetData.CIDR,
		Location:     subnetData.Location,
		LocationType: subnetData.LocationType,
		ParentID:     subnetData.ParentID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Add cloud info if provided
	if subnetData.CloudInfo != nil {
		subnet.CloudInfo = &repository.CloudInfo{
			Provider:     subnetData.CloudInfo.Provider,
			Region:       subnetData.CloudInfo.Region,
			AccountID:    subnetData.CloudInfo.AccountID,
			ResourceType: subnetData.CloudInfo.ResourceType,
			VPCId:        subnetData.CloudInfo.VPCId,
			SubnetId:     subnetData.CloudInfo.SubnetId,
		}
	}

	log.Printf("[CreateSubnetRepository] Repository model: %+v", subnet)

	// Create subnet using service layer (which will calculate details and create in repository)
	ctx := r.Context()
	err = g.serviceLayer.CreateSubnetRepository(ctx, subnet)
	if err != nil {
		log.Printf("[CreateSubnetRepository] Service layer error: %v", err)
		g.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), err)
		return
	}

	// Retrieve the created subnet with calculated details
	createdSubnet, err := g.serviceLayer.GetSubnetRepository(ctx, subnet.ID)
	if err != nil {
		log.Printf("[CreateSubnetRepository] Failed to retrieve created subnet: %v", err)
		g.writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve created subnet", err)
		return
	}

	log.Printf("[CreateSubnetRepository] Successfully created subnet: %s", subnet.ID)

	// Convert to JSON response
	jsonSubnet := RepositorySubnetToJSON(createdSubnet)
	g.writeJSON(w, http.StatusCreated, jsonSubnet)
}
