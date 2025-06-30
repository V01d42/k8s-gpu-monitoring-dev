package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"k8s-gpu-monitoring/internal/models"
)

// Client represents Prometheus client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// PrometheusResponse represents Prometheus API response
type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
	Error     string `json:"error,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
}

// NewClient creates a new Prometheus client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Query executes a PromQL query
func (c *Client) Query(ctx context.Context, query string) (*PrometheusResponse, error) {
	queryURL := fmt.Sprintf("%s/api/v1/query", c.baseURL)

	params := url.Values{}
	params.Add("query", query)
	params.Add("time", strconv.FormatInt(time.Now().Unix(), 10))

	req, err := http.NewRequestWithContext(ctx, "GET", queryURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var promResp PrometheusResponse
	if err := json.Unmarshal(body, &promResp); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	if promResp.Status != "success" {
		return nil, fmt.Errorf("prometheus query failed: %s - %s", promResp.ErrorType, promResp.Error)
	}

	return &promResp, nil
}

// GetGPUMetrics retrieves GPU metrics from Prometheus
func (c *Client) GetGPUMetrics(ctx context.Context) ([]models.GPUMetrics, error) {
	// クエリ群を並行実行
	queries := map[string]string{
		"utilization":  `nvidia_smi_utilization_gpu_ratio`,
		"memory_used":  `nvidia_smi_memory_used_bytes`,
		"memory_total": `nvidia_smi_memory_total_bytes`,
		"temperature":  `nvidia_smi_temperature_gpu_celsius`,
		"power_draw":   `nvidia_smi_power_draw_watts`,
		"power_limit":  `nvidia_smi_enforced_power_limit_watts`,
		"gpu_name":     `nvidia_smi_gpu_info`,
	}

	results := make(map[string]*PrometheusResponse)
	errors := make(chan error, len(queries))

	// 並行してクエリを実行
	for name, query := range queries {
		go func(name, query string) {
			resp, err := c.Query(ctx, query)
			if err != nil {
				errors <- fmt.Errorf("query %s failed: %w", name, err)
				return
			}
			results[name] = resp
			errors <- nil
		}(name, query)
	}

	// すべてのクエリの完了を待機
	for i := 0; i < len(queries); i++ {
		if err := <-errors; err != nil {
			return nil, err
		}
	}

	return c.parseGPUMetrics(results)
}

// parseGPUMetrics parses Prometheus response into GPUMetrics
func (c *Client) parseGPUMetrics(results map[string]*PrometheusResponse) ([]models.GPUMetrics, error) {
	// メトリクスをノード・GPU単位でグループ化
	metricsMap := make(map[string]*models.GPUMetrics) // key: "node_name:gpu_index"

	for metricType, response := range results {
		for _, result := range response.Data.Result {
			nodeName := result.Metric["node"]
			gpuIndex := result.Metric["gpu"]

			if nodeName == "" || gpuIndex == "" {
				continue
			}

			key := fmt.Sprintf("%s:%s", nodeName, gpuIndex)

			if metricsMap[key] == nil {
				idx, _ := strconv.Atoi(gpuIndex)
				metricsMap[key] = &models.GPUMetrics{
					NodeName:  nodeName,
					GPUIndex:  idx,
					Timestamp: time.Now(),
				}
			}

			// 値を取得してパース
			if len(result.Value) >= 2 {
				valueStr, ok := result.Value[1].(string)
				if !ok {
					continue
				}

				value, err := strconv.ParseFloat(valueStr, 64)
				if err != nil {
					continue
				}

				// メトリクスタイプに応じて値を設定
				switch metricType {
				case "utilization":
					metricsMap[key].Utilization = value * 100 // 0-1 を 0-100% に変換
				case "memory_used":
					metricsMap[key].MemoryUsed = value / (1024 * 1024 * 1024) // bytes to GB
				case "memory_total":
					metricsMap[key].MemoryTotal = value / (1024 * 1024 * 1024) // bytes to GB
				case "temperature":
					metricsMap[key].Temperature = value
				case "power_draw":
					metricsMap[key].PowerDraw = value
				case "power_limit":
					metricsMap[key].PowerLimit = value
				case "gpu_name":
					if gpuName, exists := result.Metric["name"]; exists {
						metricsMap[key].GPUName = gpuName
					}
				}
			}
		}
	}

	// メモリ使用率を計算
	for _, metrics := range metricsMap {
		if metrics.MemoryTotal > 0 {
			metrics.MemoryUtilization = (metrics.MemoryUsed / metrics.MemoryTotal) * 100
		}
	}

	// スライスに変換
	var gpuMetrics []models.GPUMetrics
	for _, metrics := range metricsMap {
		gpuMetrics = append(gpuMetrics, *metrics)
	}

	return gpuMetrics, nil
}

// GetGPUNodes retrieves GPU node information
func (c *Client) GetGPUNodes(ctx context.Context) ([]models.GPUNode, error) {
	query := `group by (node) (nvidia_smi_gpu_info)`

	resp, err := c.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("getting GPU nodes: %w", err)
	}

	var nodes []models.GPUNode
	nodeMap := make(map[string]*models.GPUNode)

	for _, result := range resp.Data.Result {
		nodeName := result.Metric["node"]
		if nodeName == "" {
			continue
		}

		if nodeMap[nodeName] == nil {
			nodeMap[nodeName] = &models.GPUNode{
				NodeName:  nodeName,
				GPUModels: make([]string, 0),
			}
		}

		// GPU数を取得するための別クエリ
		countQuery := fmt.Sprintf(`count by (node) (nvidia_smi_gpu_info{node="%s"})`, nodeName)
		countResp, err := c.Query(ctx, countQuery)
		if err == nil && len(countResp.Data.Result) > 0 && len(countResp.Data.Result[0].Value) >= 2 {
			if countStr, ok := countResp.Data.Result[0].Value[1].(string); ok {
				if count, err := strconv.Atoi(countStr); err == nil {
					nodeMap[nodeName].GPUCount = count
				}
			}
		}
	}

	for _, node := range nodeMap {
		nodes = append(nodes, *node)
	}

	return nodes, nil
}
