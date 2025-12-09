package gateway

import (
	"fmt"
	"io"
	"log"
	"net/http"

	pb "github.com/bananaops/ipam-bananaops/proto"
	"github.com/gorilla/mux"
)

// handleCreateSubnet handles POST /api/v1/subnets
func (g *RESTGateway) handleCreateSubnet(w http.ResponseWriter, r *http.Request) {
	log.Println("[CreateSubnet] Received request")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[CreateSubnet] Failed to read body: %v", err)
		g.writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to read request body")
		return
	}
	defer r.Body.Close()

	log.Printf("[CreateSubnet] Request body: %s", string(body))

	// Validate request body is not empty
	if len(body) == 0 {
		g.writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Request body is required")
		return
	}

	// Convert JSON to Protobuf request
	req, err := JSONToCreateSubnetRequest(body)
	if err != nil {
		g.writeError(w, http.StatusBadRequest, "INVALID_MESSAGE_FORMAT", err.Error())
		return
	}

	// Validate required fields
	if req.Cidr == "" {
		g.writeError(w, http.StatusBadRequest, "MISSING_FIELD", "CIDR is required")
		return
	}
	if req.Name == "" {
		g.writeError(w, http.StatusBadRequest, "MISSING_FIELD", "Name is required")
		return
	}

	// Call service layer
	resp, err := g.serviceLayer.CreateSubnet(r.Context(), req)
	if err != nil {
		g.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Check for service-level errors
	if resp.Error != nil {
		g.writeProtobufError(w, resp.Error)
		return
	}

	// Convert response to JSON and send
	jsonSubnet := SubnetToJSON(resp.Subnet)
	g.writeJSON(w, http.StatusCreated, jsonSubnet)
}

// handleListSubnets handles GET /api/v1/subnets
func (g *RESTGateway) handleListSubnets(w http.ResponseWriter, r *http.Request) {
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
		g.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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
func (g *RESTGateway) handleGetSubnet(w http.ResponseWriter, r *http.Request) {
	// Extract subnet ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		g.writeError(w, http.StatusBadRequest, "MISSING_FIELD", "Subnet ID is required")
		return
	}

	req := &pb.GetSubnetRequest{
		Id: id,
	}

	// Call service layer
	resp, err := g.serviceLayer.GetSubnet(r.Context(), req)
	if err != nil {
		g.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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
func (g *RESTGateway) handleUpdateSubnet(w http.ResponseWriter, r *http.Request) {
	// Extract subnet ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		g.writeError(w, http.StatusBadRequest, "MISSING_FIELD", "Subnet ID is required")
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		g.writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to read request body")
		return
	}
	defer r.Body.Close()

	// Validate request body is not empty
	if len(body) == 0 {
		g.writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Request body is required")
		return
	}

	// Convert JSON to Protobuf request
	req, err := JSONToUpdateSubnetRequest(id, body)
	if err != nil {
		g.writeError(w, http.StatusBadRequest, "INVALID_MESSAGE_FORMAT", err.Error())
		return
	}

	// Call service layer
	resp, err := g.serviceLayer.UpdateSubnet(r.Context(), req)
	if err != nil {
		g.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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
func (g *RESTGateway) handleDeleteSubnet(w http.ResponseWriter, r *http.Request) {
	// Extract subnet ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		g.writeError(w, http.StatusBadRequest, "MISSING_FIELD", "Subnet ID is required")
		return
	}

	req := &pb.DeleteSubnetRequest{
		Id: id,
	}

	// Call service layer
	resp, err := g.serviceLayer.DeleteSubnet(r.Context(), req)
	if err != nil {
		g.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
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

// parseIntParam parses an integer query parameter with a default value
func parseIntParam(s string, defaultVal int32) int32 {
	if s == "" {
		return defaultVal
	}
	var val int32
	_, err := fmt.Sscanf(s, "%d", &val)
	if err != nil {
		return defaultVal
	}
	return val
}
