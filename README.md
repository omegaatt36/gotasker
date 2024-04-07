# GoTasker

GoTasker 是一個使用 Go 語言開發 RESTful API，提供簡單而高效的任務管理功能。它允許用戶創建、檢視、更新和刪除任務，每個任務包含名稱和完成狀態。這個 API 特別適合需要快速、可靠任務追蹤的場景。GoTasker 的設計考慮了易用性和靈活性，並通過單元測試和 Docker 容器化來確保穩定性和可擴展性。

## Requirements

- Golang 1.22+
- Docker or Podman

## Configuration

## How To Use

提供兩種方法，主要差異在 redis(in-memory data storage) 的持久與否。

### 方法一

使用外部 redis 作為持久資料儲存庫，所有 CRUD 過的資料會被儲存在

1. 簡單使用 `make setup-redis` 來啟動一個持久的 redis container，假裝是外部的資料層
1. 啟動服務
    - 選項一：直接使用 `go run main.go` 來啟動服務，參數可以參考 [#Configuration](#Configuration) 段落。
    - 選項二：透過 `docker build -t gotasker:latest .` 來打包成 image，再執行 `docker run --net=host gotasker:latest` 來啟動 container。
1. 透過 `make setup-swagger` 來啟動 swagger，或是直接使用 curl/httpie 等 http client 來 call endpoint。
1. 再測試完成後可以使用 `make remove` 來刪除持久資料。

### 方法二

可以直接使用 `docker compose up` 來透過 `docker-compose.yaml` 與 `Dockerfile` 直接將 stack 跑起來，需要注意的是 redis 有做 health check，以及 api 會等待 redis 是健康的才啟動，請看到以下 log 再進行使用：

```shell
gotasker-api      | INFO        api/server.go:64        starts serving...
```

若是更改了程式碼，需要重新編譯，請使用 `docker compose up --build` 而非 `docker compose up`，如此一來 docker 才會重新拿 Dockerfile 來再次打包。

## Troubleshooting

特別需要注意的是，無論使用何種方法，看到以下 log 才是成功啟動。

```go
INFO    api/server.go:64        starts serving...
```

此程式僅有 redis 會依賴外部服務，若是看到以下 log，代表啟動 redis 的方式錯誤，或遠端 redis 並未開啟防火牆等等，請檢查 redis 的連線設定。

```go
❯ go run main.go
PANIC   database/database.go:29 connect to redis(localhost:6379) failed: dial tcp [::1]:6379: connect: connection refused
panic: connect to redis(localhost:6379) failed: dial tcp [::1]:6379: connect: connection refused
```