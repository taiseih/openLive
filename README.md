# openLive

自由にライブ配信ができるプラットフォーム

基盤プラットフォーム内にユーザーが公開のライブ配信エリアを作成し、会員またはフリーでライブ配信ができるプラットフォーム環境。
ユーザーは独自の広場（配信スペース）を構築し、ユーザーを招待し、広場内でのみ交流可能となるマルチプラットフォーム型のライブ配信アプリ。

## 技術スタック

- **言語**: Go 1.21
- **フレームワーク**: Echo v4.12.0
- **データベース**: PostgreSQL 16
- **認証**: Firebase Auth (Emulator対応)
- **リアルタイム通信**: WebSocket + WebRTC
- **アーキテクチャ**: DDD (Domain-Driven Design) + Clean Architecture

## プロジェクト構造

```
.
├── cmd/
│   └── server/           # アプリケーションエントリーポイント
├── internal/
│   ├── domain/           # ドメイン層（エンティティ、リポジトリインターフェース）
│   ├── usecase/          # ユースケース層（ビジネスロジック）
│   ├── interface/        # インターフェース層（ハンドラー、DTO、ミドルウェア）
│   ├── infrastructure/   # インフラ層（DB、Firebase、WebSocket）
│   └── config/           # 設定管理
├── pkg/                  # 共有パッケージ
├── migrations/           # DBマイグレーション
├── docker/               # Docker関連ファイル
└── firebase/             # Firebase Emulator設定
```

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

### 環境変数

`.env`ファイルを作成して以下を設定：

```env
# Server Configuration
PORT=8080

# Database Configuration
DATABASE_URL=postgres://openlive:openlive@postgres:5432/openlive_db?sslmode=disable

# Firebase Configuration
FIREBASE_PROJECT_ID=openlive-dev
FIREBASE_AUTH_EMULATOR_HOST=firebase:9099

# CORS Configuration
CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:8080
```

### Docker環境での起動

```bash
# コンテナをビルド＆起動
make up

# または
docker-compose up -d

# ログを確認
make logs

# コンテナを停止
make down
```

サーバーは `http://localhost:8080` で起動します。
Firebase Emulator UIは `http://localhost:4000` でアクセスできます。

### ローカル環境での起動

```bash
# 依存関係のインストール
go mod download

# PostgreSQLとFirebase Emulatorを起動
docker-compose up -d postgres firebase

# サーバーを起動
make run-local

# または
go run cmd/server/main.go
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

## WebSocketメッセージフォーマット

### チャットメッセージ

```json
{
  "type": "chat",
  "data": {
    "message": "Hello, World!"
  }
}
```

### WebRTCシグナリング

#### Offer

```json
{
  "type": "offer",
  "to": "user-uuid",
  "data": {
    "sdp": "...",
    "type": "offer"
  }
}
```

#### Answer

```json
{
  "type": "answer",
  "to": "user-uuid",
  "data": {
    "sdp": "...",
    "type": "answer"
  }
}
```

#### ICE Candidate

```json
{
  "type": "ice-candidate",
  "to": "user-uuid",
  "data": {
    "candidate": "...",
    "sdpMid": "...",
    "sdpMLineIndex": 0
  }
}
```

## データベース

PostgreSQL 16を使用しています。マイグレーションファイルは `migrations/001_init.sql` にあります。

### テーブル構成

- `users` - ユーザー情報
- `live_areas` - ライブエリア（配信広場）
- `live_area_members` - エリアメンバー
- `live_streams` - ライブ配信
- `live_stream_viewers` - 配信視聴者

## 開発

### ホットリロード

開発時はAirを使用してホットリロードが有効になっています：

```bash
docker-compose up
```

### テスト

```bash
make test
```

## アーキテクチャ

このプロジェクトはDDD（ドメイン駆動設計）とクリーンアーキテクチャの原則に従っています：

1. **ドメイン層**: ビジネスロジックとエンティティ
2. **ユースケース層**: アプリケーションのビジネスルール
3. **インターフェース層**: HTTP/WebSocketハンドラー、DTO
4. **インフラストラクチャ層**: データベース、外部サービス

## ライセンス

MIT
