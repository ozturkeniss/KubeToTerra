package service

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"bweng/api/proto/user"
)

// UserClient handles communication with user service
type UserClient struct {
	client user.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserClient creates a new user client
func NewUserClient(userServiceAddr string) (*UserClient, error) {
	conn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := user.NewUserServiceClient(conn)

	return &UserClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *UserClient) Close() error {
	return c.conn.Close()
}

// GetUserByID retrieves a user by ID from user service
func (c *UserClient) GetUserByID(ctx context.Context, userID uint64) (*user.User, error) {
	req := &user.GetUserByIDRequest{
		Id: userID,
	}

	resp, err := c.client.GetUserByID(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, ErrUserNotFound
	}

	return resp.User, nil
}

// ValidateUserExists checks if a user exists
func (c *UserClient) ValidateUserExists(ctx context.Context, userID uint64) error {
	_, err := c.GetUserByID(ctx, userID)
	return err
} 