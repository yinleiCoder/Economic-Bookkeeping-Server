package middlewares

import (
	"economic_api/user_web/global"
	"economic_api/user_web/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var (
	TokenExpired = errors.New("token is expired")
	TokenInvalid = errors.New("can't handle this token")
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.Request.Header.Get("x-token")
		if tokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, map[string]string{
				"message": "please login",
			})
			ctx.Abort()
			return
		}
		jwtAuth := NewJWT()
		claims, err := jwtAuth.ParseToken(tokenStr)
		if err != nil {
			if err == TokenExpired {
				ctx.JSON(http.StatusUnauthorized, map[string]string{
					"message": "token has expired, please login again",
				})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusUnauthorized, "login error:"+err.Error())
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		ctx.Set("userId", claims.ID)
		ctx.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey),
	}
}

func (j *JWT) CreateToken(claims models.CustomJwtClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenStr string) (*models.CustomJwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.CustomJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomJwtClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

func (j JWT) RefreshToken(tokenStr string) (string, error) {
	jwt.WithTimeFunc(func() time.Time {
		return time.Unix(0, 0)
	})
	token, err := jwt.ParseWithClaims(tokenStr, &models.CustomJwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomJwtClaims); ok && token.Valid {
		jwt.WithTimeFunc(func() time.Time {
			return time.Now()
		})
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
