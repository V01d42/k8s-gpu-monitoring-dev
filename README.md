# K8s GPU Monitoring Dashboard

Kubernetes上でPrometheusからGPUメトリクスを取得・表示するための統合監視ダッシュボードです。Go製のバックエンドAPIとReact製のフロントエンドで構成されています。

![GPU Dashboard](https://img.shields.io/badge/Status-Production%20Ready-green)
![Go](https://img.shields.io/badge/Go-1.24-blue)
![React](https://img.shields.io/badge/React-19-blue)
![TypeScript](https://img.shields.io/badge/TypeScript-5.7-blue)
![Vite](https://img.shields.io/badge/Vite-7.0-purple)

## 特徴

### 最新技術スタック
- **Go 1.24**: 最新の標準ライブラリServeMuxとメソッド指定ルーティングを使用
- **React 19.1.0**: 最新のReactによる高性能フロントエンド
- **TypeScript 5.7**: 型安全な開発環境
- **Vite 7.0**: 超高速ビルドツール
- **TailwindCSS 4.1**: 最新のユーティリティファーストCSS

### 高性能アーキテクチャ
- **並行処理**: Goroutineによる効率的なPrometheusクエリの並列実行
- **TanStack Query**: 高度なデータフェッチング・キャッシング・状態管理
- **TanStack Table**: 大量データを高速表示する高性能テーブル
- **リアルタイム更新**: 30秒間隔での自動データ更新と手動更新機能

### 豊富な監視機能
- **詳細メトリクス**: GPU利用率、メモリ使用量、温度、電力消費
- **統計ダッシュボード**: 総GPU数、平均利用率、アクティブ率、高温アラート
- **ビジュアル表示**: プログレスバーと色分けによる直感的な表示
- **ヘルスモニタリング**: システム全体の健全性監視とPromtheus接続状態

### モダンUI/UX
- **レスポンシブデザイン**: デスクトップ・モバイル完全対応
- **ダークモード**: 目に優しい表示（デフォルト有効）
- **アクセシビリティ**: WCAG準拠のキーボードナビゲーション
- **リアルタイム表示**: ソート機能付きテーブルでの動的データ表示

## アーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                       │
│                                                             │
│  ┌─────────────┐    ┌──────────────┐    ┌────────────────┐  │
│  │   Ingress   │    │   Frontend   │    │   Backend      │  │
│  │ Controller  │    │  React 19    │    │   Go 1.24     │  │
│  │             │───→│  TypeScript  │───→│  ServeMux     │──┼─→ Prometheus
│  │   nginx     │    │  TanStack    │    │  + CORS       │  │
│  │             │    │  Tailwind    │    │  + Recovery   │  │
│  └─────────────┘    └──────────────┘    └────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## クイックスタート

### Helmでのインストール

```bash
# Helmリポジトリ追加
helm repo add gpu-monitoring https://v01d42.github.io/k8s-gpu-monitoring-dev
helm repo update

# 基本インストール
helm install gpu-monitoring gpu-monitoring/k8s-gpu-monitoring-dev \
  --namespace gpu-monitoring \
  --create-namespace \
  --set backend.env.PROMETHEUS_URL=http://prometheus-server:9090 \
  --set ingress.hosts[0].host=gpu-monitoring.local

# Ingress IPを確認してhostsファイルに追加
kubectl get ingress -n gpu-monitoring
echo "192.168.1.100 gpu-monitoring.local" | sudo tee -a /etc/hosts
```

### アクセス

#### Web UI
```bash
# Ingressアクセス
http://gpu-monitoring.local

# Port-forwardでローカルアクセス
kubectl port-forward -n gpu-monitoring svc/gpu-monitoring-frontend 3000:80
# http://localhost:3000 でアクセス
```

#### API
```bash
# ヘルスチェック
curl http://gpu-monitoring.local/api/health

# GPUメトリクス取得
curl http://gpu-monitoring.local/api/v1/gpu/metrics

# GPU利用率のみ取得（軽量）
curl http://gpu-monitoring.local/api/v1/gpu/utilization

# GPUノード一覧
curl http://gpu-monitoring.local/api/v1/gpu/nodes
```

## プロジェクト構成

```
k8s-gpu-monitoring-dev/
├── backend/                      # Go 1.24 API サーバー
│   ├── cmd/server/main.go       # メインエントリーポイント
│   ├── internal/
│   │   ├── handlers/            # HTTPハンドラー（GPU関連API）
│   │   │   ├── gpu.go          # メインハンドラー
│   │   │   └── gpu_test.go     # 単体テスト
│   │   ├── middleware/          # HTTP ミドルウェア
│   │   │   └── middleware.go   # CORS・ログ・リカバリ
│   │   ├── models/             # データモデル定義
│   │   │   └── gpu.go          # GPU関連構造体
│   │   └── prometheus/         # Prometheusクライアント
│   │       └── client.go       # HTTP API クライアント
│   ├── go.mod                   # Go 1.24 モジュール設定
│   └── Dockerfile               # マルチステージビルド
├── frontend/                     # React 19 + TypeScript フロントエンド
│   ├── src/
│   │   ├── components/          # Reactコンポーネント
│   │   │   ├── GPUTable.tsx    # GPUテーブルコンポーネント
│   │   │   └── ui/             # 共通UIコンポーネント
│   │   │       ├── button.tsx  # ボタンコンポーネント
│   │   │       └── card.tsx    # カードコンポーネント
│   │   ├── api/                # APIクライアント
│   │   │   └── client.ts       # Axiosベースクライアント
│   │   ├── types/              # TypeScript型定義
│   │   │   └── gpu.ts          # GPU関連型定義
│   │   ├── lib/                # ユーティリティ
│   │   │   └── utils.ts        # 共通ユーティリティ関数
│   │   ├── App.tsx             # メインアプリケーション
│   │   ├── main.tsx            # エントリーポイント
│   │   └── index.css           # TailwindCSS設定
│   ├── package.json            # Node.js依存関係
│   ├── vite.config.ts          # Vite 7.0設定
│   ├── tailwind.config.js      # TailwindCSS 4.1設定
│   ├── nginx.conf              # 本番用Nginx設定
│   └── Dockerfile              # マルチステージビルド
├── charts/                     # Helm Charts
│   └── k8s-gpu-monitoring-dev/ # Helm Chart
│       ├── Chart.yaml          # チャート定義
│       ├── values.yaml         # デフォルト設定値
│       └── templates/          # Kubernetesマニフェスト
│           ├── _helpers.tpl    # Helmヘルパー
│           ├── backend/        # バックエンドリソース
│           ├── frontend/       # フロントエンドリソース
│           └── ingress.yaml    # Ingress設定
├── scripts/                    # 運用スクリプト
│   └── release.sh             # リリースプロセス自動化
├── .github/workflows/          # CI/CD
│   └── release.yml            # GitHub Actions ワークフロー
└── docs/                      # ドキュメント
    └── DEPLOYMENT.md          # デプロイメントガイド
```

## 開発環境

### 前提条件
- **Go 1.24以上**
- **Node.js 20以上** (mise.tomlでlatestを指定)
- **Docker & Docker Compose**

### ローカル開発環境

#### 1. バックエンドAPI
```bash
cd backend
go mod download
go run cmd/server/main.go
# サーバーが http://localhost:8080 で起動
```

#### 2. フロントエンド
```bash
cd frontend
npm ci
npm run dev
# Vite開発サーバーが http://localhost:3000 で起動
```

#### 3. 開発時のアクセス
- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8080
- **APIドキュメント**: http://localhost:8080/api/health

## 設定

### 環境変数

#### バックエンド
```bash
PROMETHEUS_URL=http://prometheus-server:9090  # PrometheusサーバーURL
PORT=8080                                     # APIサーバーポート
```

#### フロントエンド
```bash
VITE_API_URL=/api                            # APIベースURL（本番）
# 開発時は vite.config.ts のプロキシ設定を使用
```

### 必要なPrometheusメトリクス

以下のNVIDIA GPUメトリクス（nvidia-gpu-exporter対応）が必要です：

```promql
# GPU利用率
nvidia_gpu_utilization_percent

# メモリ関連
nvidia_gpu_used_memory_bytes
nvidia_gpu_total_memory_bytes
nvidia_gpu_free_memory_bytes
nvidia_gpu_memory_utilization_percent

# 温度
nvidia_gpu_temperature_celsius
```

## カスタマイズ

### Helm設定

#### 基本設定
```yaml
# values.yaml
global:
  imageRegistry: "ghcr.io/v01d42/k8s-gpu-monitoring-dev"

backend:
  enabled: true
  replicas: 1
  resources:
    requests:
      cpu: "250m"
      memory: "256Mi"
    limits:
      cpu: "500m"
      memory: "512Mi"
  env:
    PROMETHEUS_URL: "http://prometheus-server:9090"
    
frontend:
  enabled: true
  replicas: 1
  resources:
    requests:
      cpu: "100m"
      memory: "128Mi"
    limits:
      cpu: "200m"
      memory: "256Mi"
  
ingress:
  enabled: true
  className: "nginx"
  hosts:
    - host: gpu-monitoring.local
      paths:
        - path: /api
          pathType: Prefix
          backend:
            service: backend
            port: 8080
        - path: /
          pathType: Prefix
          backend:
            service: frontend
            port: 80
```

#### セキュリティ設定
```yaml
backend:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1001
    runAsGroup: 1001
    capabilities:
      drop:
        - ALL
    readOnlyRootFilesystem: true
    allowPrivilegeEscalation: false

frontend:
  securityContext:
    runAsNonRoot: true
    runAsUser: 101
    runAsGroup: 101
    capabilities:
      drop:
        - ALL
    readOnlyRootFilesystem: true
    allowPrivilegeEscalation: false
```

## テスト

### バックエンドテスト
```bash
cd backend
go test ./internal/handlers/...
# モックPromtheusクライアントを使用した単体テスト
```

### フロントエンドテスト
```bash
cd frontend
npm run lint          # ESLint
npm run type-check     # TypeScriptチェック
```

### Helmチャートテスト
```bash
helm test gpu-monitoring --namespace gpu-monitoring
```

## API仕様

### エンドポイント

| Method | Path | 説明 | レスポンス |
|--------|------|------|-----------|
| GET | `/api/health` | ヘルスチェック・Prometheus接続確認 | `APIResponse` |
| GET | `/api/v1/gpu/metrics` | 全GPUの詳細メトリクス | `APIResponse<GPUMetrics[]>` |
| GET | `/api/v1/gpu/nodes` | GPU搭載ノード一覧 | `APIResponse<GPUNode[]>` |
| GET | `/api/v1/gpu/utilization` | GPU利用率のみ（軽量） | `APIResponse<GPUUtilization[]>` |

### データ構造

```typescript
interface GPUMetrics {
  node_name: string;
  gpu_index: number;
  gpu_name: string;
  utilization: number;          // GPU利用率 (%)
  memory_used: number;          // 使用メモリ (GB)
  memory_total: number;         // 総メモリ (GB)
  memory_free: number;          // 空きメモリ (GB)
  memory_utilization: number;   // メモリ利用率 (%)
  temperature: number;          // 温度 (℃)
  power_draw: number;          // 電力消費 (W)
  power_limit: number;         // 電力制限 (W)
  timestamp: string;           // タイムスタンプ
}

interface APIResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}
```

## 監視・運用

### ヘルスチェック
```bash
# Backend API
kubectl exec -n gpu-monitoring deployment/gpu-monitoring-backend -- \
  wget -qO- http://localhost:8080/api/health

# Frontend
kubectl exec -n gpu-monitoring deployment/gpu-monitoring-frontend -- \
  wget -qO- http://localhost:80/health
```

### ログ確認
```bash
# Backend ログ
kubectl logs -n gpu-monitoring deployment/gpu-monitoring-backend -f

# Frontend ログ
kubectl logs -n gpu-monitoring deployment/gpu-monitoring-frontend -f

# 全体ログ
kubectl logs -n gpu-monitoring -l app.kubernetes.io/name=k8s-gpu-monitoring-dev -f
```

### リソース確認
```bash
# Pod状態確認
kubectl get pods -n gpu-monitoring

# リソース使用量確認
kubectl top pods -n gpu-monitoring

# 詳細情報
kubectl describe pods -n gpu-monitoring
```

## アップグレード

### Helmでのアップグレード
```bash
# リポジトリ更新
helm repo update

# アップグレード
helm upgrade gpu-monitoring gpu-monitoring/k8s-gpu-monitoring-dev \
  --namespace gpu-monitoring

# 特定バージョンにアップグレード
helm upgrade gpu-monitoring gpu-monitoring/k8s-gpu-monitoring-dev \
  --namespace gpu-monitoring \
  --version 1.0.1
```

### 設定変更アップグレード
```bash
# リソース設定変更
helm upgrade gpu-monitoring gpu-monitoring/k8s-gpu-monitoring-dev \
  --namespace gpu-monitoring \
  --set backend.resources.requests.cpu=500m \
  --set backend.resources.requests.memory=512Mi
```

## セキュリティ

- **CORS設定**: 適切なクロスオリジン設定とプリフライトリクエスト処理
- **入力検証**: バックエンドでの厳密な入力バリデーション
- **セキュリティヘッダー**: Nginxでの包括的なセキュリティヘッダー設定
- **非rootユーザー**: 全Dockerコンテナでの非特権実行
- **読み取り専用ファイルシステム**: セキュリティ強化のためのイミュータブル実行環境
- **シンプル設計**: Prometheus HTTPクライアントのみ使用（Kubernetes API不要）

## 開発・コントリビューション

### プロジェクト構造の拡張

1. **バックエンド新機能**: `internal/handlers/` にハンドラー追加
2. **フロントエンド新機能**: `src/components/` にコンポーネント追加
3. **API拡張**: `src/api/client.ts` にAPI関数追加
4. **新しい型定義**: `src/types/` にTypeScript型追加

### CI/CDパイプライン

- **自動ビルド**: タグプッシュ時のDockerイメージ自動ビルド
- **マルチアーキテクチャ対応**: linux/amd64のみ（GPUノード対応）
- **Helmチャート公開**: GitHub Pagesでの自動公開
- **リリース自動化**: scripts/release.sh による完全自動化

### リリースプロセス

```bash
# 新バージョンリリース
./scripts/release.sh 1.0.1
# 1. Chart.yamlとvalues.yamlのバージョン更新
# 2. Git コミット・タグ作成
# 3. GitHub Actions による自動ビルド・デプロイ
```

## トラブルシューティング

### よくある問題

1. **API接続エラー**: 
   ```bash
   kubectl exec -n gpu-monitoring deployment/gpu-monitoring-backend -- \
     wget -qO- http://prometheus-server:9090/api/v1/query?query=up
   ```

2. **メトリクス取得エラー**: NVIDIA GPU Exporterの設定確認
   ```bash
   kubectl get pods -A | grep nvidia
   kubectl logs -n monitoring nvidia-gpu-exporter-xxx
   ```

3. **フロントエンド表示エラー**: ブラウザ開発者ツールでAPI通信確認

4. **Ingressアクセス不可**: 
   ```bash
   kubectl get ingress -n gpu-monitoring
   # /etc/hostsにドメイン追加確認
   ```

### .localドメインの設定
```bash
# Ingress Controller のIPを確認
kubectl get svc -n ingress-nginx ingress-nginx-controller

# /etc/hostsに追加 (Linux/Mac)
echo "192.168.1.100 gpu-monitoring.local" | sudo tee -a /etc/hosts

# Windows
# C:\Windows\System32\drivers\etc\hosts に追加
# 192.168.1.100 gpu-monitoring.local
```

## アンインストール

```bash
# アプリケーション削除
helm uninstall gpu-monitoring --namespace gpu-monitoring

# Namespace削除
kubectl delete namespace gpu-monitoring

# /etc/hostsエントリ削除
sudo sed -i '/gpu-monitoring.local/d' /etc/hosts
```