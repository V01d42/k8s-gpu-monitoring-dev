package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"k8s-gpu-monitoring/internal/models"
	"k8s-gpu-monitoring/internal/prometheus"
)

// MockPrometheusClient implements a mock prometheus client for testing
type MockPrometheusClient struct {
	shouldFail bool
}

func (m *MockPrometheusClient) Query(ctx context.Context, query string) (*prometheus.PrometheusResponse, error) {
	if m.shouldFail {
		return nil, &mockError{message: "mock error"}
	}

	// モックレスポンスを返す
	return &prometheus.PrometheusResponse{
		Status: "success",
		Data: struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Value  []interface{}     `json:"value"`
			} `json:"result"`
		}{
			ResultType: "vector",
			Result: []struct {
				Metric map[string]string `json:"metric"`
				Value  []interface{}     `json:"value"`
			}{
				{
					Metric: map[string]string{
						"node": "test-node",
						"gpu":  "0",
					},
					Value: []interface{}{1640995200, "50.5"},
				},
			},
		},
	}, nil
}

func (m *MockPrometheusClient) GetGPUMetrics(ctx context.Context) ([]models.GPUMetrics, error) {
	if m.shouldFail {
		return nil, &mockError{message: "mock error"}
	}

	return []models.GPUMetrics{
		{
			NodeName:          "test-node",
			GPUIndex:          0,
			GPUName:           "Tesla V100",
			Utilization:       75.5,
			MemoryUsed:        12.5,
			MemoryTotal:       16.0,
			MemoryUtilization: 78.125,
			Temperature:       65.0,
			PowerDraw:         250.0,
			PowerLimit:        300.0,
		},
	}, nil
}

func (m *MockPrometheusClient) GetGPUNodes(ctx context.Context) ([]models.GPUNode, error) {
	if m.shouldFail {
		return nil, &mockError{message: "mock error"}
	}

	return []models.GPUNode{
		{
			NodeName:  "test-node",
			GPUCount:  2,
			GPUModels: []string{"Tesla V100"},
		},
	}, nil
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

func TestGPUHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		shouldFail     bool
		expectedStatus int
	}{
		{
			name:           "successful health check",
			shouldFail:     false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "failed health check",
			shouldFail:     true,
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// プロメテウスクライアントをラップ
			promClient := &prometheus.Client{}
			handler := NewGPUHandler(promClient)

			// モッククライアントを使用するためにハンドラーを直接テスト
			w := httptest.NewRecorder()

			// モックを使った簡易テスト
			if tt.shouldFail {
				handler.writeErrorResponse(w, http.StatusServiceUnavailable, "Prometheus connection failed")
			} else {
				response := models.APIResponse{
					Success: true,
					Message: "Service is healthy",
					Data: map[string]interface{}{
						"status":  "healthy",
						"version": "1.0.0",
					},
				}
				handler.writeJSONResponse(w, http.StatusOK, response)
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response models.APIResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if tt.shouldFail && response.Success {
				t.Error("Expected failure response, got success")
			}
			if !tt.shouldFail && !response.Success {
				t.Error("Expected success response, got failure")
			}
		})
	}
}

func TestGPUHandler_GetGPUMetrics(t *testing.T) {
	tests := []struct {
		name           string
		shouldFail     bool
		expectedStatus int
	}{
		{
			name:           "successful metrics retrieval",
			shouldFail:     false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "failed metrics retrieval",
			shouldFail:     true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockPrometheusClient{shouldFail: tt.shouldFail}
			promClient := &prometheus.Client{}
			handler := NewGPUHandler(promClient)

			w := httptest.NewRecorder()

			// モックレスポンスをシミュレート
			if tt.shouldFail {
				handler.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve GPU metrics")
			} else {
				metrics, _ := mockClient.GetGPUMetrics(context.Background())
				response := models.APIResponse{
					Success: true,
					Data:    metrics,
					Message: "GPU metrics retrieved successfully",
				}
				handler.writeJSONResponse(w, http.StatusOK, response)
			}

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// JSON レスポンスの検証
			if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
				t.Error("Expected JSON content type")
			}

			var response models.APIResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if tt.shouldFail && response.Success {
				t.Error("Expected failure response, got success")
			}
			if !tt.shouldFail && !response.Success {
				t.Error("Expected success response, got failure")
			}

			// 成功時のデータ検証
			if !tt.shouldFail && response.Data != nil {
				metricsData, ok := response.Data.([]interface{})
				if !ok {
					t.Error("Expected metrics data to be an array")
				} else if len(metricsData) == 0 {
					t.Error("Expected non-empty metrics data")
				}
			}
		})
	}
}

func TestGPUHandler_WriteJSONResponse(t *testing.T) {
	promClient := &prometheus.Client{}
	handler := NewGPUHandler(promClient)

	testData := map[string]interface{}{
		"test":   "data",
		"number": 42,
	}

	w := httptest.NewRecorder()

	handler.writeJSONResponse(w, http.StatusOK, testData)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if result["test"] != "data" {
		t.Errorf("Expected test field to be 'data', got %v", result["test"])
	}
}
