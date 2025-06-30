package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"k8s-gpu-monitoring/internal/models"
	"k8s-gpu-monitoring/internal/prometheus"
)

// GPUHandler handles GPU-related HTTP requests
type GPUHandler struct {
	promClient *prometheus.Client
}

// NewGPUHandler creates a new GPU handler
func NewGPUHandler(promClient *prometheus.Client) *GPUHandler {
	return &GPUHandler{
		promClient: promClient,
	}
}

// writeJSONResponse writes a JSON response with proper headers
func (h *GPUHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// writeErrorResponse writes an error response
func (h *GPUHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := models.APIResponse{
		Success: false,
		Error:   message,
	}
	h.writeJSONResponse(w, statusCode, response)
}

// GetGPUMetrics handles GET /api/v1/gpu/metrics
func (h *GPUHandler) GetGPUMetrics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	metrics, err := h.promClient.GetGPUMetrics(ctx)
	if err != nil {
		log.Printf("Error getting GPU metrics: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve GPU metrics")
		return
	}

	response := models.APIResponse{
		Success: true,
		Data:    metrics,
		Message: "GPU metrics retrieved successfully",
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetGPUNodes handles GET /api/v1/gpu/nodes
func (h *GPUHandler) GetGPUNodes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	nodes, err := h.promClient.GetGPUNodes(ctx)
	if err != nil {
		log.Printf("Error getting GPU nodes: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve GPU nodes")
		return
	}

	response := models.APIResponse{
		Success: true,
		Data:    nodes,
		Message: "GPU nodes retrieved successfully",
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// GetGPUUtilization handles GET /api/v1/gpu/utilization
func (h *GPUHandler) GetGPUUtilization(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// 利用率のみを取得するクエリ
	query := `nvidia_smi_utilization_gpu_ratio * 100`

	resp, err := h.promClient.Query(ctx, query)
	if err != nil {
		log.Printf("Error getting GPU utilization: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve GPU utilization")
		return
	}

	// レスポンスを簡素化
	var utilization []map[string]interface{}
	for _, result := range resp.Data.Result {
		if len(result.Value) >= 2 {
			util := map[string]interface{}{
				"node":        result.Metric["node"],
				"gpu_index":   result.Metric["gpu"],
				"utilization": result.Value[1],
				"timestamp":   result.Value[0],
			}
			utilization = append(utilization, util)
		}
	}

	response := models.APIResponse{
		Success: true,
		Data:    utilization,
		Message: "GPU utilization retrieved successfully",
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// HealthCheck handles GET /api/health
func (h *GPUHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Prometheusサーバーの接続性をチェック
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.promClient.Query(ctx, "up")
	if err != nil {
		log.Printf("Health check failed: %v", err)
		h.writeErrorResponse(w, http.StatusServiceUnavailable, "Prometheus connection failed")
		return
	}

	response := models.APIResponse{
		Success: true,
		Message: "Service is healthy",
		Data: map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		},
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}
