package main

import (
	"context"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "bweng/docs"
	"bweng/api/proto/user"
	"bweng/internal/user/config"
	"bweng/internal/user/handler"
	"bweng/internal/user/model"
	"bweng/internal/user/repository"
	"bweng/internal/user/service"
)

// @title User Service API
// @version 1.0
// @description This is a user service API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// gRPC server implementation
type userGRPCServer struct {
	user.UnimplementedUserServiceServer
	userService *service.UserService
}

func (s *userGRPCServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	createReq := &model.CreateUserRequest{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	userResp, err := s.userService.CreateUser(createReq)
	if err != nil {
		return &user.UserResponse{Error: err.Error()}, nil
	}

	return &user.UserResponse{
		User: &user.User{
			Id:        uint64(userResp.ID),
			Username:  userResp.Username,
			Email:     userResp.Email,
			FirstName: userResp.FirstName,
			LastName:  userResp.LastName,
		},
	}, nil
}

func (s *userGRPCServer) GetUserByID(ctx context.Context, req *user.GetUserByIDRequest) (*user.UserResponse, error) {
	userResp, err := s.userService.GetUserByID(uint(req.Id))
	if err != nil {
		return &user.UserResponse{Error: err.Error()}, nil
	}

	return &user.UserResponse{
		User: &user.User{
			Id:        uint64(userResp.ID),
			Username:  userResp.Username,
			Email:     userResp.Email,
			FirstName: userResp.FirstName,
			LastName:  userResp.LastName,
		},
	}, nil
}

func (s *userGRPCServer) GetUserByEmail(ctx context.Context, req *user.GetUserByEmailRequest) (*user.UserResponse, error) {
	userResp, err := s.userService.GetUserByEmail(req.Email)
	if err != nil {
		return &user.UserResponse{Error: err.Error()}, nil
	}

	return &user.UserResponse{
		User: &user.User{
			Id:        uint64(userResp.ID),
			Username:  userResp.Username,
			Email:     userResp.Email,
			FirstName: userResp.FirstName,
			LastName:  userResp.LastName,
		},
	}, nil
}

func (s *userGRPCServer) GetUserByUsername(ctx context.Context, req *user.GetUserByUsernameRequest) (*user.UserResponse, error) {
	userResp, err := s.userService.GetUserByUsername(req.Username)
	if err != nil {
		return &user.UserResponse{Error: err.Error()}, nil
	}

	return &user.UserResponse{
		User: &user.User{
			Id:        uint64(userResp.ID),
			Username:  userResp.Username,
			Email:     userResp.Email,
			FirstName: userResp.FirstName,
			LastName:  userResp.LastName,
		},
	}, nil
}

func (s *userGRPCServer) GetAllUsers(ctx context.Context, req *user.GetAllUsersRequest) (*user.GetAllUsersResponse, error) {
	users, err := s.userService.GetAllUsers()
	if err != nil {
		return &user.GetAllUsersResponse{}, err
	}

	var userList []*user.User
	for _, u := range users {
		userList = append(userList, &user.User{
			Id:        uint64(u.ID),
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		})
	}

	return &user.GetAllUsersResponse{Users: userList}, nil
}

func (s *userGRPCServer) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	updateReq := &model.UpdateUserRequest{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	userResp, err := s.userService.UpdateUser(uint(req.Id), updateReq)
	if err != nil {
		return &user.UserResponse{Error: err.Error()}, nil
	}

	return &user.UserResponse{
		User: &user.User{
			Id:        uint64(userResp.ID),
			Username:  userResp.Username,
			Email:     userResp.Email,
			FirstName: userResp.FirstName,
			LastName:  userResp.LastName,
		},
	}, nil
}

func (s *userGRPCServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	err := s.userService.DeleteUser(uint(req.Id))
	if err != nil {
		return &user.DeleteUserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &user.DeleteUserResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

func main() {
	// Initialize database configuration
	dbConfig := config.NewDatabaseConfig()

	// Connect to database
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Run database migrations
	if err := userRepo.Migrate(); err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handler
	userHandler := handler.NewUserHandler(userService)

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen for gRPC: %v", err)
		}

		s := grpc.NewServer()
		user.RegisterUserServiceServer(s, &userGRPCServer{userService: userService})
		reflection.Register(s)

		log.Println("gRPC server starting on port 50051...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

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
		// User routes
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.GET("/email", userHandler.GetUserByEmail)
			users.GET("/username", userHandler.GetUserByUsername)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "user-service",
		})
	})

	// Start HTTP server
	log.Println("User service starting on port 8080...")
	log.Println("Swagger documentation available at: http://localhost:8080/swagger/index.html")
	log.Println("gRPC server available at: localhost:50051")
	
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 