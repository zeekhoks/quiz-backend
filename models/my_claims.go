package models

import "github.com/dgrijalva/jwt-go"

type MyUserClaims struct {
	User User
	jwt.StandardClaims
}
