package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	echo "github.com/labstack/echo/v4"
)

var mySigningKey = []byte("ct-secret-key")

func GenerateToken(username string, expireDuration time.Duration) (string, error) {
	// Create the Claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "ct-backend-course",
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return ss, nil

}

func ExtractUserNameFn(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) (err error) {
		auth := c.Get("auth").(*jwt.Token)
		claims := auth.Claims.(*jwt.RegisteredClaims)
		username := claims.Subject
		c.Set("username", username)

		next(c)

		return nil
	}
}

func AuthMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		ContextKey: "auth",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
		SigningKey: mySigningKey,
	}
	return echojwt.WithConfig(config)
}
