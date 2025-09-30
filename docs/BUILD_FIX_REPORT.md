# Payment Service 建置修正報告

## 概述
本報告詳細說明在建置 payment-service 專案過程中遇到的語法錯誤及其修正方法。

## 修正項目

### 1. Logger 變數名稱衝突 (pkg/logger/logger.go:97)

#### 問題描述
```go
// 錯誤代碼
core := zapcore.NewCore(encoder, writeSyncer, level)
zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

return &zapLogger{zap: zapLogger}, nil
```

#### 錯誤信息
```
zapLogger is not a type
```

#### 問題分析
- **衝突原因**: 變數名稱 `zapLogger` 與結構體類型 `zapLogger` 同名
- **語法問題**: Go 編譯器無法區分是指變數還是類型
- **作用域影響**: 在同一作用域內，標識符不能同時作為變數名和類型名使用

#### 修正方法
```go
// 正確代碼
core := zapcore.NewCore(encoder, writeSyncer, level)
logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

return &zapLogger{zap: logger}, nil
```

#### 修正詳解
- **變數重命名**: 將變數 `zapLogger` 重命名為 `logger`
- **避免衝突**: 消除了變數名與類型名的命名衝突
- **保持語義**: 變數名仍然清晰表達其用途
- **類型安全**: 確保 `&zapLogger{zap: logger}` 能正確初始化結構體

### 2. Logger 方法調用錯誤 (cmd/server/main.go)

#### 問題描述
```go
// 錯誤代碼
logger.Fatal("Failed to connect to database", logger.Error(err))
logger.Fatal("Failed to start server", logger.Error(err))
logger.Fatal("Server forced to shutdown", logger.Error(err))
```

#### 錯誤信息
```
logger.Error(err) (no value) used as value
cannot use err (variable of interface type error) as string value in argument to logger.Error
```

#### 問題分析
- **方法簽名錯誤**: `logger.Error()` 是 void 方法，不返回值
- **參數類型錯誤**: `logger.Error()` 期望字符串消息作為第一個參數
- **鏈式調用錯誤**: 試圖將無返回值的方法調用作為參數傳遞

#### Logger 介面分析
```go
type Logger interface {
    Debug(msg string, fields ...zap.Field)
    Info(msg string, fields ...zap.Field)
    Warn(msg string, fields ...zap.Field)
    Error(msg string, fields ...zap.Field)
    Fatal(msg string, fields ...zap.Field)
    With(fields ...zap.Field) Logger
}
```

#### 修正方法
```go
// 正確代碼
logger.Fatal("Failed to connect to database", zap.Error(err))
logger.Fatal("Failed to start server", zap.Error(err))
logger.Fatal("Server forced to shutdown", zap.Error(err))
```

#### 修正詳解
- **使用 zap.Error()**: 這是 zap 包提供的字段構造函數
- **返回 zap.Field**: `zap.Error(err)` 返回 `zap.Field` 類型
- **正確參數傳遞**: 符合 Logger 介面的方法簽名要求
- **錯誤處理**: 正確封裝錯誤信息為結構化日誌字段

### 3. 缺少 Import 聲明 (cmd/server/main.go)

#### 問題描述
```go
// 缺少 import
import (
    // ... 其他 imports
    "github.com/company/payment-service/pkg/logger"
    "github.com/joho/godotenv"
    // 缺少 "go.uber.org/zap"
)
```

#### 錯誤信息
```
undefined: zap
```

#### 問題分析
- **未導入包**: 使用了 `zap.Error()` 但未導入 zap 包
- **編譯器無法解析**: 無法找到 `zap` 標識符的定義
- **依賴關係**: 雖然間接依賴 zap，但需要顯式導入才能直接使用

#### 修正方法
```go
// 正確代碼
import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    httpdelivery "github.com/company/payment-service/internal/delivery/http"
    "github.com/company/payment-service/internal/domain/usecase"
    "github.com/company/payment-service/internal/infrastructure/config"
    "github.com/company/payment-service/internal/infrastructure/database"
    "github.com/company/payment-service/pkg/logger"
    "github.com/joho/godotenv"
    "go.uber.org/zap"  // 新增的 import
)
```

#### 修正詳解
- **顯式導入**: 添加 `"go.uber.org/zap"` 到 import 列表
- **包可見性**: 使 `zap` 包的導出標識符在當前包中可見
- **依賴解析**: 編譯器能正確解析 `zap.Error()` 函數調用
- **模組管理**: 確保 go.mod 中的依賴能被正確使用

## 技術細節

### Zap Logger 架構
```go
// 核心組件關係
zapcore.Core -> zap.Logger -> zapLogger (custom wrapper) -> Logger (interface)
```

### 錯誤處理最佳實踐
```go
// 推薦方式：結構化日誌
logger.Error("operation failed",
    zap.Error(err),
    zap.String("operation", "database_connection"),
    zap.Duration("timeout", 30*time.Second))

// 避免方式：字符串拼接
logger.Error("operation failed: " + err.Error())
```

### Go 語言關鍵概念

#### 1. 標識符作用域
- 同一作用域內標識符必須唯一
- 類型名和變數名不能重複
- 包級別標識符對整個包可見

#### 2. 方法簽名匹配
- 介面實現必須完全匹配方法簽名
- 參數類型和返回值類型必須一致
- 可變參數 `...T` 的正確使用

#### 3. 包導入規則
- 必須顯式導入使用的包
- 間接依賴不會自動可見
- Import 路徑必須與 go.mod 中的模組路徑一致

## 建置結果

### 修正前
```
# github.com/company/payment-service/pkg/logger
pkg\logger\logger.go:97:10: zapLogger is not a type

# github.com/company/payment-service/cmd/server
cmd\server\main.go:55:49: logger.Error(err) (no value) used as value
cmd\server\main.go:55:62: cannot use err (variable of interface type error) as string value in argument to logger.Error
cmd\server\main.go:82:43: logger.Error(err) (no value) used as value
cmd\server\main.go:82:56: cannot use err (variable of interface type error) as string value in argument to logger.Error
cmd\server\main.go:97:45: logger.Error(err) (no value) used as value
cmd\server\main.go:97:58: cannot use err (variable of interface type error) as string value in argument to logger.Error

cmd\server\main.go:55:49: undefined: zap
cmd\server\main.go:82:43: undefined: zap
cmd\server\main.go:97:45: undefined: zap
```

### 修正後
```
$ go build -o payment-service -v ./cmd/server
github.com/company/payment-service/cmd/server

$ ls -la payment-service
-rwxr-xr-x 1 shiou 197121 15378944 九月 20 22:17 payment-service
```

## 總結

本次修正解決了三個關鍵的 Go 語言編譯錯誤：

1. **命名衝突**: 透過重命名變數避免與類型名衝突
2. **方法調用**: 正確使用 zap 日誌庫的 API
3. **包導入**: 確保所有使用的包都被正確導入

這些修正不僅解決了編譯問題，還提升了代碼的可讀性和維護性，遵循了 Go 語言的最佳實踐。

## 建議

### 開發最佳實踐
1. **命名規範**: 避免變數名與類型名衝突
2. **API 使用**: 仔細閱讀第三方庫的文檔和方法簽名
3. **依賴管理**: 明確導入所有直接使用的包
4. **錯誤處理**: 使用結構化日誌記錄錯誤信息

### 工具建議
1. 使用 `go vet` 進行靜態分析
2. 配置 IDE 的語法檢查和自動導入
3. 定期運行 `go mod tidy` 清理依賴
4. 使用 linter 工具如 `golangci-lint` 提升代碼質量