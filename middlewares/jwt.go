package middlewares

import (
	"errors"
	"gotrack/helpers/common"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	jwt.StandardClaims
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := GetJwtTokenFromHeader(c)
		if err != nil {
			common.GenerateErrorResponse(c, err.Error())
			return
		}

		data, ok := DummyRedis[tokenString]
		if !ok {
			common.GenerateErrorResponse(c, "token invalid, please log in again")
			return
		}

		if time.Now().After(data.ExpiredAt) {
			common.GenerateErrorResponse(c, "token expired, please log in again")
			return
		}

		c.Set("auth", data)

		c.Next()
	}
}

func GetJwtTokenFromHeader(c *gin.Context) (tokenString string, err error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		return tokenString, errors.New("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return tokenString, errors.New("invalid Authorization header format")
	}

	return parts[1], nil
}

func GenerateJwtToken() (token string, err error) {
	expirationTime := time.Now().Add(2 * time.Hour)

	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	generatedTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = generatedTokenJwt.SignedString([]byte(os.Getenv("jwt_secret_key")))
	if err != nil {
		return
	}

	return
}

func AuthorizeRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("auth")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		loginData, ok := user.(UserLoginRedis)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			c.Abort()
			return
		}

		// Log untuk debugging
		// fmt.Printf("User role: %s, Required role: %s\n", loginData.Role, role)

		if loginData.Role != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
