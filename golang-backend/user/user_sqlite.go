package user

import (
	"database/sql"
	"errors"
)

// UserSQLiteRepository handles all database interactions related to the User model.
type UserSQLiteRepository struct {
	db *sql.DB
}

// NewUserSQLiteRepository creates a new instance of UserRepository.
func NewUserSQLiteRepository(db *sql.DB) *UserSQLiteRepository {
	return &UserSQLiteRepository{db: db}
}

// GetUserByID retrieves a user by its ID.
func (r *UserSQLiteRepository) GetUserByID(id uint64) (*User, error) {
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
func (r *UserSQLiteRepository) GetUserByMyshopify(myshopifyDomain string) (*User, error) {
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
func (r *UserSQLiteRepository) CreateUser(user *User) (*User, error) {
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
func (r *UserSQLiteRepository) UpdateUser(user *User) (*User, error) {
	query := "UPDATE users SET accesstoken = ?, scopes = ? WHERE myshopifydomain = ?"
	_, err := r.db.Exec(query, user.AccessToken, user.Scopes, user.MyshopifyDomain)
	if err != nil {
		return nil, err
	}

	return user, nil
}
