# 🚀 openLive クイックスタートガイド

このガイドに従って、5分で開発環境を立ち上げることができます！

## 方法1: Dev Container（推奨）⭐

**最も簡単で確実な方法！**

### ステップ

1. **前提条件**
   - Docker Desktop をインストール
   - VS Code をインストール
   - VS Code拡張機能「Dev Containers」をインストール

2. **プロジェクトを開く**
   
   **重要**: 必ず`openLive`ディレクトリから開いてください！
   
   ```bash
   # 方法1: ターミナルから（推奨）
   cd /Users/taiseihayashizaki/dev_private/openLive
   code .
   
   # 方法2: ワークスペースファイルを開く
   code /Users/taiseihayashizaki/dev_private/openLive/openLive.code-workspace
   ```
   
   ❌ 間違い: `cd /Users/taiseihayashizaki/dev_private && code .`
   ✅ 正しい: `cd /Users/taiseihayashizaki/dev_private/openLive && code .`

3. **Dev Containerで再起動**
   - `Cmd+Shift+P` (macOS) または `Ctrl+Shift+P` (Windows/Linux)
   - `Dev Containers: Reopen in Container` を選択
   - 初回は5-10分かかります（イメージのビルド）

4. **サービスを起動**
   コンテナ内のターミナルで：
   ```bash
   make all
   ```

5. **アクセス**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Firebase Emulator UI: http://localhost:4000

### メリット

✅ 環境構築が自動
✅ チーム全員が同じ環境
✅ ローカルマシンを汚さない
✅ Go、Node.js、PostgreSQLなど全部コンテナ内
✅ デバッガーも使える
✅ すべての依存関係が自動インストール

---

## 方法2: ローカル + Docker（従来の方法）

### ステップ

1. **前提条件のインストール**
   ```bash
   # macOS (Homebrew)
   brew install go node postgresql docker
   
   # 確認
   go version    # 1.21以上
   node --version # 18以上
   ```

2. **依存関係のインストール**
   ```bash
   # Server
   cd server
   go mod download
   
   # Frontend
   cd ../frontend
   npm install
   ```

3. **Dockerサービスを起動**
   ```bash
   cd ..
   make all
   ```

4. **アクセス**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Firebase Emulator UI: http://localhost:4000

---

## 方法3: Docker Composeのみ

### ステップ

1. **Docker Desktopのみインストール**

2. **全サービスを起動**
   ```bash
   make all
   ```

3. **完了！**

すべてのサービスがコンテナ内で動作します。

---

## よくある質問

### Q: どの方法がおすすめ？

**A: Dev Container（方法1）が最もおすすめです。**

理由：
- 環境構築が自動
- チーム全員が同じ環境
- トラブルが少ない

### Q: ポートが使用中エラーが出る

**A: 既存のコンテナを停止してください**

```bash
make down
# または
docker-compose down
```

### Q: 変更が反映されない

**A: ホットリロードが動作しているか確認**

- Server: Airが自動リロード（コンテナログで確認）
- Frontend: Next.js Dev Serverが自動リロード

手動で再起動する場合：
```bash
make restart-server   # Serverのみ
make restart-frontend # Frontendのみ
```

### Q: データベースに接続できない

**A: PostgreSQLが起動しているか確認**

```bash
make status  # ステータス確認

# 直接接続して確認
make db
```

### Q: Firebase Emulatorが動かない

**A: ポート9099と4000が空いているか確認**

```bash
lsof -i :9099
lsof -i :4000

# 該当プロセスを停止してから再起動
make restart
```

---

## 便利なコマンド

```bash
# ステータス確認
make status

# ログ表示
make logs-server    # Serverのログ
make logs-frontend  # Frontendのログ
make logs           # 全ログ

# 再起動
make restart-server
make restart-frontend
make restart

# 停止
make down

# 完全クリーンアップ
make clean

# DB接続
make db
```

---

## 次のステップ

開発環境が起動したら：

1. **フロントエンドを確認**
   - http://localhost:3000 にアクセス
   - Googleログインを試す

2. **APIを確認**
   - http://localhost:8080/health にアクセス
   - `{"status":"ok"}` が返ってくればOK

3. **Firebase Emulatorを確認**
   - http://localhost:4000 にアクセス
   - Authenticationタブを確認

4. **コードを編集**
   - ファイルを編集すると自動的にリロードされます
   - Server: Go (Air)
   - Frontend: Next.js

5. **デバッグ**
   - VS CodeでF5を押すとGoのデバッガーが起動
   - ブレークポイントを設定して実行

---

## トラブルシューティング

問題が発生した場合：

1. **全コンテナを停止**
   ```bash
   make down
   ```

2. **完全クリーンアップ**
   ```bash
   make clean
   ```

3. **再起動**
   ```bash
   make all
   ```

それでも解決しない場合は、README.mdの詳細なトラブルシューティングセクションを参照してください。

---

## サポート

問題が解決しない場合：
- [README.md](README.md) を確認
- [.devcontainer/README.md](.devcontainer/README.md) を確認
- GitHubのIssuesで質問

Happy Coding! 🎉

