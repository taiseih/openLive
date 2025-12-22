# 開発用コンテナ（Go + Node.js）
FROM golang:1.21-bullseye

# Node.jsのインストール
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - \
    && apt-get update \
    && apt-get install -y nodejs \
    && apt-get install -y \
        git \
        vim \
        nano \
        curl \
        wget \
        zsh \
        postgresql-client \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# 非rootユーザーの作成
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && apt-get update \
    && apt-get install -y sudo \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# Goツールのインストール
RUN go install github.com/cosmtrek/air@v1.49.0 \
    && go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest

# グローバルnpmパッケージのインストール
RUN npm install -g npm@latest \
    && npm install -g firebase-tools

USER $USERNAME

WORKDIR /workspace

