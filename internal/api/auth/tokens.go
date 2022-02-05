package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
)

var jwtKey = []byte("sectre key")                           //TODO: change the key
var accessTokenDuration = time.Duration(60) * time.Minute   //30 Mins
var refreshTokenDuration = time.Duration(40*24) * time.Hour //30 days

// CustomClaims - wraps the jwt standard claims, so User info can be added.
type CustomClaims struct {
	UserID model.UserID `json:"userID,omitempty"`
	Name   string       `json:"name,omitempty"`
	Role   string       `json:"role,omitempty"`
	Type   string       `json:"type,omitempty"`
	jwt.StandardClaims
}

// Tokens is a wrapper for access and refresh tokens
type Tokens struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresAt  int64  `json:"expiresAt,omitempty"` //we return only the access Token's expires at time
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresAt int64  `json:"-"` //to be stored in the database with refresh tokens
}

//IssueToken generate access and refresh token
func IssueToken(principal model.Principal) (*Tokens, error) {

	if principal.UserID == model.NilUserID {
		return nil, errors.New("invalid principal")
	}

	accessToken, accessTokenExpiresAt, err := generateToken(principal, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenExpiresAt, err := generateToken(principal, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}
	return tokens, nil
}

func generateToken(principal model.Principal, duration time.Duration) (string, int64, error) {

	now := time.Now()

	// Generate Access Tokens
	claims := &CustomClaims{
		UserID: principal.UserID,
		Name:   principal.Name,
		Role:   principal.Role,
		Type:   principal.Type,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, claims.ExpiresAt, nil

}

// VerifyToken -  checks the the token submitted is valid
func VerifyToken(accessToken string) (*model.Principal, error) {

	claims := &CustomClaims{}

	tkn, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	principal := &model.Principal{
		UserID: claims.UserID,
		Name:   claims.Name,
		Role:   claims.Role,
		Type:   claims.Type,
	}

	// return principal even if token is invalid because we need to get the UserID
	if !tkn.Valid {
		return principal, err
	}

	return principal, nil
}
