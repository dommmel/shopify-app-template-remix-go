package user

// A Shop represents a shop record in the database
type User struct {
	ID              int64
	AccessToken     string
	MyshopifyDomain string
	Scopes          string
}
