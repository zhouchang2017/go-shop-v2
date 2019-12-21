package auth

// 授权对象
type Authenticatable interface {
	// Get the name of the unique identifier for the user.
	GetAuthIdentifierName() string

	// Get the unique identifier for the user.
	GetAuthIdentifier() string

	// Get the password for the user.
	GetAuthPassword() string
}
