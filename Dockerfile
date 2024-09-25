FROM golang:1.21 as builder
WORKDIR /go/src/app
ENV GO111MODULE=on
RUN groupadd -g 10001 app \
  && useradd -u 10001 -g app app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip .
RUN ls -la .
RUN chmod 755 ./zoneinfo.zip
RUN make build

# ubuntuイメージ（「--platform=linux/amd64」はM1 Mac環境で実行するために追加）
FROM --platform=linux/amd64 ubuntu:22.04

# 非対話的インストールの設定（Dockerイメージをビルドする際に通常対話的な操作を避ける必要があるため）
ENV DEBIAN_FRONTEND=noninteractive

# リポジトリをクリアしてからアップデート
RUN apt-get clean && apt-get update

# wgetとgnupgをインストール
RUN apt-get install -y wget gnupg

# パッケージをインストール・アップデート
RUN wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | apt-key add -
RUN echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list
  
# 再度リポジトリをアップデートしてGoogle Chromeをインストール
RUN apt-get update
RUN apt-get install -y google-chrome-stable

# "--disable-web-security", "--allow-running-insecure-content", "--disable-features=IsolateOrigins,site-per-process":secureモードを無効にする※RANがhttp://のため
# disable-dev-shm-usage:https://qiita.com/yoshi10321/items/8b7e6ed2c2c15c3344c6
CMD ["google-chrome", "--no-sandbox", "--disable-dev-shm-usage", "--disable-setuid-sandbox", "--no-first-run", "--no-zygote", "--single-process", "--user-data-dir=/data", "--disable-web-security", "--allow-running-insecure-content", "--disable-features=IsolateOrigins,site-per-process"]
ENV GOROOT /usr/local/go
COPY --from=builder /go/bin/app /go/bin/app
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/app/.conf /.conf
COPY --from=builder /go/src/app/.migrations /.migrations
COPY --from=builder /go/src/app/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

ENV APP_ENV=prd
EXPOSE 8080
ENTRYPOINT ["/go/bin/app"]