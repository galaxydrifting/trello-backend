# Trello-Backend

## 專案簡介
本專案為 Trello 類型的看板系統後端，支援「看板 > 清單 > 卡片」的階層式管理，並具備卡片拖拉、CRUD 操作、用戶認證等功能。

## 技術棧
- 語言：Go 1.23
- Web 框架：Gin
- GraphQL：gqlgen
- ORM：GORM
- 資料庫：PostgreSQL
- Swagger 文件：swaggo
- JWT 驗證

## 快速開始

### 1. 環境需求
- Docker & Docker Compose
- Go 1.23（如需本機開發）

### 2. 啟動服務
```bash
# 啟動後端與資料庫（於專案根目錄）
docker compose -f .devcontainer/docker-compose.yml up -d
```
> 預設會啟動 app（後端）與 db（PostgreSQL）

### 3. 設定環境變數
請參考 `.devcontainer/.env`，主要變數如下：
- `POSTGRES_HOSTNAME`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_PORT`
- `POSTGRES_DB`
- `JWT_SECRET`
- `CORS_ALLOW_ORIGINS`

### 4. 資料庫初始化與部署
- 專案啟動時會自動執行 GORM 的 AutoMigrate，對應程式碼請見 `internal/app/migrations.go`
- 不需手動執行 migration 指令，資料表會自動建立/更新
- 若需自訂 migration，請於 `migrations.go` 增修

### 5. API 文件
- Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- GraphQL Playground: [http://localhost:8080/api/graphql/playground](http://localhost:8080/api/graphql/playground)
- GraphQL CRUD 查詢: [http://localhost:8080/api/graphql/query](http://localhost:8080/api/graphql/query)

## 專案結構
- `cmd/`：主程式入口
- `internal/`：商業邏輯、資料庫、服務、路由
- `graph/`：GraphQL schema 與 resolver
- `docs/`：Swagger 文件與技術文件

## 技術文件
- [Null 安全性：Banishing Null with Non-Nullable References](docs/null-safety.md)
  - 了解如何在 Go 專案中實踐非空引用原則
  - 減少 nil 引用錯誤，提高程式碼品質

## 測試
建議於專案根目錄執行下列指令，獲得詳細測試摘要：

```bash
# 先安裝 gotestsum（如尚未安裝）
go install gotest.tools/gotestsum@latest
# 執行所有測試
gotestsum
```
這會自動搜尋所有子目錄的測試並顯示 summary（總測試數、成功/失敗數）。

如未安裝 gotestsum，可改用：
```bash
go test -v ./...
```

### 安裝 golangci-lint

如尚未安裝，請執行：

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 執行 lint 檢查

於專案根目錄執行：

```bash
golangci-lint run ./...
```

> 建議將此流程納入開發習慣，或於 CI/CD pipeline 自動檢查。