package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "bweng/docs"
	"bweng/internal/order/config"
	"bweng/internal/order/handler"
	"bweng/internal/order/repository"
	"bweng/internal/order/service"
)

// @title Order Service API
// @version 1.0
// @description This is an order service API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Initialize database configuration
	dbConfig := config.NewDatabaseConfig()

	// Connect to database
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repository
	orderRepo := repository.NewOrderRepository(db)

	// Run database migrations
	if err := orderRepo.Migrate(); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize user client
	userClient, err := service.NewUserClient("localhost:50051")
	if err != nil {
		log.Fatal("Failed to connect to user service:", err)
	}
	defer userClient.Close()
	log.Println("Connected to user service via gRPC")

	// Initialize service
	orderService := service.NewOrderService(orderRepo, userClient)

	// Initialize handler
	orderHandler := handler.NewOrderHandler(orderService)

	// Setup Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Order routes
		orders := api.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetAllOrders)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.GET("/user/:user_id", orderHandler.GetOrdersByUserID)
			orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			orders.DELETE("/:id", orderHandler.DeleteOrder)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "order-service",
		})
	})

	// Start server
	log.Println("Order service starting on port 8081...")
	log.Println("Swagger documentation available at: http://localhost:8081/swagger/index.html")
	
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 