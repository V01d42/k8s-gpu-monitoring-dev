# GPU Monitoring API

K8s上でPrometheusからGPUメトリクスを取得・表示するためのGoで実装されたREST APIサーバーです。

## 特徴

- **Go 1.22の新しいServeMux**: 最新のHTTPルーティング機能を使用
- **並行処理**: Goroutineを使用した効率的なPrometheusクエリ
- **RESTful API**: 標準的なHTTP APIデザイン
- **適切なエラーハンドリング**: 包括的なエラー処理とログ
- **CORS対応**: クロスオリジンリクエストをサポート
- **Graceful Shutdown**: 安全なサーバー停止処理

## API エンドポイント

### ヘルスチェック
```
GET /api/health
```
サーバーとPrometheus接続の健全性をチェック

### GPUメトリクス取得
```
GET /api/v1/gpu/metrics
```
全GPUの詳細なメトリクス情報を取得

### GPU搭載ノード一覧
```
GET /api/v1/gpu/nodes
```
GPU搭載ノードの情報を取得

### GPU利用率
```
GET /api/v1/gpu/utilization
```
GPUの利用率のみを取得（軽量）

## 設定

環境変数で設定できます：

- `PROMETHEUS_URL`: PrometheusサーバーのURL（デフォルト: `http://localhost:9090`）
- `PORT`: APIサーバーのポート（デフォルト: `8080`）

## レスポンス形式

すべてのAPIレスポンスは以下の形式です：

```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully",
  "error": null
}
```

## 実行方法

### 開発環境
```bash
# 依存関係を取得
go mod download

# サーバーを起動
go run cmd/server/main.go
```

### 本番環境
```bash
# ビルド
go build -o gpu-monitoring-api cmd/server/main.go

# 実行
./gpu-monitoring-api
```

### Docker
```bash
# イメージをビルド
docker build -t gpu-monitoring-api .

# コンテナを実行
docker run -p 8080:8080 \
  -e PROMETHEUS_URL=http://prometheus:9090 \
  gpu-monitoring-api
```

## テスト

```bash
# すべてのテストを実行
go test ./...

# カバレッジ付きでテスト
go test -cover ./...

# ベンチマークテスト
go test -bench=. ./...
```

## 必要なPrometheusメトリクス

このAPIは以下のNVIDIA GPUメトリクスを期待します：

- `nvidia_smi_utilization_gpu_ratio` - GPU利用率
- `nvidia_smi_memory_used_bytes` - GPU メモリ使用量
- `nvidia_smi_memory_total_bytes` - GPU メモリ総容量
- `nvidia_smi_temperature_gpu_celsius` - GPU 温度
- `nvidia_smi_power_draw_watts` - GPU 電力消費
- `nvidia_smi_enforced_power_limit_watts` - GPU 電力制限
- `nvidia_smi_gpu_info` - GPU 情報

## アーキテクチャ

```
┌─────────────────┐
│   HTTP Client   │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│   Middleware    │ ← CORS, Logging, Recovery
│   (Chain)       │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│   GPUHandler    │ ← HTTP Request Handlers
└─────────────────┘
         │
         ▼
┌─────────────────┐
│ PrometheusClient│ ← Concurrent Queries
└─────────────────┘
         │
         ▼
┌─────────────────┐
│  Prometheus     │
│    Server       │
└─────────────────┘
```

## パフォーマンス

- **並行クエリ**: 複数のPrometheusクエリを並行実行
- **コネクションプール**: HTTPクライアントの再利用
- **タイムアウト**: 適切なタイムアウト設定
- **メモリ効率**: 効率的なデータ構造使用

## セキュリティ

- **入力検証**: 適切な入力バリデーション
- **エラー情報**: 機密情報を含まないエラーメッセージ
- **リソース制限**: 適切なタイムアウトとリソース制限
- **CORS設定**: セキュアなクロスオリジン設定 