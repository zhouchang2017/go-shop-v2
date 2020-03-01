package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	err2 "go-shop-v2/pkg/err"
	"log"
	"net/http"
	"time"
)

type JWTGuard struct {
	name      string
	secretKey string
	exp       int64        // 过期时间,单位分钟
	provider  UserProvider // The user provider implementation.
	jwt       *JWT
	ctx       *gin.Context
	user      Authenticatable // The currently authenticated user.
}

func (this *JWTGuard) SetContext(ctx *gin.Context) {
	this.ctx = ctx
}

func (this *JWTGuard) GetContext() *gin.Context {
	return this.ctx
}

func (this *JWTGuard) Check() bool {
	if this.user != nil {
		return true
	}
	return false
}

func (this *JWTGuard) authenticate() (user Authenticatable, err error) {
	if this.user != nil {
		return this.user, nil
	}
	return nil, err2.NewFromCode(401)
}

func (this *JWTGuard) setUser(user Authenticatable) {
	this.user = user
}

func (this *JWTGuard) User() (user Authenticatable, err error) {
	user, err = this.authenticate()
	if err == nil {
		return user, nil
	}
	token, err := this.jwt.SetContext(this.ctx).GetToken()
	if err != nil && !this.jwt.Check() {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if validationError.Errors == jwt.ValidationErrorExpired {
				// token过期
				return nil, err2.New(http.StatusUnauthorized, validationError.Error())
			}
		}
		return nil, err2.New(http.StatusUnauthorized, err.Error())
	}
	payload := token.Claims.(jwt.MapClaims)

	authenticatable, err := this.provider.RetrieveById(payload["sub"])
	if err != nil {
		return nil, err2.New(http.StatusUnauthorized, err.Error())
	}
	this.setUser(authenticatable)
	return this.user, nil
}

func (this *JWTGuard) Id() (id string, err error) {
	user, err := this.authenticate()
	if err != nil {
		return "", err
	}
	return user.GetAuthIdentifier(), nil
}

func (this *JWTGuard) Validate(credentials map[string]string) bool {
	_, ok := this.Attempt(credentials, false)
	return ok
}

func (this *JWTGuard) Attempt(credentials map[string]string, login bool) (res interface{}, ok bool) {
	authenticatable, err := this.provider.RetrieveByCredentials(credentials)
	if err != nil {
		log.Printf("jwt guard attempt err:%s\n", err)
		return nil, false
	}
	if this.hasValidCredentials(authenticatable, credentials) {
		if login {
			data, err := this.Login(authenticatable)
			if err != nil {
				return nil, false
			}
			return data, true
		}
		return nil, true
	}
	return nil, false
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (l loginResponse) String() string {
	return l.TokenType + " " + l.AccessToken
}

// Make a token for a user.
func (this *JWTGuard) Login(user Authenticatable) (data interface{}, err error) {
	token, err := this.jwt.FromUser(user)
	if err != nil {
		return nil, err
	}
	return loginResponse{AccessToken: token, TokenType: "Bearer"}, nil
}

func (this *JWTGuard) GetProvider() UserProvider {
	return this.provider
}

func (this *JWTGuard) SetProvider(provider UserProvider) {
	this.provider = provider
}

func (this JWTGuard) Name() string {
	return this.name
}

func (this *JWTGuard) hasValidCredentials(user Authenticatable, credentials map[string]string) bool {
	if user == nil {
		return false
	}
	return this.provider.ValidateCredentials(user, credentials)
}

func NewJwtGuard(name string, secretKey string, exp int64, provider UserProvider) *JWTGuard {
	if exp == 0 {
		exp = 2
	}
	return &JWTGuard{
		name:      name,
		secretKey: secretKey,
		exp:       exp,
		provider:  provider,
		jwt:       &JWT{secretKey: secretKey, exp: time.Minute},
	}
}

func (this *JWTGuard) Logout(token string) error {
	panic("implement me")
}

// Token
type JWTToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
