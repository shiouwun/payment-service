# Payment Service 專案結構

**最後更新**: 2025-09-30

## 📂 專案目錄結構

```
payment-service/
├── .claude/                    # Claude Code 設定
│   └── settings.local.json     # 本地權限設定
├── cmd/                        # 應用程式入口
│   └── server/
│       └── main.go            # 主程式
├── configs/                    # 配置檔案
│   └── config.yaml            # 應用程式設定（含預設值）
├── docs/                       # 📚 專案文件
│   ├── README.md              # 文件導覽
│   ├── BUILD_FIX_REPORT.md    # 建置問題修復記錄
│   ├── CLONE_BUILD_VERIFICATION.md  # Clone 建置驗證
│   ├── PRACTICE.md            # 開發實踐記錄
│   ├── PRACTICE_LOCAL.md      # 本地開發記錄
│   ├── SECURITY_CHECK.md      # 安全檢查報告
│   └── TEST_REPORT.md         # API 測試報告
├── internal/                   # 內部套件（不對外開放）
│   ├── delivery/              # 交付層
│   │   └── http/              # HTTP 處理器
│   │       ├── middleware.go  # 中介軟體（認證、CORS、日誌）
│   │       ├── payment_handler.go  # 支付 API 處理器
│   │       └── router.go      # 路由設定
│   ├── domain/                # 領域層
│   │   ├── entity/            # 實體定義
│   │   │   └── payment.go     # Payment, Merchant, Customer
│   │   ├── repository/        # Repository 介面
│   │   │   └── payment_repository.go
│   │   └── usecase/           # 業務邏輯
│   │       ├── payment_usecase.go
│   │       └── payment_usecase_test.go
│   └── infrastructure/        # 基礎設施層
│       ├── config/            # 配置管理
│       │   └── config.go      # Viper 配置載入
│       └── database/          # 資料庫實作
│           ├── postgres.go    # PostgreSQL 連接
│           ├── payment_repository_impl.go
│           ├── merchant_repository_impl.go
│           └── customer_repository_impl.go
├── pkg/                        # 可重用套件
│   ├── errors/                # 錯誤處理
│   │   └── errors.go
│   └── logger/                # 日誌
│       └── logger.go          # Zap logger 封裝
├── scripts/                    # 腳本
│   └── migrations/            # 資料庫遷移
│       └── 001_initial_schema.sql
├── .dockerignore              # Docker 忽略檔案
├── .env.example               # 環境變數範例
├── .gitignore                 # Git 忽略檔案
├── docker-compose.yml         # Docker Compose 設定
├── Dockerfile                 # Docker 建置檔案
├── go.mod                     # Go 模組定義
├── go.sum                     # Go 模組校驗
├── Makefile                   # Make 指令
└── README.md                  # 專案說明文件
```

## 🏗️ 架構分層

### 1. Domain Layer (領域層)
**位置**: `internal/domain/`

**職責**:
- 定義業務實體（Entity）
- 定義 Repository 介面
- 實作業務邏輯（UseCase）

**特點**:
- 不依賴任何外部套件
- 純粹的業務邏輯
- 高可測試性

### 2. Infrastructure Layer (基礎設施層)
**位置**: `internal/infrastructure/`

**職責**:
- 實作 Repository 介面
- 資料庫連接管理
- 外部服務整合

**特點**:
- 實作 Domain 層定義的介面
- 處理技術細節
- 可抽換實作

### 3. Delivery Layer (交付層)
**位置**: `internal/delivery/`

**職責**:
- HTTP API 處理
- 請求驗證
- 回應格式化

**特點**:
- 薄薄一層，僅處理 HTTP 相關邏輯
- 呼叫 UseCase 執行業務邏輯
- 錯誤轉換與回應

### 4. Application Layer (應用層)
**位置**: `cmd/server/`

**職責**:
- 應用程式入口
- 依賴注入組裝
- 啟動與關閉管理

## 📦 重要檔案說明

### 配置相關

#### `configs/config.yaml`
應用程式主要配置檔案，包含：
- Server 設定（host, port, timeout）
- Database 設定（connection info）
- Logger 設定（level, format）
- App 設定（name, version, environment）

#### `.env.example`
環境變數範例檔案，說明可用的環境變數。

**注意**: `.env` 檔案已被 `.gitignore` 忽略。

### 資料庫相關

#### `scripts/migrations/001_initial_schema.sql`
資料庫初始化腳本：
- 建立 merchants, customers, payments 表
- 建立索引
- 建立觸發器（自動更新 updated_at）
- 插入測試資料

### Docker 相關

#### `docker-compose.yml`
定義多個服務：
- **postgres**: PostgreSQL 資料庫
- **payment-service**: Payment Service 應用
- **redis**: Redis 快取（可選）

#### `Dockerfile`
使用 multi-stage build：
1. Builder stage: 編譯 Go 應用
2. Runtime stage: 複製執行檔到最小映像

## 🔧 開發工具檔案

### `Makefile`
提供常用指令：
```makefile
make build      # 建置應用程式
make test       # 執行測試
make run        # 本地運行
make docker     # 建置 Docker 映像
```

### `.claude/settings.local.json`
Claude Code 的權限設定，允許執行特定指令。

## 📚 文件組織

所有專案文件集中在 `docs/` 目錄：

- **開發文件**: BUILD_FIX_REPORT.md, PRACTICE.md
- **測試文件**: TEST_REPORT.md, CLONE_BUILD_VERIFICATION.md
- **安全文件**: SECURITY_CHECK.md

每個文件都記錄了專案開發過程的真實經驗。

## 🚫 被忽略的檔案

### `.gitignore` 涵蓋範圍

**編譯產物**:
- `*.exe`, `*.dll`, `*.so`
- `/payment-service`, `/payment-service_unix`

**開發環境**:
- `.env` (環境變數)
- `.vscode/`, `.idea/` (IDE 設定)

**測試與日誌**:
- `*.log`, `logs/`
- `*.test`, `coverage.out`

**資料庫**:
- `*.db`, `*.sqlite`

## 🎯 專案特點

### 1. Clean Architecture
- 依賴方向明確：外層依賴內層
- 領域層完全獨立
- 高可測試性

### 2. 依賴注入
- 手動組裝依賴（無 DI 容器）
- 透過建構函數注入
- 便於測試與抽換實作

### 3. 配置管理
- 多層級配置（檔案 + 環境變數）
- 環境變數可覆蓋檔案設定
- 內建預設值

### 4. 容器化
- Docker multi-stage build
- Docker Compose 多服務編排
- 開發與生產環境分離

## 📊 程式碼統計

**總行數**: ~2500 行（含測試）

**分布**:
- Domain Layer: ~40%
- Infrastructure Layer: ~35%
- Delivery Layer: ~15%
- Tests: ~10%

**測試覆蓋率**: 43.2% (UseCase 層)

## 🔗 相關連結

- **GitHub**: https://github.com/shiouwun/payment-service
- **文件首頁**: [docs/README.md](README.md)
- **API 文件**: [README.md](../README.md#api-文檔與測試)

---

這個專案結構展現了現代 Go 應用的最佳實踐，適合用於履歷展示和技術面試。