package user

import (
	"os"

	"github.com/dommmel/shopify-app-template-remix-go/pkg/shopify"
)

type UserService struct {
	repo *UserSQLiteRepository
}

func NewUserService(repo *UserSQLiteRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindOrCreateUserByEncodedSessionToken(sessionToken string) (*User, error) {
	decoded, err := shopify.DecodeSessionToken(sessionToken)
	if err != nil {
		return nil, err
	}
	myshopifyDomain := decoded.MyshopifyDomain

	// Try to get user by Shopify domain
	u, err := s.repo.GetUserByMyshopify(myshopifyDomain)
	if err != nil {
		// No user found, call Shopify API to create one
		resp, err := shopify.GetAccessTokenFromShopify(myshopifyDomain, sessionToken)
		if err != nil {
			return nil, err
		}

		// Create a new user
		newUser := &User{
			AccessToken:     resp.AccessToken,
			MyshopifyDomain: myshopifyDomain,
			Scopes:          resp.Scope,
		}

		return s.repo.CreateUser(newUser)
	}

	// If user found but scopes have changed, update the user
	if u.Scopes != os.Getenv("SCOPES") {
		resp, err := shopify.GetAccessTokenFromShopify(myshopifyDomain, sessionToken)
		if err != nil {
			return nil, err
		}

		u.AccessToken = resp.AccessToken
		u.Scopes = resp.Scope

		return s.repo.UpdateUser(u)
	}

	return u, nil
}
