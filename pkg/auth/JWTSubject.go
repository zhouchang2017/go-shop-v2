package auth

type JWTSubject interface {
	// Get the identifier that will be stored in the subject claim of the JWT.
	GetJWTIdentifier() string

	//  Return a interface, containing any custom claims to be added to the JWT.
	GetJWTCustomClaims() interface{}
}
