// GPU メトリクスのデータ型定義
export interface GPUMetrics {
  node_name: string;
  gpu_index: number;
  gpu_name: string;
  utilization: number;
  memory_used: number;
  memory_total: number;
  memory_utilization: number;
  temperature: number;
  power_draw: number;
  power_limit: number;
  timestamp: string;
}

// GPU ノード情報の型定義
export interface GPUNode {
  node_name: string;
  gpu_count: number;
  gpu_models: string[];
}

// API レスポンスの型定義
export interface APIResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

// GPU 利用率の型定義
export interface GPUUtilization {
  node: string;
  gpu_index: string;
  utilization: string;
  timestamp: number;
}

// テーブルの設定タイプ
export interface TableConfig {
  pageSize: number;
  refreshInterval: number;
  autoRefresh: boolean;
}

// フィルタリングの設定
export interface FilterConfig {
  nodeFilter: string;
  utilizationRange: [number, number];
  temperatureRange: [number, number];
  showOnlyHighUtilization: boolean;
} 