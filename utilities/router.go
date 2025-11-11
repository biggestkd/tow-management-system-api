package utilities

import (
	"net/http"
	"tow-management-system-api/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	userHandler    *handler.UserHandler
	companyHandler *handler.CompanyHandler
	towHandler     *handler.TowHandler
	metricHandler  *handler.MetricHandler
	priceHandler   *handler.PriceHandler
	paymentHandler *handler.PaymentHandler
}

func NewRouter(user *handler.UserHandler, company *handler.CompanyHandler, towHandler *handler.TowHandler, metricHandler *handler.MetricHandler, priceHandler *handler.PriceHandler, paymentHandler *handler.PaymentHandler) *Router {
	return &Router{
		userHandler:    user,
		companyHandler: company,
		towHandler:     towHandler,
		metricHandler:  metricHandler,
		priceHandler:   priceHandler,
		paymentHandler: paymentHandler,
	}
}

// InitializeRouter builds the gin.Engine and registers routes/middleware.
func (r *Router) InitializeRouter() *gin.Engine {

	engine := gin.Default()

	// Set CORs policy
	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-User-Id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
	})

	// Health
	engine.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// ==== User routes ====
	engine.POST("/user", r.userHandler.PostUser)       // Create a user
	engine.GET("/user/:userId", r.userHandler.GetUser) // Get a user
	engine.PUT("/user/:userId", r.userHandler.PutUser) // Update a user

	// ==== Company routes ====
	engine.POST("/company", r.companyHandler.PostCompany)   // Create a company
	engine.GET("/company/:id", r.companyHandler.GetCompany) // Get a company
	engine.PUT("/company/:id", r.companyHandler.PutCompany) // Update a company

	// ==== Tow routes ====
	engine.GET("/tows/company/:companyId", r.towHandler.GetTowHistory) // Get tow history
	engine.POST("/tows/:companyId", r.towHandler.PostTow)              // Create tow
	engine.PUT("/tows/:towId", r.towHandler.PutUpdateTow)              // Update tow

	// ==== Metric routes ====
	engine.GET("/metrics/:companyId", r.metricHandler.GetCompanyMetrics) // Get metrics

	// ==== Price routes ====
	engine.GET("/pricing/company/:companyId", r.priceHandler.GetPrices) // Get prices by company
	engine.PUT("/pricing", r.priceHandler.PutPrices)                    // Set prices

	// ==== Payment routes ====
	engine.POST("/payment/webhook", r.paymentHandler.PostStripeWebhook) // Handle payment webhooks

	return engine
}
