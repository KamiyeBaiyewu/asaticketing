package model

// Tokens is a wrapper for access and refresh tokens
type Tokens struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresAt  int64  `json:"expiresAt,omitempty"` //we return only the access Token's expires at time
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresAt int64  `json:"refreshTokenExpiration,omitempty"` //to be stored in the database with refresh tokens
}
