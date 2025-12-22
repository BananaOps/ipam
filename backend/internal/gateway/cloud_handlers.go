package gateway

import (
	"encoding/json"
	"net/http"
)

// CloudSyncRequest represents a cloud sync request
type CloudSyncRequest struct {
	Provider string `json:"provider,omitempty"`
	Region   string `json:"region,omitempty"`
}

// CloudSyncResponse represents a cloud sync response
type CloudSyncResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CloudStatusResponse represents cloud provider status
type CloudStatusResponse struct {
	Enabled   bool                    `json:"enabled"`
	Providers map[string]ProviderInfo `json:"providers"`
}

// ProviderInfo represents cloud provider information
type ProviderInfo struct {
	Enabled bool     `json:"enabled"`
	Regions []string `json:"regions"`
}

// HandleCloudSync handles cloud synchronization requests
func (g *Gateway) HandleCloudSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CloudSyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		g.writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err)
		return
	}

	ctx := r.Context()

	// Check if cloud providers are enabled
	if !g.cloudManager.IsEnabled() {
		g.writeErrorResponse(w, http.StatusServiceUnavailable, "CLOUD_DISABLED", "Cloud providers are disabled", nil)
		return
	}

	var err error
	var message string

	switch req.Provider {
	case "aws":
		if req.Region != "" {
			err = g.cloudManager.SyncAWSRegion(ctx, req.Region)
			message = "AWS region " + req.Region + " synchronized successfully"
		} else {
			err = g.cloudManager.SyncAll(ctx)
			message = "All AWS regions synchronized successfully"
		}
	case "":
		// Sync all providers
		err = g.cloudManager.SyncAll(ctx)
		message = "All cloud providers synchronized successfully"
	default:
		g.writeErrorResponse(w, http.StatusBadRequest, "UNSUPPORTED_PROVIDER", "Unsupported cloud provider: "+req.Provider, nil)
		return
	}

	if err != nil {
		g.writeErrorResponse(w, http.StatusInternalServerError, "SYNC_FAILED", "Cloud synchronization failed", err)
		return
	}

	response := CloudSyncResponse{
		Success: true,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleCloudStatus handles cloud provider status requests
func (g *Gateway) HandleCloudStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providers := make(map[string]ProviderInfo)

	// AWS status
	if g.cloudManager.IsAWSEnabled() {
		providers["aws"] = ProviderInfo{
			Enabled: true,
			Regions: g.cloudManager.ListAWSRegions(),
		}
	} else {
		providers["aws"] = ProviderInfo{
			Enabled: false,
			Regions: []string{},
		}
	}

	response := CloudStatusResponse{
		Enabled:   g.cloudManager.IsEnabled(),
		Providers: providers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleUpdateUtilization handles utilization update requests
func (g *Gateway) HandleUpdateUtilization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// Check if cloud providers are enabled
	if !g.cloudManager.IsEnabled() {
		g.writeErrorResponse(w, http.StatusServiceUnavailable, "CLOUD_DISABLED", "Cloud providers are disabled", nil)
		return
	}

	err := g.cloudManager.UpdateUtilization(ctx)
	if err != nil {
		g.writeErrorResponse(w, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update utilization", err)
		return
	}

	response := CloudSyncResponse{
		Success: true,
		Message: "Utilization data updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
