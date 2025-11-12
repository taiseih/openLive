FROM node:18-alpine

WORKDIR /app

# Firebase CLIをインストール
RUN npm install -g firebase-tools

# 設定ファイルをコピー
COPY docker/firebase/firebase.json .

EXPOSE 9099 4000

CMD ["firebase", "emulators:start", "--only", "auth", "--project", "openlive-dev"]

