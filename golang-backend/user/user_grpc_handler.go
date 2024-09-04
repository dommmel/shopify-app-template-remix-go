package user

import (
	"context"
	"log"

	pb "github.com/dommmel/shopify-app-template-remix-go/generated/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler implements the gRPC UserServiceServer interface.
type UserHandler struct {
	repo *UserRepository
	pb.UnimplementedUserServiceServer
}

// NewUserHandler creates a new instance of UserHandler with the provided repository.
func NewUserHandler(repo *UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetUser handles the GetUser gRPC request and uses the repository directly.
func (h *UserHandler) GetUser(ctx context.Context, userRequest *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := h.repo.GetUserByID(uint64(userRequest.ID))
	if err != nil {
		if err.Error() == "user not found" {
			log.Printf("No user found with ID %d", userRequest.ID)
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		log.Printf("Error fetching user with ID %d: %v", userRequest.ID, err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return &pb.UserResponse{
		ID:              int64(user.ID),
		AccessToken:     user.AccessToken,
		MyshopifyDomain: user.MyshopifyDomain,
		Scopes:          user.Scopes,
	}, nil
}

// FindOrCreateUserByEncodedSessionToken handles the token-based user creation.
func (h *UserHandler) FindOrCreateUserByEncodedSessionToken(ctx context.Context, tokenRequest *pb.TokenRequest) (*pb.UserResponse, error) {
	user, err := h.repo.FindOrCreateUserByEncodedSessionToken(tokenRequest.Token)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
	return &pb.UserResponse{
		ID:              user.ID,
		AccessToken:     user.AccessToken,
		MyshopifyDomain: user.MyshopifyDomain,
		Scopes:          user.Scopes,
	}, nil

}
