package utils

import (
	common "boiler-platecode/src/common/const"
	"boiler-platecode/src/common/const/exception"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SendRestResponse[T any](ctx *gin.Context, output common.ServiceOutput[T]) {
	if output.Exception != nil {
		ctx.JSON(output.Exception.HttpStatusCode, common.ApiResponse[any]{
			Code:    output.Exception.Code,
			Message: fallbackIfEmpty(output.Message, output.Exception.Message),
		})
		return
	}

	ctx.JSON(output.HttpStatusCode, common.ApiResponse[T]{
		Code:    "000000",
		Message: fallbackIfEmpty(output.Message, "SUCCESS"),
		Data:    output.OutputData,
	})
}

func fallbackIfEmpty(preferred string, fallback string) string {
	if preferred != "" {
		return preferred
	}
	return fallback
}

func HashPassword(password string) string {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic("Error While hashing password: ", err)
	}
	return string(hashedPassword)
}
func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateJwtToken(tokenType common.TokenType, userId uint, expTimeInMinutes int, secret string) (string, error) {

	// secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"userId": userId,
		"type":   tokenType,
		"exp":    time.Now().Add(time.Minute * time.Duration(expTimeInMinutes)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJwtToken(tokenStr string, secret string) (jwt.MapClaims, error) {
	// secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Panic("JWT SECRET is not set in .env file")
	}

	// Parse token
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Ensure signing method is expected
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, common.ErrInvalidSigning
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, common.ErrExpiredToken
		}
		return nil, common.ErrInvalidToken
	}

	// Extract and assert claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, common.ErrInvalidToken

	}

	return claims, nil
}

func ServiceError[T any](code exception.ErrorCode) common.ServiceOutput[T] {
	return common.ServiceOutput[T]{
		Exception: exception.GetException(code),
	}
}
