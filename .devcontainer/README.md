# Dev Container 開発環境

このディレクトリには、VS Code Dev Containersを使用した開発環境の設定が含まれています。

## 使い方

### 1. 前提条件

- Docker Desktop がインストールされている
- VS Code がインストールされている
- VS Code拡張機能「Dev Containers」がインストールされている

### 2. Dev Containerで開く

**⚠️ 重要: 正しいディレクトリから開いてください！**

```bash
# 正しい方法
cd /Users/taiseihayashizaki/dev_private/openLive
code .

# または、ワークスペースファイルを開く
code openLive.code-workspace
```

1. VS Codeで**openLiveプロジェクトディレクトリ**を開く
2. コマンドパレット（`Cmd+Shift+P` または `Ctrl+Shift+P`）を開く
3. `Dev Containers: Reopen in Container` を選択

**エラーが出る場合**: VS Codeが親ディレクトリから開かれている可能性があります。一度閉じて、`openLive`ディレクトリから開き直してください。

### 3. 初回セットアップ

コンテナが起動したら、以下を実行：

```bash
# Server側の依存関係をインストール
cd /workspace/server
go mod download

# Frontend側の依存関係をインストール
cd /workspace/frontend
npm install
```

### 4. サービスの起動

Dev Container内のターミナルで：

```bash
# インフラを起動
make infra

# Serverを起動
make server

# Frontendを起動
make frontend

# または全部まとめて
make all
```

### 5. アクセスURL

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Firebase Emulator UI: http://localhost:4000
- PostgreSQL: localhost:5432

## 構成

### devcontainer.json

- Go 1.21 + Node.js 18 の開発環境
- 必要なVS Code拡張機能が自動インストール
- ポートフォワーディング設定済み
- PostgreSQL、Firebase Emulatorと自動接続

### インストール済みツール

**Go関連:**
- Air (ホットリロード)
- gopls (言語サーバー)
- dlv (デバッガー)
- staticcheck (リンター)

**Node.js関連:**
- npm
- Firebase CLI

**その他:**
- Git
- PostgreSQL Client
- vim/nano

## トラブルシューティング

### ポートが使用中の場合

```bash
# ローカルのコンテナを停止
docker-compose down
```

### ネットワークエラーの場合

```bash
# ネットワークを作り直す
docker network rm openlive_openlive-network
make infra
```

### 依存関係の問題

```bash
# Server
cd /workspace/server
go mod tidy

# Frontend
cd /workspace/frontend
rm -rf node_modules package-lock.json
npm install
```

## 開発のヒント

### デバッグ

- Go: VS Codeのデバッガーでブレークポイントを設定可能
- Next.js: ブラウザのDevToolsを使用

### ログ確認

```bash
make logs-server    # Serverログ
make logs-frontend  # Frontendログ
make logs           # 全ログ
```

### DB接続

```bash
make db  # PostgreSQLに接続
```

### コンテナの再起動

```bash
make restart-server   # Serverのみ
make restart-frontend # Frontendのみ
make restart          # 全サービス
```

