FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY api/go.mod api/go.sum ./

COPY bfs/ ../bfs/
COPY dfs/ ../dfs/
COPY scraping/ ../scraping/

RUN go mod download

COPY api/ .

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

COPY bfs/ ./bfs/
COPY dfs/ ./dfs/
COPY scraping/ ./scraping/

EXPOSE 8080

CMD ["./main"]
