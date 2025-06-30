# GPU Monitoring Frontend

K8s上でPrometheusからGPUメトリクスを取得・表示するためのReact + TypeScriptフロントエンドアプリケーションです。

## 特徴

- **Modern React**: React 18 + TypeScript + Vite で高速開発
- **TanStack Table**: 高性能なテーブルコンポーネントで大量データを効率表示
- **TanStack Query**: 効率的なデータフェッチとキャッシュ機能
- **Tailwind CSS**: モダンでレスポンシブなUI
- **リアルタイム更新**: 30秒間隔での自動データ更新
- **ダークモード**: 目に優しいダークテーマ対応

## 主要機能

### 📊 ダッシュボード
- **統計カード**: 総GPU数、平均利用率、高温アラート、アクティブ率
- **ヘルスステータス**: バックエンドAPI接続状態の監視
- **自動更新**: リアルタイムでのメトリクス更新

### 📋 GPUテーブル
- **詳細メトリクス**: 利用率、メモリ使用量、温度、電力消費など
- **ビジュアル表示**: プログレスバーによる直感的な利用率表示
- **色分け表示**: 利用率と温度に応じた色分け（緑/黄/赤）
- **ソート機能**: 各列でのソート対応
- **レスポンシブ**: モバイルデバイス対応

### 🎨 UI/UX
- **モダンデザイン**: Clean で見やすいインターフェース
- **ローディング状態**: 適切なローディングインジケーター
- **エラーハンドリング**: 分かりやすいエラーメッセージ
- **アクセシビリティ**: キーボードナビゲーション対応

## 技術スタック

- **React 18**: 最新のReact機能
- **TypeScript**: 型安全な開発
- **Vite**: 高速ビルドツール
- **TanStack Table v8**: 高性能テーブル
- **TanStack Query v5**: データフェッチ & キャッシュ
- **Tailwind CSS**: ユーティリティファーストCSS
- **Lucide React**: アイコンライブラリ
- **Axios**: HTTP クライアント

## 開発環境セットアップ

### 前提条件
- Node.js 18以上
- npm または yarn

### インストール
```bash
# 依存関係をインストール
npm install

# 開発サーバーを起動
npm run dev
```

### 利用可能なスクリプト
```bash
# 開発サーバー起動（ポート3000）
npm run dev

# プロダクションビルド
npm run build

# 静的解析（ESLint）
npm run lint

# プレビューサーバー起動
npm run preview
```

## 環境変数

`.env`ファイルで設定可能：

```env
# APIベースURL（オプション）
VITE_API_URL=http://localhost:8080/api
```

## API連携

フロントエンドは以下のAPIエンドポイントを使用：

- `GET /api/health` - ヘルスチェック
- `GET /api/v1/gpu/metrics` - 全GPUメトリクス
- `GET /api/v1/gpu/nodes` - GPU搭載ノード一覧
- `GET /api/v1/gpu/utilization` - GPU利用率（軽量）

## ディレクトリ構成

```
src/
├── api/            # API クライアント
├── components/     # React コンポーネント
│   ├── ui/        # 再利用可能UIコンポーネント
│   └── GPUTable.tsx   # メインテーブルコンポーネント
├── lib/           # ユーティリティ関数
├── types/         # TypeScript型定義
├── App.tsx        # メインアプリケーション
├── main.tsx       # エントリーポイント
└── index.css      # グローバルスタイル
```

## Docker

### ビルド
```bash
# イメージをビルド
docker build -t gpu-monitoring-frontend .
```

### 実行
```bash
# コンテナを実行
docker run -p 3000:80 gpu-monitoring-frontend
```

### docker-compose
```yaml
version: '3.8'
services:
  frontend:
    build: .
    ports:
      - "3000:80"
    depends_on:
      - backend
```

## 本番デプロイ

### Nginx設定
本番環境では、含まれている`nginx.conf`を使用してNginxでホスティング。

### 最適化
- **Gzip圧縮**: 有効化済み
- **キャッシュ設定**: 静的アセットの長期キャッシュ
- **セキュリティヘッダー**: 適切なHTTPセキュリティヘッダー
- **バンドル最適化**: Viteによる最適化されたビルド

## パフォーマンス

- **初回ロード**: < 2秒
- **データ更新**: 30秒間隔の自動更新
- **キャッシュ**: TanStack Query による効率的なキャッシュ
- **バンドルサイズ**: gzip圧縮後 < 500KB

## ブラウザ対応

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- モバイルブラウザ対応

## トラブルシューティング

### よくある問題

1. **API接続エラー**
   - バックエンドサーバーが起動しているか確認
   - CORSの設定を確認

2. **ビルドエラー**
   - Node.js バージョンを確認（18以上）
   - `npm ci` で依存関係を再インストール

3. **表示が崩れる**
   - ブラウザのキャッシュをクリア
   - ダークモードの設定を確認

### デバッグ

開発者ツールで以下を確認：
- Network タブでAPI通信状況
- Console タブでエラーメッセージ
- React Query Devtools でデータ状態 