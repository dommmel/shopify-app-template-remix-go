package user

import (
	"database/sql"
	"errors"
	"os"

	"github.com/dommmel/shopify-app-template-remix-go/pkg/shopify"
)

// UserRepository handles all database interactions related to the User model.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID retrieves a user by its ID.
func (r *UserRepository) GetUserByID(id uint64) (*User, error) {
	user := &User{}
	query := "SELECT id, accesstoken, myshopifydomain, scopes FROM users WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.AccessToken, &user.MyshopifyDomain, &user.Scopes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by its ID.
func (r *UserRepository) GetUserByMyshopify(myshopifyDomain string) (*User, error) {
	user := &User{}
	query := "SELECT id, accesstoken, myshopifydomain, scopes FROM users WHERE myshopifydomain = ?"
	err := r.db.QueryRow(query, myshopifyDomain).Scan(&user.ID, &user.AccessToken, &user.MyshopifyDomain, &user.Scopes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user in the database and returns the created User struct.
func (r *UserRepository) CreateUser(user *User) (*User, error) {
	query := "INSERT INTO users (accesstoken, myshopifydomain, scopes) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, user.AccessToken, user.MyshopifyDomain, user.Scopes)
	if err != nil {
		return nil, err
	}

	// Retrieve the last inserted ID and set it in the user struct
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = id

	// Return the created user
	return user, nil
}
func (r *UserRepository) UpdateUser(user *User) (*User, error) {
	query := "UPDATE users SET accesstoken = ?, scopes = ? WHERE myshopifydomain = ?"
	_, err := r.db.Exec(query, user.AccessToken, user.Scopes, user.MyshopifyDomain)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *UserRepository) FindOrCreateUserByEncodedSessionToken(sessionToken string) (*User, error) {
	decoded, err := shopify.DecodeSessionToken(sessionToken)
	if err != nil {
		return nil, err
	}
	myshopifyDomain := decoded.MyshopifyDomain
	user, err := r.GetUserByMyshopify(myshopifyDomain)

	// No user Found, create a new one
	if err != nil {
		// Call Shopify API to exchange the session token for an access token
		resp, err := shopify.GetAccessTokenFromShopify(myshopifyDomain, sessionToken)
		if err != nil {
			return nil, err
		}
		newUser := &User{
			AccessToken:     resp.AccessToken,
			MyshopifyDomain: myshopifyDomain,
			Scopes:          resp.Scope,
		}
		// Create the user in the database
		return r.CreateUser(newUser)
	}

	// User found but scopes have changed
	if user.Scopes != os.Getenv("SCOPES") {
		// Call Shopify API to exchange the session token for an access token
		resp, err := shopify.GetAccessTokenFromShopify(myshopifyDomain, sessionToken)
		if err != nil {
			return nil, err
		}
		user.AccessToken = resp.AccessToken
		user.Scopes = resp.Scope

		return r.UpdateUser(user)
	}

	return user, nil
}
