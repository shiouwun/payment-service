# Go 企業級金流服務 - 本機開發練習指南

這是專為本機開發設計的練習指南，不依賴 Docker，讓您專注於 Go 語言學習。

## 📚 第一階段：Go 環境準備

### 步驟 1: 安裝 Go
```bash
# 檢查是否已安裝 Go
go version

# 如果沒有安裝，請到官網下載：
# https://golang.org/dl/
# 下載 go1.21.x.windows-amd64.msi 並安裝
```

### 步驟 2: 驗證 Go 環境
```bash
# 檢查 Go 路徑
go env GOPATH
go env GOROOT

# 檢查模組支援
go version
```

### 步驟 3: 安裝 PostgreSQL (本機版)
有幾種選擇：

**選項 A: 安裝完整版 PostgreSQL**
1. 下載：https://www.postgresql.org/download/windows/
2. 安裝時記住密碼（建議設為 `postgres`）
3. 預設埠號：5432

**選項 B: 使用 SQLite (更簡單)**
我們可以先用 SQLite 練習，不需要安裝 PostgreSQL

### 步驟 4: 進入專案並初始化
```bash
cd payment-service

# 初始化模組（如果還沒做）
go mod init github.com/company/payment-service

# 下載依賴
go mod download
go mod tidy
```

---

## 🔧 第二階段：修改為 SQLite 版本（簡化開發）

### 步驟 5: 修改 go.mod 支援 SQLite
```bash
# 添加 SQLite 驅動
go get github.com/mattn/go-sqlite3
```

### 步驟 6: 創建 SQLite 版本的資料庫設定
```bash
# 創建本機資料庫目錄
mkdir -p data
```

### 步驟 7: 修改配置支援 SQLite
編輯 `configs/config.yaml`：

```yaml
server:
  host: "localhost"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"

database:
  driver: "sqlite"
  dsn: "./data/payment.db"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: "5m"

logger:
  level: "info"
  format: "console"  # 開發環境用 console 更易讀
  output_path: "stdout"

app:
  name: "payment-service"
  version: "1.0.0"
  environment: "development"
```

### 步驟 8: 創建 SQLite 初始化腳本
創建 `scripts/init_sqlite.sql`：

```sql
-- SQLite 版本的初始化腳本
PRAGMA foreign_keys = ON;

-- 創建商戶表
CREATE TABLE IF NOT EXISTS merchants (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    api_key TEXT UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 創建客戶表
CREATE TABLE IF NOT EXISTS customers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 創建支付表
CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY,
    merchant_id TEXT NOT NULL,
    customer_id TEXT NOT NULL,
    amount INTEGER NOT NULL,
    currency TEXT NOT NULL DEFAULT 'USD',
    method TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    description TEXT,
    reference TEXT UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (merchant_id) REFERENCES merchants(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

-- 創建索引
CREATE INDEX IF NOT EXISTS idx_payments_merchant_id ON payments(merchant_id);
CREATE INDEX IF NOT EXISTS idx_payments_customer_id ON payments(customer_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at);

-- 插入測試資料
INSERT OR IGNORE INTO merchants (id, name, email, api_key, is_active) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Test Merchant 1', 'merchant1@example.com', 'api_key_merchant_1', 1),
    ('550e8400-e29b-41d4-a716-446655440002', 'Test Merchant 2', 'merchant2@example.com', 'api_key_merchant_2', 1);

INSERT OR IGNORE INTO customers (id, name, email, phone) VALUES
    ('550e8400-e29b-41d4-a716-446655440101', 'John Doe', 'john@example.com', '+1234567890'),
    ('550e8400-e29b-41d4-a716-446655440102', 'Jane Smith', 'jane@example.com', '+1234567891');
```

---

## 🎯 第三階段：本機運行與測試

### 步驟 9: 初始化資料庫
```bash
# 安裝 SQLite 命令列工具（如果沒有）
# Windows: 下載 sqlite-tools 從 https://sqlite.org/download.html

# 初始化資料庫
sqlite3 data/payment.db < scripts/init_sqlite.sql

# 檢查表格是否創建成功
sqlite3 data/payment.db ".tables"
```

### 步驟 10: 編譯並運行服務
```bash
# 編譯
go build -o payment-service.exe cmd/server/main.go

# 運行
./payment-service.exe

# 或直接運行（不編譯）
go run cmd/server/main.go
```

### 步驟 11: 測試健康檢查
打開新的命令列視窗：

```bash
# 測試健康檢查（使用 curl 或瀏覽器）
curl http://localhost:8080/health

# 如果沒有 curl，可以在瀏覽器打開：
# http://localhost:8080/health
```

---

## 💻 第四階段：API 測試實作

### 步驟 12: 使用 PowerShell 測試 API

**創建支付**：
```powershell
# PowerShell 版本的 API 測試
$headers = @{
    'Content-Type' = 'application/json'
    'X-API-Key' = 'api_key_merchant_1'
}

$body = @{
    merchant_id = '550e8400-e29b-41d4-a716-446655440001'
    customer_id = '550e8400-e29b-41d4-a716-446655440101'
    amount = 10000
    currency = 'USD'
    method = 'credit_card'
    description = '我的第一筆測試支付'
    reference = 'PRACTICE_001'
} | ConvertTo-Json

Invoke-RestMethod -Uri 'http://localhost:8080/api/v1/payments' -Method Post -Headers $headers -Body $body
```

### 步驟 13: 使用 Postman 或瀏覽器測試

**Postman 設定**：
1. 下載 Postman: https://www.postman.com/downloads/
2. 創建新的 Collection: "Payment Service"
3. 添加請求：

```
POST http://localhost:8080/api/v1/payments
Headers:
  Content-Type: application/json
  X-API-Key: api_key_merchant_1

Body (JSON):
{
  "merchant_id": "550e8400-e29b-41d4-a716-446655440001",
  "customer_id": "550e8400-e29b-41d4-a716-446655440101",
  "amount": 10000,
  "currency": "USD",
  "method": "credit_card",
  "description": "我的第一筆測試支付",
  "reference": "PRACTICE_001"
}
```

---

## 🧪 第五階段：程式碼學習與測試

### 步驟 14: 運行單元測試
```bash
# 運行所有測試
go test ./...

# 運行特定包的測試並顯示詳細資訊
go test -v ./internal/domain/usecase

# 生成測試覆蓋率
go test -cover ./...
```

### 步驟 15: 理解測試程式碼
打開 `internal/domain/usecase/payment_usecase_test.go`，學習：

1. **Mock 的使用**：
```go
type MockPaymentRepository struct {
    mock.Mock
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
    args := m.Called(ctx, payment)
    return args.Error(0)
}
```

2. **測試案例結構**：
```go
tests := []struct {
    name          string
    request       CreatePaymentRequest
    setupMocks    func(...)
    expectedError string
}{
    // 測試案例...
}
```

**🎯 練習任務 1**: 添加一個新的測試案例

---

## 📊 第六階段：查看和操作資料

### 步驟 16: 使用 SQLite 瀏覽器
推薦工具：
- DB Browser for SQLite: https://sqlitebrowser.org/
- 或使用命令列：

```bash
# 進入 SQLite 互動模式
sqlite3 data/payment.db

# 查看所有支付記錄
.mode column
.headers on
SELECT * FROM payments;

# 查看商戶資料
SELECT * FROM merchants;

# 離開
.quit
```

### 步驟 17: 分析資料庫設計
觀察表格關係：

```sql
-- 查看外鍵約束
.schema payments

-- 統計數據
SELECT status, COUNT(*) as count FROM payments GROUP BY status;

-- 查看支付歷史
SELECT
    p.id,
    m.name as merchant_name,
    c.name as customer_name,
    p.amount/100.0 as amount_dollars,
    p.status,
    p.created_at
FROM payments p
JOIN merchants m ON p.merchant_id = m.id
JOIN customers c ON p.customer_id = c.id
ORDER BY p.created_at DESC;
```

---

## 🛠️ 第七階段：實作新功能

### 步驟 18: 實現客戶管理 API
**🎯 練習任務 2**: 添加客戶相關的 HTTP 端點

1. 創建 `internal/delivery/http/customer_handler.go`
2. 實現以下端點：
   - `POST /api/v1/customers` - 創建客戶
   - `GET /api/v1/customers/{id}` - 查詢客戶
   - `PUT /api/v1/customers/{id}` - 更新客戶

### 步驟 19: 添加支付統計功能
在 `PaymentUseCase` 中添加：

```go
type PaymentStats struct {
    TotalPayments    int64   `json:"total_payments"`
    TotalAmount      int64   `json:"total_amount"`
    CompletedCount   int64   `json:"completed_count"`
    PendingCount     int64   `json:"pending_count"`
    AverageAmount    float64 `json:"average_amount"`
}
```

**🎯 練習任務 3**: 實現統計功能

---

## 📈 第八階段：性能測試與優化

### 步驟 20: 簡單的性能測試
```bash
# 安裝 hey 工具（Go 寫的壓力測試工具）
go install github.com/rakyll/hey@latest

# 測試健康檢查端點
hey -n 100 -c 5 http://localhost:8080/health

# 測試 API（需要先有有效的 payment ID）
hey -n 50 -c 2 -H "X-API-Key: api_key_merchant_1" http://localhost:8080/api/v1/payments/{payment_id}
```

### 步驟 21: 添加日誌分析
在程式碼中添加更多日誌：

```go
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        log.Printf("CreatePayment took %v", duration)
    }()

    // 處理邏輯...
}
```

---

## 🎯 進階挑戰（選做）

### 挑戰 1: 實現配置熱重載
監控配置檔案變化，自動重新載入設定

### 挑戰 2: 添加快取層
使用 Go 的 sync.Map 或第三方套件實現記憶體快取

### 挑戰 3: 實現優雅關閉
處理 Ctrl+C 信號，確保正在處理的請求完成後才關閉服務

```go
// 在 main.go 中添加
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)

go func() {
    <-c
    log.Println("Shutting down gracefully...")
    // 實現優雅關閉邏輯
    os.Exit(0)
}()
```

---

## 📝 學習檢查清單

### Go 基礎 ✅
- [ ] 理解 Go 模組系統
- [ ] 掌握結構體和介面
- [ ] 了解錯誤處理模式
- [ ] 熟悉 Go 的併發基礎

### Web 開發 ✅
- [ ] 理解 HTTP 路由
- [ ] 掌握 JSON 處理
- [ ] 了解中間件概念
- [ ] 熟悉 RESTful API 設計

### 資料庫操作 ✅
- [ ] 理解 SQL 基礎
- [ ] 掌握 Go 資料庫操作
- [ ] 了解交易處理
- [ ] 熟悉資料庫設計原則

### 企業級開發 ✅
- [ ] 理解分層架構
- [ ] 掌握依賴注入
- [ ] 了解測試策略
- [ ] 熟悉配置管理

---

## 🚀 下一步學習建議

1. **深入 Go 語言**：
   - 學習 Goroutines 和 Channels
   - 理解 Go 的記憶體模型
   - 掌握 Go 工具鏈

2. **Web 框架精通**：
   - 深入 Gin 框架特性
   - 學習中間件開發
   - 了解效能調優

3. **資料庫進階**：
   - 學習 PostgreSQL 特性
   - 理解資料庫效能調優
   - 掌握資料庫遷移策略

4. **部署與運維**：
   - 學習 Docker 容器化
   - 了解 CI/CD 流程
   - 掌握監控和日誌

---

## 💡 開發技巧

### 1. 使用 Go 工具
```bash
# 格式化程式碼
go fmt ./...

# 檢查程式碼
go vet ./...

# 生成文檔
go doc -http=:6060
```

### 2. 偵錯技巧
```bash
# 編譯時加入偵錯資訊
go build -gcflags="-N -l" -o debug-service cmd/server/main.go

# 使用 delve 偵錯器
go install github.com/go-delve/delve/cmd/dlv@latest
dlv exec ./debug-service
```

### 3. 效能分析
```bash
# 添加 pprof 支援到 main.go
import _ "net/http/pprof"

# 在程式中啟動 pprof 伺服器
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

---

**開始您的 Go 企業級開發學習之旅！**
記住：每個階段都很重要，不要急於求成，紮實地完成每一步 🎯