package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"time"
)

type JWT struct {
	secretKey string
	exp       time.Duration
	ctx       *gin.Context
	token     *jwt.Token
}

// Set the gin context.
func (this *JWT) SetContext(ctx *gin.Context) *JWT {
	this.ctx = ctx
	return this
}

// Set the token.
func (this *JWT) SetToken(token interface{}) (*JWT, error) {
	switch token.(type) {
	case string:
		t, err := this.decode(token.(string))
		if err != nil {
			return this, err
		}
		this.token = t
	case *jwt.Token:
		this.token = token.(*jwt.Token)
	default:
		return this, fmt.Errorf("jwt SetToken token type error:%+v", token)
	}
	return this, nil
}

// Unset the current token.
func (this *JWT) UnsetToken() *JWT {
	this.token = nil
	return this
}

// Alias to generate a token for a given user.
func (this *JWT) FromUser(user Authenticatable) (string, error) {
	return this.FromSubject(user)
}

// Generate a token for a given subject.
func (this *JWT) FromSubject(user Authenticatable) (string, error) {
	payload := this.MakePayload(user, time.Now().Add(this.exp))
	return this.encode(payload)
}

// Refresh an expired token.
func (this *JWT) Refresh(user Authenticatable) (string, error) {
	if err := this.requireToken(); err != nil {
		return "", err
	}

	payload := this.MakePayload(user, time.Now().Add(time.Hour*24*7))
	return this.encode(payload)
}

// Make a Payload instance.
func (this *JWT) MakePayload(user Authenticatable, exp time.Time) jwt.Claims {
	claims := make(jwt.MapClaims)
	claims["exp"] = exp.Unix()
	claims["iat"] = time.Now().Unix()
	if jwtsubUser, ok := user.(JWTSubject); ok {
		claims["sub"] = jwtsubUser.GetJWTIdentifier()
		this.setCustomClaims(jwtsubUser, claims)
	} else {
		claims["sub"] = user.GetAuthIdentifier()
	}

	return claims
}

// Get the raw Payload instance.
func (this *JWT) GetPayload() (claims jwt.Claims, err error) {
	if err := this.requireToken(); err != nil {
		return nil, err
	}
	return this.token.Claims, nil
}

// Set custom attr in claims
func (this *JWT) setCustomClaims(user JWTSubject, claims map[string]interface{}) {
	claims["custom"] = user.GetJWTCustomClaims()
}

// Get the token.
func (this *JWT) GetToken() (*jwt.Token, error) {
	if this.token == nil {
		_, err := this.ParseToken()
		if err != nil {
			return nil, err
		}
	}
	return this.token, nil
}

func (this *JWT) Check() bool {
	if this.token != nil {
		if this.token.Claims.Valid() != nil {
			return false
		}
		return true
	}
	return false
}

// Iterate through the parsers and attempt to retrieve
func (this *JWT) ParseToken() (*JWT, error) {
	token, err := request.AuthorizationHeaderExtractor.ExtractToken(this.ctx.Request)
	if err != nil {
		return this, err
	}
	return this.SetToken(token)
}

func (this *JWT) encode(claims jwt.Claims) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	jwtToken.Claims = claims
	return jwtToken.SignedString([]byte(this.secretKey))
}

func (this *JWT) decode(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(this.secretKey), nil
	})
}

func (this *JWT) requireToken() error {
	if this.token == nil {
		return errors.New("A token is required")
	}
	return nil
}
