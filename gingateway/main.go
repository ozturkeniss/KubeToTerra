package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

// ServiceConfig represents a microservice configuration
type ServiceConfig struct {
	Name     string
	BasePath string
	Target   string
}

// Gateway represents the API Gateway
type Gateway struct {
	services map[string]*ServiceConfig
}

// NewGateway creates a new API Gateway
func NewGateway() *Gateway {
	return &Gateway{
		services: make(map[string]*ServiceConfig),
	}
}

// RegisterService registers a microservice with the gateway
func (g *Gateway) RegisterService(config *ServiceConfig) {
	g.services[config.Name] = config
	log.Printf("Registered service: %s -> %s", config.BasePath, config.Target)
}

// ProxyHandler handles the proxy routing
func (g *Gateway) ProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// Find the service for this path
		var targetService *ServiceConfig
		for _, service := range g.services {
			if strings.HasPrefix(path, service.BasePath) {
				targetService = service
				break
			}
		}

		if targetService == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Service not found",
				"path":  path,
			})
			return
		}

		// Create reverse proxy
		targetURL, err := url.Parse(targetService.Target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid target URL",
			})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		
		// Add custom director to modify the request
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = targetURL.Host
			
			// Log the request
			log.Printf("Proxying request: %s %s -> %s", req.Method, req.URL.Path, targetService.Target)
		}

		// Add error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Service unavailable"))
		}

		// Serve the request
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// HealthCheckHandler handles health check requests
func (g *Gateway) HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "api-gateway",
			"timestamp": time.Now().Unix(),
			"services": len(g.services),
		})
	}
}

// ServicesHandler returns information about registered services
func (g *Gateway) ServicesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		services := make([]gin.H, 0)
		for name, service := range g.services {
			services = append(services, gin.H{
				"name":      name,
				"base_path": service.BasePath,
				"target":    service.Target,
			})
		}
		
		c.JSON(http.StatusOK, gin.H{
			"services": services,
		})
	}
}

func main() {
	// Create gateway
	gateway := NewGateway()

	// Register microservices
	gateway.RegisterService(&ServiceConfig{
		Name:     "user-service",
		BasePath: "/api/v1/users",
		Target:   "http://localhost:8080",
	})

	gateway.RegisterService(&ServiceConfig{
		Name:     "order-service",
		BasePath: "/api/v1/orders",
		Target:   "http://localhost:8081",
	})

	// Setup Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add request logging middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[GATEWAY] %v | %3d | %13v | %15s | %-7s %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	})

	// Health check endpoint
	r.GET("/health", gateway.HealthCheckHandler())

	// Services info endpoint
	r.GET("/services", gateway.ServicesHandler())

	// API routes - proxy to microservices
	api := r.Group("/api/v1")
	{
		// All requests to /api/v1/* will be proxied
		api.Any("/*path", gateway.ProxyHandler())
	}

	// Start server
	log.Println("API Gateway starting on port 8082...")
	log.Println("Health check: http://localhost:8082/health")
	log.Println("Services info: http://localhost:8082/services")
	
	if err := r.Run(":8082"); err != nil {
		log.Fatal("Failed to start gateway:", err)
	}
} 