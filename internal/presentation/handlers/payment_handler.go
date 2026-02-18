package handlers

import (
	"net/http"
	"payment-system/internal/application/services"
	"payment-system/internal/presentation/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// ProcessInternalPayment handles POST /api/v1/payments/internal
func (h *PaymentHandler) ProcessInternalPayment(c *gin.Context) {
	var req dto.InternalPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate amount
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than zero"})
		return
	}

	serviceReq := services.InternalPaymentRequest{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Amount:     req.Amount,
		Currency:   req.Currency,
		Reference:  req.Reference,
	}

	payment, err := h.paymentService.ProcessInternalPayment(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PaymentToResponse(payment))
}

// CreateExternalPaymentRequest handles POST /api/v1/payments/external
func (h *PaymentHandler) CreateExternalPaymentRequest(c *gin.Context) {
	var req dto.ExternalPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount format"})
		return
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than zero"})
		return
	}

	// Validate external account details
	if req.ExternalAccountID == nil {
		if req.ExternalAccountNumber == "" && req.ExternalIBAN == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "external account details required"})
			return
		}
	}

	serviceReq := services.ExternalPaymentRequest{
		FromUserID:                req.FromUserID,
		ExternalAccountID:         req.ExternalAccountID,
		Amount:                     amount,
		Currency:                   req.Currency,
		TargetCurrency:             req.TargetCurrency,
		ExternalAccountNumber:      req.ExternalAccountNumber,
		ExternalRoutingNumber:      req.ExternalRoutingNumber,
		ExternalIBAN:               req.ExternalIBAN,
		ExternalSwiftCode:          req.ExternalSwiftCode,
		ExternalBankName:           req.ExternalBankName,
		ExternalAccountHolderName: req.ExternalAccountHolderName,
		ExternalCountryCode:        req.ExternalCountryCode,
		Reference:                  req.Reference,
		IdempotencyKey:             req.IdempotencyKey,
		SaveAccount:                req.SaveAccount,
	}

	paymentRequest, err := h.paymentService.CreateExternalPaymentRequest(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.PaymentRequestToResponse(paymentRequest))
}

// ProcessPaymentRequest handles POST /api/v1/payments/requests/:id/process
func (h *PaymentHandler) ProcessPaymentRequest(c *gin.Context) {
	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	payment, err := h.paymentService.ProcessPaymentRequest(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PaymentToResponse(payment))
}
