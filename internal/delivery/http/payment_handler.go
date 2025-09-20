package http

import (
	"net/http"
	"strconv"

	"github.com/company/payment-service/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentUseCase usecase.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

type CreatePaymentResponse struct {
	Success bool   `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req usecase.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CreatePaymentResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	payment, err := h.paymentUseCase.CreatePayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CreatePaymentResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, CreatePaymentResponse{
		Success: true,
		Data:    payment,
		Message: "Payment created successfully",
	})
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreatePaymentResponse{
			Success: false,
			Error:   "Invalid payment ID format",
		})
		return
	}

	payment, err := h.paymentUseCase.GetPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, CreatePaymentResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CreatePaymentResponse{
		Success: true,
		Data:    payment,
	})
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreatePaymentResponse{
			Success: false,
			Error:   "Invalid payment ID format",
		})
		return
	}

	err = h.paymentUseCase.ProcessPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CreatePaymentResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CreatePaymentResponse{
		Success: true,
		Message: "Payment processed successfully",
	})
}

func (h *PaymentHandler) CancelPayment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreatePaymentResponse{
			Success: false,
			Error:   "Invalid payment ID format",
		})
		return
	}

	err = h.paymentUseCase.CancelPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CreatePaymentResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CreatePaymentResponse{
		Success: true,
		Message: "Payment cancelled successfully",
	})
}

func (h *PaymentHandler) GetMerchantPayments(c *gin.Context) {
	merchantIDParam := c.Param("merchantId")
	merchantID, err := uuid.Parse(merchantIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, CreatePaymentResponse{
			Success: false,
			Error:   "Invalid merchant ID format",
		})
		return
	}

	// 分頁參數
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	payments, err := h.paymentUseCase.GetMerchantPayments(c.Request.Context(), merchantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CreatePaymentResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CreatePaymentResponse{
		Success: true,
		Data:    payments,
	})
}