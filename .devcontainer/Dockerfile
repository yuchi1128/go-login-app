# .devcontainer/Dockerfile
FROM golang:1.23-alpine

RUN apk update && apk add --no-cache git

RUN git config --global user.name "yuchi1128" && \
    git config --global user.email "yuchi141128@gmail.com"

# Goのツールやリンターなどをインストール
RUN go install golang.org/x/tools/gopls@latest \
    && go install honnef.co/go/tools/cmd/staticcheck@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/ramya-rao-a/go-outline@latest \
    && go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest

# ポートを開放 (docker-compose.ymlでも指定)
EXPOSE 8080

# 非rootユーザーを作成し、利用 (セキュリティ向上のため)
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN addgroup -g $USER_GID $USERNAME \
    && adduser -D -u $USER_UID -h /home/$USERNAME -G $USERNAME $USERNAME

# (Optional) sudo権限を付与する場合
# RUN apk add --no-cache sudo
# RUN echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
#     && chmod 0440 /etc/sudoers.d/$USERNAME

# ワークスペースディレクトリを作成し、所有権を変更
ARG WORKSPACE=/workspace
RUN mkdir -p ${WORKSPACE} \
    && chown -R $USERNAME:$USERNAME ${WORKSPACE}

USER $USERNAME
WORKDIR ${WORKSPACE}

# デフォルトの起動コマンド
CMD ["sleep", "infinity"] 