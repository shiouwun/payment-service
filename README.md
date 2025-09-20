# Payment Service - 企業級 Go 金流服務

一個採用 Clean Architecture 設計的高可維護性金流服務系統，使用 Go 語言開發，提供完整的支付處理功能。

## 🏗️ 專案架構

### Clean Architecture 分層

```
payment-service/
├── cmd/                    # 應用程式入口點
│   └── server/
│       └── main.go        # 主程式
├── internal/              # 內部包（不對外開放）
│   ├── domain/           # 領域層 (Domain Layer)
│   │   ├── entity/       # 實體定義
│   │   ├── repository/   # 資料庫接口定義
│   │   └── usecase/      # 業務邏輯
│   ├── infrastructure/   # 基礎設施層
│   │   ├── database/     # 資料庫實現
│   │   ├── http/         # HTTP 相關
│   │   └── config/       # 配置管理
│   └── delivery/         # 交付層
│       └── http/         # HTTP 處理器
├── pkg/                  # 可重用的包
│   ├── logger/          # 日誌包
│   └── errors/          # 錯誤處理
├── configs/             # 配置文件
├── scripts/            # 腳本文件
│   └── migrations/     # 資料庫遷移
└── docs/              # 文檔
```

### 依賴關係圖

```
Delivery Layer (HTTP Handlers)
      ↓
Use Case Layer (Business Logic)
      ↓
Domain Layer (Entities & Interfaces)
      ↓
Infrastructure Layer (Database, External APIs)
```

## 🚀 快速開始

### 1. 環境準備

確保您已安裝：
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (可選)

### 2. 克隆專案

```bash
git clone <your-repository>
cd payment-service
```

### 3. 安裝依賴

```bash
go mod download
go mod verify
```

### 4. 配置環境

複製環境變數範例檔案：
```bash
cp .env.example .env
```

編輯 `.env` 檔案，設定您的資料庫連接資訊。

### 5. 資料庫設置

#### 使用 Docker Compose (推薦)

```bash
# 啟動所有服務
docker-compose up -d

# 查看服務狀態
docker-compose ps

# 查看日誌
docker-compose logs payment-service
```

#### 手動設置

1. 建立資料庫：
```sql
CREATE DATABASE payment_service;
```

2. 執行遷移：
```bash
psql -U postgres -d payment_service -f scripts/migrations/001_initial_schema.sql
```

### 6. 運行服務

```bash
# 開發模式
go run cmd/server/main.go

# 編譯並運行
go build -o payment-service cmd/server/main.go
./payment-service
```

服務會在 `http://localhost:8080` 上運行。

## 📝 API 文檔

### 健康檢查

```bash
GET /health
```

### 支付相關 API

所有支付 API 都需要 API Key 驗證，請在請求標頭中加入：
```
X-API-Key: your-api-key
```
或
```
Authorization: Bearer your-api-key
```

#### 1. 創建支付

```bash
POST /api/v1/payments
Content-Type: application/json
X-API-Key: api_key_merchant_1

{
  "merchant_id": "550e8400-e29b-41d4-a716-446655440001",
  "customer_id": "550e8400-e29b-41d4-a716-446655440101",
  "amount": 10000,
  "currency": "USD",
  "method": "credit_card",
  "description": "購買商品",
  "reference": "ORDER_001"
}
```

#### 2. 查詢支付

```bash
GET /api/v1/payments/{payment_id}
X-API-Key: api_key_merchant_1
```

#### 3. 處理支付

```bash
POST /api/v1/payments/{payment_id}/process
X-API-Key: api_key_merchant_1
```

#### 4. 取消支付

```bash
POST /api/v1/payments/{payment_id}/cancel
X-API-Key: api_key_merchant_1
```

#### 5. 查詢商戶支付記錄

```bash
GET /api/v1/merchants/{merchant_id}/payments?limit=20&offset=0
X-API-Key: api_key_merchant_1
```

## 🔧 配置管理

### 配置文件

主要配置檔案位於 `configs/config.yaml`：

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "payment_service"
  sslmode: "disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: "5m"
```

### 環境變數

系統支援通過環境變數覆蓋配置，前綴為 `PAYMENT_`：

- `PAYMENT_DATABASE_HOST`
- `PAYMENT_DATABASE_PORT`
- `PAYMENT_SERVER_PORT`
- 等...

## 🧪 測試

### 運行單元測試

```bash
# 運行所有測試
go test ./...

# 運行特定包的測試
go test ./internal/domain/usecase

# 運行測試並顯示覆蓋率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 運行集成測試

```bash
# 確保資料庫已運行
docker-compose up -d postgres

# 運行集成測試
go test -tags=integration ./...
```

## 🐳 Docker 部署

### 建構映像

```bash
docker build -t payment-service:latest .
```

### 使用 Docker Compose

```bash
# 啟動所有服務
docker-compose up -d

# 僅啟動特定服務
docker-compose up -d postgres
docker-compose up -d payment-service

# 停止服務
docker-compose down

# 查看日誌
docker-compose logs -f payment-service
```

## 📊 監控與日誌

### 日誌配置

系統使用 Uber Zap 作為日誌庫，支援：

- **格式**：JSON (生產環境) / Console (開發環境)
- **級別**：debug, info, warn, error, fatal
- **輸出**：stdout 或文件

### 健康檢查

服務提供健康檢查端點：

```bash
curl http://localhost:8080/health
```

回應：
```json
{
  "status": "ok",
  "service": "payment-service"
}
```

## 🔐 安全性

### API 認證

- 使用 API Key 進行身份驗證
- 支援 `X-API-Key` header 或 `Authorization: Bearer` header
- 每個請求都會驗證 merchant 的活躍狀態

### 資料安全

- 密碼和敏感資訊使用環境變數儲存
- API Key 在回應中不會暴露
- 支援 HTTPS (在生產環境中配置)

## 🚧 開發指南

### 添加新功能

1. **領域層**：在 `internal/domain/entity` 定義新實體
2. **資料庫層**：在 `internal/domain/repository` 定義接口
3. **業務邏輯層**：在 `internal/domain/usecase` 實現業務邏輯
4. **基礎設施層**：在 `internal/infrastructure/database` 實現資料庫操作
5. **交付層**：在 `internal/delivery/http` 實現 HTTP 處理器

### 程式碼規範

- 遵循 Go 官方程式碼風格
- 使用 `gofmt` 格式化程式碼
- 使用 `golint` 檢查程式碼品質
- 為所有公開的函數添加文檔註釋

### 提交規範

```bash
# 功能提交
git commit -m "feat: add payment processing functionality"

# 修復提交
git commit -m "fix: resolve database connection issue"

# 文檔提交
git commit -m "docs: update API documentation"
```

## 🔄 CI/CD

### GitHub Actions

專案包含 GitHub Actions 工作流程：

- **測試**：自動運行所有測試
- **建構**：建構 Docker 映像
- **部署**：部署到指定環境

### 本地開發流程

1. 建立功能分支
2. 開發功能
3. 運行測試
4. 提交程式碼
5. 建立 Pull Request

## 📈 性能優化

### 資料庫優化

- 使用適當的索引
- 連接池管理
- 查詢優化

### HTTP 優化

- 請求/回應壓縮
- 連接重用
- 超時設置

## 🛠️ 故障排除

### 常見問題

1. **資料庫連接失敗**
   - 檢查資料庫是否運行
   - 驗證連接字串
   - 檢查防火牆設置

2. **API 認證失敗**
   - 驗證 API Key 正確性
   - 檢查 merchant 狀態
   - 確認 header 格式

3. **服務啟動失敗**
   - 檢查端口是否被占用
   - 驗證配置檔案
   - 查看錯誤日誌

## 📚 相關資源

- [Go 官方文檔](https://golang.org/doc/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [PostgreSQL 文檔](https://www.postgresql.org/docs/)
- [Docker 文檔](https://docs.docker.com/)

## 🤝 貢獻指南

1. Fork 專案
2. 建立功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交變更 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 開啟 Pull Request

## 📄 授權

本專案使用 MIT 授權 - 詳見 [LICENSE](LICENSE) 檔案

## 👥 聯絡方式

- 專案維護者：[Your Name]
- 電子郵件：[your.email@example.com]
- 專案連結：[https://github.com/company/payment-service]