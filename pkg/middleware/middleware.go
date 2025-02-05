package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"

	"myproject/pkg/config"
	"myproject/pkg/model"

	"github.com/golang-jwt/jwt"
)

type MiddlewareJwt interface {
	GenerateAdminToken(username string) (string, error)
	AdminAuthMiddleware() fiber.Handler
}

type MiddlewareJWT struct {
	Config config.Config
}

func (s MiddlewareJWT) GenerateAdminToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := model.AdminClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.Config.VnJWTKey))
}

// verify Admin Token
func AdminAuthentication(tokenString string, jwtKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*model.AdminClaims); ok && token.Valid {
		return claims.Username, nil
	}
	return "", errors.New("invalid token")
}

// Admin Auth middleware
func (s MiddlewareJWT) AdminAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing the authorization header"})
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
		}

		username, err := AdminAuthentication(tokenString, s.Config.AdJWTKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals("username", username)
		return c.Next()
	}
}
