import { useState, useEffect } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { RefreshCw, AlertCircle, Activity } from 'lucide-react';

import { gpuApi, queryKeys } from '@/api/client';
import { GPUTable } from '@/components/GPUTable';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { cn } from '@/lib/utils';

function App() {
  const [autoRefresh, setAutoRefresh] = useState(true);
  const queryClient = useQueryClient();

  // GPUメトリクスデータを取得
  const {
    data: gpuData,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: queryKeys.gpuMetrics,
    queryFn: gpuApi.getGPUMetrics,
    refetchInterval: autoRefresh ? 30000 : false, // 30秒間隔で自動更新
    retry: 3,
    retryDelay: 1000,
  });

  // ヘルスチェック
  const { data: healthData } = useQuery({
    queryKey: queryKeys.health,
    queryFn: gpuApi.checkHealth,
    refetchInterval: 60000, // 1分間隔
    retry: 1,
  });

  // 手動更新
  const handleRefresh = () => {
    refetch();
    queryClient.invalidateQueries({ queryKey: queryKeys.gpuMetrics });
  };

  // 統計情報を計算
  const stats = gpuData?.data ? {
    totalGPUs: gpuData.data.length,
    activeGPUs: gpuData.data.filter(gpu => gpu.utilization > 5).length,
    averageUtilization: gpuData.data.reduce((acc, gpu) => acc + gpu.utilization, 0) / gpuData.data.length,
    highTempGPUs: gpuData.data.filter(gpu => gpu.temperature > 80).length,
  } : null;

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto p-6 space-y-6">
        {/* ヘッダー */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">GPU監視ダッシュボード</h1>
            <p className="text-muted-foreground">
              Kubernetes上のGPUリソースをリアルタイムで監視
            </p>
          </div>
          
          <div className="flex items-center space-x-4">
            {/* ヘルスステータス */}
            <div className="flex items-center space-x-2">
              {healthData?.success ? (
                <div className="flex items-center space-x-2 text-green-600">
                  <Activity className="h-4 w-4" />
                  <span className="text-sm">接続OK</span>
                </div>
              ) : (
                <div className="flex items-center space-x-2 text-red-600">
                  <AlertCircle className="h-4 w-4" />
                  <span className="text-sm">接続エラー</span>
                </div>
              )}
            </div>

            {/* 自動更新トグル */}
            <Button
              variant={autoRefresh ? "default" : "outline"}
              size="sm"
              onClick={() => setAutoRefresh(!autoRefresh)}
            >
              {autoRefresh ? "自動更新ON" : "自動更新OFF"}
            </Button>

            {/* 手動更新ボタン */}
            <Button
              variant="outline"
              size="sm"
              onClick={handleRefresh}
              disabled={isLoading}
            >
              <RefreshCw className={cn("h-4 w-4 mr-2", isLoading && "animate-spin")} />
              更新
            </Button>
          </div>
        </div>

        {/* 統計カード */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">総GPU数</CardTitle>
                <Activity className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.totalGPUs}</div>
                <p className="text-xs text-muted-foreground">
                  使用中: {stats.activeGPUs}
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">平均利用率</CardTitle>
                <Activity className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {stats.averageUtilization.toFixed(1)}%
                </div>
                <p className="text-xs text-muted-foreground">
                  全GPU平均
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">高温アラート</CardTitle>
                <AlertCircle className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className={cn(
                  "text-2xl font-bold",
                  stats.highTempGPUs > 0 ? "text-red-500" : "text-green-500"
                )}>
                  {stats.highTempGPUs}
                </div>
                <p className="text-xs text-muted-foreground">
                  80°C以上
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">アクティブ率</CardTitle>
                <Activity className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {stats.totalGPUs > 0 ? ((stats.activeGPUs / stats.totalGPUs) * 100).toFixed(1) : 0}%
                </div>
                <p className="text-xs text-muted-foreground">
                  利用中のGPU
                </p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* GPUテーブル */}
        <GPUTable
          data={gpuData?.data || []}
          isLoading={isLoading}
          error={error}
        />
      </div>
    </div>
  );
}

export default App; 