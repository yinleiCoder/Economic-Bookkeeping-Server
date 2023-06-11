package models

import "github.com/golang-jwt/jwt/v5"

type CustomJwtClaims struct {
	ID       uint
	NickName string
	RoleId   uint // role
	jwt.RegisteredClaims
}
