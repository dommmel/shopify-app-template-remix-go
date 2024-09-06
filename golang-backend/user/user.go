package user

// A Shop represents a shop record in the database
type User struct {
	ID              int64
	AccessToken     string
	MyshopifyDomain string
	Scopes          string
}
type UserRepository interface {
	GetUserByID(id uint64) (*User, error)
	GetUserByMyshopify(myshopifyDomain string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(user *User) (*User, error)
}
