# openLive

自由にライブ配信ができるプラットフォーム

基盤プラットフォーム内にユーザーが公開のライブ配信エリアを作成し、会員またはフリーでライブ配信ができるプラットフォーム環境。
ユーザーは独自の広場（配信スペース）を構築し、ユーザーを招待し、広場内でのみ交流可能となるマルチプラットフォーム型のライブ配信アプリ。

## プロジェクト構成

このプロジェクトはモノレポ構成で、以下のディレクトリで構成されています：

```
.
├── server/              # Goバックエンド
│   ├── cmd/            # アプリケーションエントリーポイント
│   ├── internal/       # 内部パッケージ
│   │   ├── domain/     # ドメイン層（DDD）
│   │   ├── usecase/    # ユースケース層
│   │   ├── interface/  # インターフェース層
│   │   ├── infrastructure/ # インフラ層
│   │   └── config/     # 設定管理
│   ├── migrations/     # DBマイグレーション
│   └── go.mod
├── frontend/           # Next.jsフロントエンド
│   ├── src/
│   │   ├── app/       # Next.js App Router
│   │   ├── components/ # UIコンポーネント（SOLID原則）
│   │   ├── hooks/     # カスタムフック
│   │   ├── services/  # APIサービス層
│   │   ├── store/     # 状態管理（Zustand）
│   │   └── types/     # TypeScript型定義
│   └── package.json
├── docker/             # Docker関連ファイル
│   ├── server/
│   └── firebase/
└── docker-compose.yml  # 統合Docker Compose設定
```

## 技術スタック

### Backend (Server)

- **言語**: Go 1.21
- **フレームワーク**: Echo v4.12.0
- **データベース**: PostgreSQL 16
- **認証**: Firebase Auth (Emulator対応)
- **リアルタイム通信**: WebSocket + WebRTC
- **アーキテクチャ**: DDD (Domain-Driven Design) + Clean Architecture

### Frontend

- **フレームワーク**: Next.js 14 (App Router)
- **言語**: TypeScript
- **UIライブラリ**: Tailwind CSS
- **状態管理**: Zustand
- **APIクライアント**: Axios
- **認証**: Firebase Authentication
- **リアルタイム通信**: WebSocket, WebRTC (simple-peer)
- **設計原則**: SOLID原則に準拠

## 主な機能

### 1. 認証機能
- Firebase AuthによるGoogle OAuth2認証
- トークンベースの認証
- エミュレーター環境での開発対応

### 2. ライブエリア管理
- 配信広場の作成・管理
- 公開/非公開の設定
- メンバー招待・権限管理

### 3. ライブ配信機能
- WebRTCによるリアルタイム配信
- 配信開始・終了管理
- 視聴者管理

### 4. リアルタイム通信
- WebSocketによるチャット機能
- WebRTCシグナリング（Offer/Answer/ICE Candidate交換）
- 視聴者数のリアルタイム更新

## セットアップ

### 前提条件

- Docker & Docker Compose
- Go 1.21+ (ローカル実行の場合)
- Node.js 18+ (ローカル実行の場合)

### Docker環境での起動（推奨）

**分離版構成**: サービスごとに独立して起動・停止できます

```bash
# 方法1: Makefileを使用（推奨）
make all          # 全サービスを一括起動
# または
make infra        # 1. 共通インフラを起動
make server       # 2. Serverを起動
make frontend     # 3. Frontendを起動

# ログを確認
make logs         # 全サービスのログ
make logs-server  # Serverのみ
make logs-frontend # Frontendのみ

# サービスの停止
make down         # 全サービス停止
make server-down  # Serverのみ停止
make frontend-down # Frontendのみ停止

# 方法2: docker-composeコマンドを直接使用
docker-compose up -d                                # 共通インフラのみ
docker-compose -f docker-compose.server.yml up -d   # Server追加
docker-compose -f docker-compose.frontend.yml up -d # Frontend追加
```

起動後、以下のURLでアクセスできます：
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080
- Firebase Emulator UI: http://localhost:4000

**Makefileコマンド一覧**:
```bash
make help          # 全コマンドを表示
make status        # サービスのステータス確認
make restart-server # Serverのみ再起動
make rebuild       # 全サービスをリビルド
make clean         # 全コンテナとボリュームを削除
make db            # PostgreSQLに接続
```

### ローカル環境での起動

#### Backend

```bash
cd server

# 依存関係のインストール
go mod download

# 環境変数の設定
cp env.example .env
# .envファイルを編集

# PostgreSQLとFirebase Emulatorを起動
cd ..
docker-compose up -d postgres firebase

# サーバーを起動
cd server
go run cmd/server/main.go
```

#### Frontend

```bash
cd frontend

# 依存関係のインストール
npm install

# 開発サーバーを起動
npm run dev
```

## Docker Compose構成（分離版）

### ファイル構成

このプロジェクトは**分離版**のDocker Compose構成を採用しています：

1. **`docker-compose.yml`** - 共通インフラ（PostgreSQL + Firebase Emulator）
2. **`docker-compose.server.yml`** - Goバックエンドサーバー
3. **`docker-compose.frontend.yml`** - Next.jsフロントエンド

### 分離版のメリット

- ✅ **独立した起動・停止**: フロントエンド/バックエンドを個別に管理
- ✅ **リソース効率**: Frontend開発時にServerを起動しなくてもOK
- ✅ **並行開発**: 異なるチームが別々に開発可能
- ✅ **柔軟な再起動**: Serverだけ、Frontendだけを再起動可能
- ✅ **個別デプロイ**: CI/CDパイプラインで個別にデプロイ可能
- ✅ **スケーリング**: サービスごとに独立してスケール可能

### 使用例

```bash
# Frontend開発の場合
make infra          # インフラだけ起動
make server         # Serverを起動
make frontend       # Frontendを起動・開発
make restart-frontend # Frontendだけ再起動（Serverは影響なし）

# Backend開発の場合
make infra          # インフラだけ起動
make server         # Serverを起動・開発
make restart-server # Serverだけ再起動（Frontendは影響なし）

# 全サービス起動
make all            # 全部まとめて起動
```

## API エンドポイント

### 認証

すべての認証が必要なエンドポイントには、Authorizationヘッダーに`Bearer <Firebase ID Token>`を含める必要があります。

### ユーザー

- `GET /api/v1/users/me` - 自分の情報取得
- `PUT /api/v1/users/me` - プロフィール更新

### ライブエリア

- `POST /api/v1/liveareas` - エリア作成 (認証必須)
- `GET /api/v1/liveareas` - エリア一覧取得
- `GET /api/v1/liveareas/:id` - エリア詳細取得
- `PUT /api/v1/liveareas/:id` - エリア更新 (認証必須)
- `DELETE /api/v1/liveareas/:id` - エリア削除 (認証必須)
- `POST /api/v1/liveareas/:id/invite` - メンバー招待 (認証必須)
- `DELETE /api/v1/liveareas/:id/members/:userId` - メンバー削除 (認証必須)
- `GET /api/v1/liveareas/:id/members` - メンバー一覧取得

### ライブ配信

- `POST /api/v1/liveareas/:id/streams` - 配信作成 (認証必須)
- `GET /api/v1/liveareas/:id/streams` - エリアの配信一覧取得
- `GET /api/v1/streams/live` - 配信中の配信一覧取得
- `GET /api/v1/streams/:id` - 配信詳細取得
- `POST /api/v1/streams/:id/start` - 配信開始 (認証必須)
- `DELETE /api/v1/streams/:id` - 配信終了 (認証必須)
- `GET /api/v1/streams/:id/viewers` - 視聴者一覧取得
- `GET /api/v1/streams/:id/viewers/count` - 視聴者数取得

### WebSocket

- `WS /api/v1/ws/streams/:id` - WebSocket接続 (認証必須)

## アーキテクチャ

### Backend（DDD + Clean Architecture）

1. **ドメイン層（Domain）**: ビジネスロジックとエンティティ、リポジトリインターフェース
2. **ユースケース層（UseCase）**: アプリケーションのビジネスルール
3. **インターフェース層（Interface）**: HTTP/WebSocketハンドラー、DTO、ミドルウェア
4. **インフラストラクチャ層（Infrastructure）**: データベース、外部サービス、リポジトリ実装

### Frontend（SOLID原則）

1. **Single Responsibility Principle（単一責任の原則）**
   - 各コンポーネント、サービスは単一の責任を持つ
   - 例：`Button`コンポーネントはボタンの表示のみ、`UserService`はユーザー管理のみ

2. **Open/Closed Principle（開放閉鎖の原則）**
   - 拡張に対して開かれており、変更に対して閉じている
   - 例：共通コンポーネントは拡張可能だが、既存の実装を変更する必要がない

3. **Liskov Substitution Principle（リスコフの置換原則）**
   - インターフェースの実装は置き換え可能
   - 例：`IHttpClient`の実装を別の実装に置き換えても動作する

4. **Interface Segregation Principle（インターフェース分離の原則）**
   - インターフェースは細かく分割されている
   - 例：`IUserService`, `ILiveAreaService`, `ILiveStreamService`など

5. **Dependency Inversion Principle（依存性逆転の原則）**
   - 上位モジュールは下位モジュールに依存しない、抽象に依存する
   - 例：サービスは`IHttpClient`インターフェースに依存し、具体的な実装には依存しない

## 開発

### ホットリロード

開発時はホットリロードが有効になっています：

**Backend**: Airを使用（自動リロード）
**Frontend**: Next.js Dev Server（自動リロード）

### データベースマイグレーション

```bash
# Docker環境の場合、初回起動時に自動実行されます

# 手動で実行する場合
docker-compose exec postgres psql -U openlive -d openlive_db -f /docker-entrypoint-initdb.d/001_init.sql
```

### コードフォーマット

```bash
# Backend
cd server
go fmt ./...

# Frontend
cd frontend
npm run lint
```

## トラブルシューティング

### ポートが使用されている場合

```bash
# 使用中のポートを確認
lsof -i :3000  # Frontend
lsof -i :8080  # Backend
lsof -i :5432  # PostgreSQL

# プロセスを停止
kill -9 <PID>
```

### Dockerコンテナのリセット

```bash
# 全コンテナを停止して削除
docker-compose down -v

# イメージも削除
docker-compose down -v --rmi all

# 再ビルドして起動
docker-compose up --build
```

## ライセンス

MIT
