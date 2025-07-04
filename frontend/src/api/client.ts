import axios from 'axios';
import { APIResponse, GPUMetrics, GPUNode, GPUUtilization } from '@/types/gpu';

// APIベースURL（環境変数または開発時のデフォルト値）
const API_BASE_URL = import.meta.env.VITE_API_URL || '/api';

// Axiosインスタンスを作成
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// レスポンスインターセプター（エラーハンドリング）
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error);
    
    if (error.code === 'ECONNABORTED') {
      throw new Error('Request timeout');
    }
    
    if (!error.response) {
      throw new Error('Network error');
    }
    
    const status = error.response.status;
    switch (status) {
      case 404:
        throw new Error('API endpoint not found');
      case 500:
        throw new Error('Server error');
      case 503:
        throw new Error('Service unavailable');
      default:
        throw new Error(`HTTP error: ${status}`);
    }
  }
);

// API関数群
export const gpuApi = {
  // ヘルスチェック
  async checkHealth(): Promise<APIResponse<any>> {
    const response = await apiClient.get<APIResponse<any>>('/health');
    return response.data;
  },

  // 全GPUメトリクスを取得
  async getGPUMetrics(): Promise<APIResponse<GPUMetrics[]>> {
    const response = await apiClient.get<APIResponse<GPUMetrics[]>>('/v1/gpu/metrics');
    return response.data;
  },

  // GPU搭載ノード一覧を取得
  async getGPUNodes(): Promise<APIResponse<GPUNode[]>> {
    const response = await apiClient.get<APIResponse<GPUNode[]>>('/v1/gpu/nodes');
    return response.data;
  },

  // GPU利用率のみを取得（軽量）
  async getGPUUtilization(): Promise<APIResponse<GPUUtilization[]>> {
    const response = await apiClient.get<APIResponse<GPUUtilization[]>>('/v1/gpu/utilization');
    return response.data;
  },
};

// TanStack Query用のキー
export const queryKeys = {
  health: ['health'] as const,
  gpuMetrics: ['gpu', 'metrics'] as const,
  gpuNodes: ['gpu', 'nodes'] as const,
  gpuUtilization: ['gpu', 'utilization'] as const,
} as const; 