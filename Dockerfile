FROM golang:1.25.0

# 作業ディレクトリ
WORKDIR /app

# Go Modules をコピーして依存関係を取得
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションコードをコピー
COPY . .

# Linux 用バイナリにビルド
RUN GOOS=linux GOARCH=amd64 go build -mod=readonly -v -o server

# Cloud Run が使用するポートを EXPOSE
EXPOSE 8080

# 環境変数は Cloud Run 側で設定（Dockerfile では不要）
# CMD では直接バイナリを実行
CMD ["./server"]

# $ docker build -t us-west1-docker.pkg.dev/ginapp-470911/go-api/api-image:latest .
# $ docker push us-west1-docker.pkg.dev/ginapp-470911/go-api/api-image:latest