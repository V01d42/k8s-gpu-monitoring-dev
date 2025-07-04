# GPU Monitoring Frontend

K8s上でPrometheusからGPUメトリクスを取得・表示するためのReact 19 + TypeScript 5.7 + Vite 7.0フロントエンドアプリケーションです。

## 特徴

- **最新技術スタック**: React 19.1.0 + TypeScript 5.7 + Vite 7.0で超高速開発
- **TanStack Table v8**: 高性能なテーブルコンポーネントで大量データを効率表示
- **TanStack Query v5**: 効率的なデータフェッチ・キャッシュ・状態管理
- **TailwindCSS 4.1**: 最新のユーティリティファーストCSS
- **リアルタイム更新**: 30秒間隔での自動データ更新と手動更新機能
- **ダークモードデフォルト**: 目に優しいダークテーマをデフォルト採用
- **レスポンシブデザイン**: デスクトップ・モバイル完全対応

## 主要機能

### 統合ダッシュボード
- **統計カード**: 総GPU数、平均利用率、高温アラート数、アクティブ率
- **ヘルスステータス**: バックエンドAPI接続状態のリアルタイム監視
- **自動更新**: 30秒間隔でのメトリクス自動更新
- **手動更新**: ユーザー操作による即座のデータ更新

### 高性能GPUテーブル
- **詳細メトリクス**: 利用率、メモリ使用量、温度、電力消費等の包括的表示
- **ビジュアル表示**: プログレスバーによる直感的な利用率・メモリ使用率表示
- **インテリジェント色分け**: 利用率と温度に応じた色分け（緑/黄/赤）
- **全列ソート機能**: 各列での昇順・降順ソート対応
- **レスポンシブテーブル**: モバイルデバイスでの最適化表示

### モダンUI/UX
- **WCAG準拠デザイン**: アクセシビリティを重視したインターフェース
- **ローディング状態**: skeleton UIとスピナーによる適切なローディング表示
- **エラーハンドリング**: 分かりやすいエラーメッセージと回復操作
- **キーボードナビゲーション**: 完全なキーボード操作対応

## 技術スタック

### コアフレームワーク
- **React 19.1.0**: 最新のReact機能（Concurrent Features、Suspense等）
- **TypeScript 5.7**: 最新の型安全機能とパフォーマンス向上
- **Vite 7.0**: 超高速ビルドツールとHMR（Hot Module Replacement）

### 状態管理・データフェッチ
- **TanStack Query v5**: サーバー状態管理・キャッシュ・背景更新
- **TanStack Table v8**: 高性能テーブル仮想化・ソート・フィルタリング

### スタイリング・UI
- **TailwindCSS 4.1**: ユーティリティファーストCSS、最新の機能
- **Lucide React**: モダンなアイコンライブラリ
- **shadcn/ui**: 再利用可能なコンポーネントライブラリ

### HTTP・通信
- **Axios**: 型安全なHTTPクライアント
- **自動プロキシ**: 開発時のvite.config.tsプロキシ設定

## 開発環境セットアップ

### 前提条件
- **Node.js 20以上** (mise.tomlでlatestを指定)
- **npm または yarn または pnpm**

### インストール
```bash
# 依存関係をインストール（npm ciで確定的インストール）
npm ci

# 開発サーバーを起動
npm run dev
# サーバーが http://localhost:3000 で起動
```

### 利用可能なスクリプト
```bash
# 開発サーバー起動（ポート3000、HMR有効）
npm run dev

# プロダクションビルド（最適化済み）
npm run build

# 静的解析（ESLint）
npm run lint

# TypeScript型チェック
npm run type-check

# プレビューサーバー起動（ビルド後の確認用）
npm run preview
```

## 環境変数

本番環境では Nginx プロキシを使用するため、開発時のみ環境変数を使用：

```env
# 開発時のAPIベースURL（オプション）
VITE_API_URL=http://localhost:8080/api
```

**注意**: 本番環境では `nginx.conf` の設定によりAPIプロキシが動作します。

## API連携

フロントエンドは以下のAPIエンドポイントを使用：

### REST API
- `GET /api/health` - ヘルスチェック・Prometheus接続状態
- `GET /api/v1/gpu/metrics` - 全GPUの詳細メトリクス
- `GET /api/v1/gpu/nodes` - GPU搭載ノード一覧
- `GET /api/v1/gpu/utilization` - GPU利用率のみ（軽量）

### データフロー
```
TanStack Query → Axios → Backend API → Prometheus
     ↓
 React State → TanStack Table → UI Components
```

## ディレクトリ構成

```
src/
├── api/                    # API クライアント
│   └── client.ts          # Axiosベースの型安全APIクライアント
├── components/             # React コンポーネント
│   ├── ui/                # 再利用可能UIコンポーネント
│   │   ├── button.tsx     # ボタンコンポーネント
│   │   └── card.tsx       # カードコンポーネント
│   └── GPUTable.tsx       # メインテーブルコンポーネント
├── lib/                   # ユーティリティ関数
│   └── utils.ts           # clsx・className結合等
├── types/                 # TypeScript型定義
│   └── gpu.ts             # GPU関連の型定義
├── App.tsx                # メインアプリケーション
├── main.tsx               # React エントリーポイント
├── index.css              # TailwindCSS 設定・カスタムスタイル
└── vite-env.d.ts          # Vite 型定義
```

## Docker

### ビルド
```bash
# マルチステージビルドでイメージを構築
docker build -t gpu-monitoring-frontend .

# 最適化されたプロダクションビルド
docker build --target production -t gpu-monitoring-frontend:prod .
```

### 実行
```bash
# コンテナを実行
docker run -p 3000:80 gpu-monitoring-frontend

# ヘルスチェック付きで実行
docker run -p 3000:80 \
  --health-cmd="wget --no-verbose --tries=1 --spider http://localhost:80/health || exit 1" \
  --health-interval=30s \
  --health-timeout=10s \
  --health-retries=3 \
  gpu-monitoring-frontend
```

### docker-compose
```yaml
version: '3.8'
services:
  frontend:
    build: 
      context: .
      target: production
    ports:
      - "3000:80"
    depends_on:
      - backend
    environment:
      - NODE_ENV=production
```

## 本番デプロイ

### Nginx設定
本番環境では、最適化された `nginx.conf` を使用：

```nginx
# 主要設定内容
- Gzip圧縮有効
- 静的アセットの長期キャッシュ（1年）
- セキュリティヘッダー設定
- APIプロキシ設定 (/api/* → backend:8080)
- SPA用のfallback設定
```

### ビルド最適化
- **Tree Shaking**: 使用されないコードの自動削除
- **Code Splitting**: 動的インポートによる最適化
- **Asset Optimization**: 画像・フォント等の最適化
- **Bundle Analysis**: webpack-bundle-analyzerによる分析

### セキュリティ
```nginx
# セキュリティヘッダー（nginx.conf内）
add_header X-Frame-Options "SAMEORIGIN";
add_header X-Content-Type-Options "nosniff";
add_header X-XSS-Protection "1; mode=block";
add_header Referrer-Policy "strict-origin-when-cross-origin";
```

## パフォーマンス

### ベンチマーク指標
- **初回ロード**: < 1.5秒（Fast 3G環境）
- **データ更新**: 30秒間隔の自動更新（バックグラウンド実行）
- **キャッシュ**: TanStack Query による効率的なデータキャッシュ
- **バンドルサイズ**: gzip圧縮後 < 400KB（最適化済み）

### 最適化技術
- **React 19 Concurrent Features**: 非ブロッキングレンダリング
- **TanStack Table 仮想化**: 大量データの効率的表示
- **メモ化**: React.memo、useMemo、useCallbackの適切な使用
- **Suspense**: 非同期データローディングの最適化

## ブラウザ対応

### サポート対象
- **Chrome/Edge**: 90+
- **Firefox**: 88+
- **Safari**: 14+
- **モバイル**: iOS Safari 14+, Android Chrome 90+

### 機能サポート
- **ES2022**: 最新のJavaScript機能
- **CSS Grid & Flexbox**: モダンレイアウト
- **CSS Custom Properties**: テーマ対応
- **Web APIs**: Fetch API、ResizeObserver等

## トラブルシューティング

### よくある問題

1. **API接続エラー**
   ```bash
   # バックエンドサーバーの状態確認
   curl http://localhost:8080/api/health
   
   # プロキシ設定確認（vite.config.ts）
   npm run dev -- --debug
   ```

2. **ビルドエラー**
   ```bash
   # Node.js バージョン確認（20以上必要）
   node --version
   
   # 依存関係を再インストール
   rm -rf node_modules package-lock.json
   npm ci
   
   # TypeScript型エラー確認
   npm run type-check
   ```

3. **表示が崩れる・スタイルエラー**
   ```bash
   # Tailwind CSS ビルド確認
   npm run build
   
   # ブラウザキャッシュクリア
   # Ctrl+Shift+R (Windows/Linux) または Cmd+Shift+R (Mac)
   
   # ダークモードの設定確認
   # ブラウザ開発者ツール → Elements → <html class="dark">
   ```

4. **メモリ不足・パフォーマンス問題**
   ```bash
   # Node.js メモリ制限を増加
   export NODE_OPTIONS="--max_old_space_size=4096"
   npm run build
   
   # バンドルサイズ分析
   npm run build -- --analyze
   ```

### デバッグ

#### 開発時のデバッグ
```bash
# 詳細ログでVite開発サーバーを起動
npm run dev -- --debug

# React DevTools でコンポーネント状態確認
# TanStack Query DevTools でクエリ状態確認
```

#### 本番環境のデバッグ
```bash
# Nginx ログ確認
docker logs <container_id>

# Static ファイル確認
curl -I http://localhost:3000/assets/index.js

# API プロキシ確認
curl http://localhost:3000/api/health
```

## 開発ガイド

### 新機能の追加

1. **新しいコンポーネント**
   ```bash
   # src/components/ に追加
   # TypeScript + Props インターフェース定義
   # shadcn/ui パターンに従った実装
   ```

2. **新しいAPI統合**
   ```bash
   # src/api/client.ts にAPI関数追加
   # src/types/ に型定義追加
   # TanStack Query フックの作成
   ```

3. **スタイリング**
   ```bash
   # TailwindCSS ユーティリティクラス使用
   # カスタムCSSは最小限に抑制
   # ダークモード対応必須
   ```
