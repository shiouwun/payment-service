# Go 企業級金流服務 - 實戰練習指南

這是一個循序漸進的練習指南，幫助您深入理解和掌握企業級 Go 開發。

## 📚 第一階段：環境準備與基礎理解

### 步驟 1: 環境檢查
```bash
# 檢查 Go 版本 (需要 1.21+)
go version

# 如果沒有安裝 Go，請到官網下載：https://golang.org/dl/
```

### 步驟 2: 項目初始化
```bash
# 進入專案目錄
cd payment-service

# 初始化 Go 模組
go mod init github.com/company/payment-service

# 安裝依賴
go mod download
go mod tidy
```

### 步驟 3: 理解專案結構
仔細閱讀每個目錄的作用：

```
📁 cmd/server/          # 程式入口點
📁 internal/domain/     # 核心業務邏輯（最重要）
📁 internal/infrastructure/ # 外部依賴實現
📁 internal/delivery/   # API 接口層
📁 pkg/                # 可重用工具包
```

**🎯 練習任務 1**: 用自己的話說明每一層的職責

---

## 📖 第二階段：理解 Clean Architecture

### 步驟 4: 分析領域模型
打開 `internal/domain/entity/payment.go`，理解：

```go
// 仔細觀察這些設計決策：
type Payment struct {
    ID       uuid.UUID     `json:"id" db:"id"`
    Amount   int64         `json:"amount" db:"amount"` // 為什麼用 int64？
    Status   PaymentStatus `json:"status" db:"status"`
    // ...
}
```

**🎯 練習任務 2**:
1. 為什麼金額使用 `int64` 而不是 `float64`？
2. 為什麼使用 `PaymentStatus` 類型而不是 `string`？
3. 為什麼需要 `CreatedAt` 和 `UpdatedAt` 欄位？

### 步驟 5: 理解接口設計
查看 `internal/domain/repository/payment_repository.go`：

```go
type PaymentRepository interface {
    Create(ctx context.Context, payment *entity.Payment) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
    // 為什麼所有方法都接受 context.Context？
}
```

**🎯 練習任務 3**:
1. 為什麼使用接口而不是直接實現？
2. `context.Context` 的作用是什麼？

---

## 🔧 第三階段：實際操作 - 啟動服務

### 步驟 6: 使用 Docker 啟動環境
```bash
# 複製環境變數文件
cp .env.example .env

# 啟動資料庫
docker-compose up -d postgres

# 等待資料庫啟動（約10秒）
sleep 10

# 檢查資料庫狀態
docker-compose ps
```

### 步驟 7: 執行資料庫遷移
```bash
# 手動執行 SQL
docker exec -i payment-postgres psql -U postgres -d payment_service < scripts/migrations/001_initial_schema.sql

# 或使用 make 命令
make db-migrate
```

### 步驟 8: 啟動服務
```bash
# 開發模式啟動
go run cmd/server/main.go

# 或使用 make
make run-dev
```

### 步驟 9: 測試健康檢查
```bash
# 測試服務是否正常
curl http://localhost:8080/health

# 預期回應：
# {"status":"ok","service":"payment-service"}
```

**🎯 練習任務 4**:
服務啟動後，訪問 http://localhost:8080/health，截圖證明服務正常運行

---

## 💻 第四階段：API 實戰操作

### 步驟 10: 創建第一筆支付
```bash
# 使用測試 API Key
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -H "X-API-Key: api_key_merchant_1" \
  -d '{
    "merchant_id": "550e8400-e29b-41d4-a716-446655440001",
    "customer_id": "550e8400-e29b-41d4-a716-446655440101",
    "amount": 10000,
    "currency": "USD",
    "method": "credit_card",
    "description": "我的第一筆測試支付",
    "reference": "PRACTICE_001"
  }'
```

**🎯 練習任務 5**:
記錄回應中的 `payment_id`，我們接下來會用到

### 步驟 11: 查詢支付狀態
```bash
# 替換 {payment_id} 為上一步得到的 ID
curl -H "X-API-Key: api_key_merchant_1" \
  http://localhost:8080/api/v1/payments/{payment_id}
```

### 步驟 12: 處理支付
```bash
# 處理支付
curl -X POST \
  -H "X-API-Key: api_key_merchant_1" \
  http://localhost:8080/api/v1/payments/{payment_id}/process
```

### 步驟 13: 再次查詢確認狀態變化
```bash
curl -H "X-API-Key: api_key_merchant_1" \
  http://localhost:8080/api/v1/payments/{payment_id}
```

**🎯 練習任務 6**:
觀察支付狀態從 "pending" 變為 "completed" 的過程

---

## 🧪 第五階段：測試與程式碼理解

### 步驟 14: 運行單元測試
```bash
# 運行所有測試
go test ./...

# 運行特定測試
go test ./internal/domain/usecase -v

# 生成測試覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 步驟 15: 分析測試程式碼
打開 `internal/domain/usecase/payment_usecase_test.go`，理解：

1. Mock 的使用方式
2. 測試案例的設計
3. 斷言的編寫

**🎯 練習任務 7**:
編寫一個新的測試案例，測試「取消已完成的支付」應該失敗

```go
func TestPaymentUseCase_CancelCompletedPayment(t *testing.T) {
    // 你來實現這個測試
}
```

---

## 🔨 第六階段：實作新功能

### 步驟 16: 添加客戶儲存庫實現
目前 `main.go` 中的 customerRepo 是 nil，讓我們實現它：

```bash
# 創建新文件
touch internal/infrastructure/database/customer_repository_impl.go
```

**🎯 練習任務 8**:
參考 `merchant_repository_impl.go`，實現完整的 `CustomerRepository`

### 步驟 17: 添加支付歷史統計功能
在 `PaymentUseCase` 中添加新方法：

```go
func (uc *paymentUseCase) GetPaymentStatistics(ctx context.Context, merchantID uuid.UUID) (*PaymentStats, error) {
    // 實現統計邏輯
}
```

**🎯 練習任務 9**:
1. 定義 `PaymentStats` 結構體
2. 實現統計方法
3. 添加對應的 HTTP 端點
4. 編寫測試

---

## 📈 第七階段：性能優化與監控

### 步驟 18: 添加日誌記錄
在 `payment_handler.go` 中添加詳細的日誌：

```go
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
    logger := c.MustGet("logger").(logger.Logger)
    logger.Info("Creating payment", zap.String("merchant_id", req.MerchantID.String()))
    // ...
}
```

### 步驟 19: 性能測試
```bash
# 安裝 hey 工具進行壓力測試
go install github.com/rakyll/hey@latest

# 執行壓力測試
hey -n 1000 -c 10 \
  -H "X-API-Key: api_key_merchant_1" \
  http://localhost:8080/health
```

**🎯 練習任務 10**:
分析壓力測試結果，找出性能瓶頸

---

## 🔧 第八階段：進階配置與部署

### 步驟 20: 自定義配置
修改 `configs/config.yaml`：

```yaml
# 添加自定義配置
app:
  max_payment_amount: 1000000  # 最大支付金額
  supported_currencies: ["USD", "EUR", "TWD"]
```

### 步驟 21: Docker 容器化
```bash
# 建構 Docker 映像
docker build -t my-payment-service .

# 使用 Docker Compose 完整部署
docker-compose up -d
```

### 步驟 22: 監控設置
```bash
# 查看服務日誌
docker-compose logs -f payment-service

# 監控資料庫
docker exec -it payment-postgres psql -U postgres -d payment_service -c "SELECT COUNT(*) FROM payments;"
```

**🎯 練習任務 11**:
設置一個簡單的監控指標，統計每分鐘的支付請求數量

---

## 🎯 最終挑戰：完整功能實現

### 挑戰 1: 實現退款功能
1. 設計退款實體和狀態
2. 實現退款業務邏輯
3. 添加 HTTP API
4. 編寫測試
5. 更新文檔

### 挑戰 2: 添加支付方式驗證
1. 為不同支付方式添加驗證邏輯
2. 實現信用卡號碼格式驗證
3. 添加支付限額檢查

### 挑戰 3: 實現資料庫連接池監控
1. 添加連接池指標
2. 實現健康檢查詳細資訊
3. 添加效能監控端點

---

## 📝 學習檢查清單

完成每個階段後，請檢查：

### 基礎理解 ✅
- [ ] 理解 Clean Architecture 各層職責
- [ ] 了解 Repository Pattern 的優勢
- [ ] 掌握 Go 模組管理

### 實際操作 ✅
- [ ] 成功啟動服務
- [ ] 完成基本 API 操作
- [ ] 運行測試通過

### 程式碼實現 ✅
- [ ] 實現新的 Repository
- [ ] 添加新的業務邏輯
- [ ] 編寫有效的測試

### 部署運維 ✅
- [ ] Docker 容器化成功
- [ ] 監控和日誌配置完成
- [ ] 性能測試執行

---

## 🚀 下一步建議

1. **深入學習**: 研究 Go 的併發模式，優化支付處理性能
2. **架構擴展**: 學習微服務架構，拆分單體應用
3. **安全強化**: 實現 OAuth2、JWT 等進階認證機制
4. **監控完善**: 集成 Prometheus、Grafana 等監控工具

## 💡 常見問題解答

**Q: 為什麼使用 Clean Architecture？**
A: 分離關注點，提高可測試性和可維護性，降低技術債務。

**Q: Repository Pattern 的核心價值？**
A: 抽象化資料存取，讓業務邏輯不依賴特定的資料庫實現。

**Q: 為什麼使用 Context？**
A: 傳遞請求範圍的資料、取消信號和超時控制。

---

開始您的企業級 Go 開發之旅吧！記得每完成一個階段就給自己一個讚 👍