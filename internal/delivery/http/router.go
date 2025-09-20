package http

import (
	"github.com/company/payment-service/internal/domain/repository"
	"github.com/company/payment-service/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	paymentUseCase usecase.PaymentUseCase,
	merchantRepo repository.MerchantRepository,
) *gin.Engine {
	// 設置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 中間件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	router.Use(RequestIDMiddleware())

	// 健康檢查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "payment-service",
		})
	})

	// API 路由組
	api := router.Group("/api/v1")

	// 初始化處理器
	paymentHandler := NewPaymentHandler(paymentUseCase)
	authMiddleware := NewAuthMiddleware(merchantRepo)

	// 支付相關路由 - 需要API密鑰驗證
	payments := api.Group("/payments")
	payments.Use(authMiddleware.APIKeyAuth())
	{
		payments.POST("", paymentHandler.CreatePayment)
		payments.GET("/:id", paymentHandler.GetPayment)
		payments.POST("/:id/process", paymentHandler.ProcessPayment)
		payments.POST("/:id/cancel", paymentHandler.CancelPayment)
	}

	// 商戶相關路由
	merchants := api.Group("/merchants")
	merchants.Use(authMiddleware.APIKeyAuth())
	{
		merchants.GET("/:merchantId/payments", paymentHandler.GetMerchantPayments)
	}

	return router
}