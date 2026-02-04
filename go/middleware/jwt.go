package middleware

import (
	"net/http"
	"solace/orm"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Auth struct {
	Claims jwt.MapClaims
	User   *orm.User
}

func JWTMiddleware(db *gorm.DB, secret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return next(c)
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected jwt signing method")
				}
				return secret, nil
			})
			if err != nil || !token.Valid {
				return next(c)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return next(c)
			}

			sub, _ := claims["sub"].(string)
			if sub == "" {
				return next(c)
			}

			var user orm.User
			if err := user.ByID(db, sub); err != nil {
				return next(c)
			}

			c.Set("auth", Auth{
				Claims: claims,
				User:   &user,
			})

			return next(c)
		}
	}
}

func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Get("auth")
			if auth == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
			}
			return next(c)
		}
	}
}
