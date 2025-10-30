package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"tow-management-system-api/handler"
	"tow-management-system-api/service"
	"tow-management-system-api/utilities"
)

var ginLambda *ginadapter.GinLambda

// buildApp composes the full dependency graph and returns a ready gin.Engine.
func buildApp() (*gin.Engine, error) {
	// 1) DB
	db, err := utilities.NewDatabaseConnection()
	if err != nil {
		return nil, err
	}

	// 2) Repositories
	userRepo := db.CreateUserRepository()
	companyRepo := db.CreateCompanyRepository()
	towRepo := db.CreateTowRepository()

	// 3) Services
	userSvc := service.NewUserServiceWithMongo(userRepo)
	companySvc := service.NewCompanyService(companyRepo)
	towSvc := service.NewTowService(towRepo)

	// 4) Handlers
	userHandler := handler.NewUserHandler(userSvc)
	companyHandler := handler.NewCompanyHandler(companySvc)
	towHandler := handler.NewTowHandler(towSvc)

	// 5) Router
	router := utilities.NewRouter(userHandler, companyHandler, towHandler)
	engine := router.InitializeRouter()
	return engine, nil
}

func main() {

	// Initialize structured logging library
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Build once for both local and lambda
	engine, err := buildApp()

	if err != nil {
		log.Fatalf("failed to build application: %v", err)
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "local" // default to local if not set
	}

	if env == "local" {
		// LOCAL HTTP SERVER
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("[local] starting Gin on :%s", port)
		if err := engine.Run(":" + port); err != nil {
			log.Fatalf("failed to start local server: %v", err)
		}
		return
	}

	// LAMBDA RUNTIME
	ginLambda = ginadapter.New(engine)
	log.Println("[lambda] handler is ready")
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
