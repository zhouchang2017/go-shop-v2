package auth

type UserProvider interface {
	// 根据用户的唯一标识符检索用户.
	RetrieveById(identifier interface{}) (Authenticatable, error)

	// 根据给定的凭据检索用户.
	RetrieveByCredentials(credentials map[string]string) (Authenticatable, error)

	// 根据给定的凭据验证用户.
	ValidateCredentials(user Authenticatable, credentials map[string]string) bool
}
