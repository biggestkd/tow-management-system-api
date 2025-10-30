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
}

func NewRouter(user *handler.UserHandler, company *handler.CompanyHandler, towHandler *handler.TowHandler, metricHandler *handler.MetricHandler) *Router {
	return &Router{userHandler: user, companyHandler: company, towHandler: towHandler, metricHandler: metricHandler}
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

	// ==== Tow routes ====
	engine.GET("/tows/company/:companyId", r.towHandler.GetTowHistory) // Get a company

	// ==== Metric routes ====
	engine.GET("/metrics/:companyId", r.metricHandler.GetCompanyMetrics) // Get company metrics

	return engine
}
