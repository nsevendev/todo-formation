package auth

import "github.com/golang-jwt/jwt/v5"

type tokenClaims struct {
	IdUser string `json:"id_user"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}