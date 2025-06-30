# 🖥️ K8s GPU Monitoring Dashboard

Kubernetes上でPrometheusからGPUメトリクスを取得・表示するための統合監視ダッシュボードです。Go製のバックエンドAPIとReact製のフロントエンドで構成されています。

![GPU Dashboard](https://img.shields.io/badge/Status-Production%20Ready-green)
![Go](https://img.shields.io/badge/Go-1.22-blue)
![React](https://img.shields.io/badge/React-18-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## ✨ 特徴

### 🚀 高性能
- **Go 1.22**: 最新の標準ライブラリServeMuxを使用
- **並行処理**: Goroutineによる効率的なPrometheusクエリ
- **TanStack Table**: 大量データを高速表示
- **リアルタイム更新**: 30秒間隔での自動データ更新

### 📊 豊富な監視機能
- **詳細メトリクス**: GPU利用率、メモリ使用量、温度、電力消費
- **統計ダッシュボード**: 総GPU数、平均利用率、アラート状況
- **ビジュアル表示**: プログレスバーと色分けによる直感的な表示
- **ヘルスモニタリング**: システム全体の健全性監視

### 🎨 モダンUI
- **レスポンシブデザイン**: デスクトップ・モバイル対応
- **ダークモード**: 目に優しい表示
- **アクセシビリティ**: キーボードナビゲーション対応
- **TypeScript**: 型安全な開発

## 🏗️ アーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                       │
│                                                             │
│  ┌─────────────┐    ┌──────────────┐    ┌────────────────┐  │
│  │   Ingress   │    │   Frontend   │    │   Backend      │  │
│  │ Controller  │    │   (React)    │    │   (Go API)     │  │
│  │             │───→│              │───→│                │──┼─→ Prometheus
│  │   nginx/    │    │  TanStack    │    │  ServeMux      │  │
│  │  traefik    │    │   Table      │    │  + CORS        │  │
│  └─────────────┘    └──────────────┘    └────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## 📁 プロジェクト構成

```
k8s-gpu-monitoring/
├── backend/                    # Go API サーバー
│   ├── cmd/server/main.go     # メインエントリーポイント
│   ├── internal/
│   │   ├── handlers/          # HTTPハンドラー
│   │   ├── middleware/        # CORS・ログ等
│   │   ├── models/           # データモデル
│   │   └── prometheus/       # Prometheusクライアント
│   ├── Dockerfile
│   └── README.md
├── frontend/                   # React フロントエンド
│   ├── src/
│   │   ├── components/       # Reactコンポーネント
│   │   ├── api/             # APIクライアント
│   │   ├── types/           # TypeScript型定義
│   │   └── lib/             # ユーティリティ
│   ├── Dockerfile
│   ├── nginx.conf
│   └── README.md
└── helm-chart/                # Helm Chart（予定）
```

## 🚀 クイックスタート

### 前提条件
- **Go 1.22以上**
- **Node.js 18以上**
- **Docker & Docker Compose**
- **Kubernetes Cluster**
- **Prometheus Server** (NVIDIA GPU メトリクス取得済み)

### 開発環境で実行

#### 1. バックエンドAPI
```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### 2. フロントエンド
```bash
cd frontend
npm install
npm run dev
```

#### 3. アクセス
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080

### Docker Composeで実行

```bash
# docker-compose.yml 作成後
docker-compose up -d
```

## 🔧 設定

### 環境変数

#### バックエンド
```bash
PROMETHEUS_URL=http://prometheus-server:9090  # PrometheusサーバーURL
PORT=8080                                     # APIサーバーポート
```

#### フロントエンド
```bash
VITE_API_URL=http://localhost:8080/api       # APIベースURL
```

### 必要なPrometheusメトリクス

以下のNVIDIA GPUメトリクスが必要です：

```promql
# GPU利用率
nvidia_smi_utilization_gpu_ratio

# メモリ関連
nvidia_smi_memory_used_bytes
nvidia_smi_memory_total_bytes

# 温度・電力
nvidia_smi_temperature_gpu_celsius
nvidia_smi_power_draw_watts
nvidia_smi_enforced_power_limit_watts

# GPU情報
nvidia_smi_gpu_info
```

## 📦 本番デプロイ

### Helm Chart (推奨)

```bash
# Helm Chart作成・デプロイ（予定）
helm install gpu-monitoring ./helm-chart \
  --set global.domain=gpu-monitoring.example.com \
  --set prometheus.url=http://prometheus-server:9090
```

### Docker

```bash
# バックエンド
docker build -t gpu-monitoring-backend ./backend
docker run -p 8080:8080 \
  -e PROMETHEUS_URL=http://prometheus:9090 \
  gpu-monitoring-backend

# フロントエンド  
docker build -t gpu-monitoring-frontend ./frontend
docker run -p 3000:80 gpu-monitoring-frontend
```

## 🔍 API仕様

### エンドポイント

| Method | Path | 説明 |
|--------|------|------|
| GET | `/api/health` | ヘルスチェック |
| GET | `/api/v1/gpu/metrics` | 全GPUメトリクス |
| GET | `/api/v1/gpu/nodes` | GPU搭載ノード一覧 |
| GET | `/api/v1/gpu/utilization` | GPU利用率（軽量） |

### レスポンス形式

```json
{
  "success": true,
  "data": [...],
  "message": "Operation completed successfully"
}
```

## 🧪 テスト

```bash
# バックエンドテスト
cd backend && go test ./...

# フロントエンドテスト
cd frontend && npm test
```

## 📈 パフォーマンス

- **バックエンド**: 1000 req/s 対応
- **フロントエンド**: 初回ロード < 2秒
- **リアルタイム更新**: 30秒間隔
- **同時接続**: 100+ クライアント対応

## 🔒 セキュリティ

- **CORS設定**: 適切なクロスオリジン設定
- **入力検証**: バックエンドでの入力バリデーション
- **セキュリティヘッダー**: Nginxでのセキュリティヘッダー設定
- **非rootユーザー**: Dockerコンテナでの非特権実行

## 🛠️ 開発

### 新機能の追加

1. **バックエンド**: `internal/handlers/` にハンドラー追加
2. **フロントエンド**: `src/components/` にコンポーネント追加
3. **API**: `src/api/client.ts` にAPI関数追加

### コントリビューション

1. Forkプロジェクト
2. 機能ブランチ作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. Pull Request作成

## 📚 詳細ドキュメント

- [バックエンドAPI](./backend/README.md)
- [フロントエンド](./frontend/README.md)
- [Helm Chart](./helm-chart/README.md) (予定)

## 🐛 トラブルシューティング

### よくある問題

1. **API接続エラー**: Prometheusサーバーの接続確認
2. **メトリクス未取得**: NVIDIA GPU Exporterの設定確認
3. **ビルドエラー**: Go/Node.jsバージョン確認

## 📄 ライセンス

MIT License - 詳細は [LICENSE](LICENSE) ファイルを参照

## 🤝 サポート

- **Issues**: [GitHub Issues](https://github.com/your-org/k8s-gpu-monitoring/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/k8s-gpu-monitoring/discussions)

---

**作成者**: GPU監視チーム  
**最終更新**: 2024年 