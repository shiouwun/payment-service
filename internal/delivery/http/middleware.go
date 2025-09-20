package http

import (
	"net/http"
	"strings"

	"github.com/company/payment-service/internal/domain/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	merchantRepo repository.MerchantRepository
}

func NewAuthMiddleware(merchantRepo repository.MerchantRepository) *AuthMiddleware {
	return &AuthMiddleware{
		merchantRepo: merchantRepo,
	}
}

func (m *AuthMiddleware) APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// 嘗試從 Authorization header 獲取
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "API key is required",
			})
			c.Abort()
			return
		}

		merchant, err := m.merchantRepo.GetByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid API key",
			})
			c.Abort()
			return
		}

		if !merchant.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Merchant account is inactive",
			})
			c.Abort()
			return
		}

		// 將商戶信息存儲在上下文中
		c.Set("merchant", merchant)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-Key, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

func generateRequestID() string {
	// 簡單的請求ID生成，實際應用中可能需要更復雜的實現
	return "req_" + strings.Replace(uuid.New().String(), "-", "", -1)[:16]
}