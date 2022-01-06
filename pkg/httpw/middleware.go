package httpw

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

//getJWTKey - gives a unique JWT key
func getJWTKey() string {
	jwtKey := os.Getenv("SAUSE_JWT_KEY")
	if len(jwtKey) == 0 {
		jwtKey = uuid.NewString()
	}
	return jwtKey
}

//getToken - gets token from context or from header
func getToken(ctx echo.Context) (token *jwt.Token, err error) {
	itk := ctx.Get("token")
	if itk != nil {
		var ok bool
		if token, ok = itk.(*jwt.Token); !ok {
			err = fmt.Errorf("Invalid token found in context")
		}
	} else {
		header := ctx.Request().Header.Get("Authorization")
		authSchemeLen := len("Bearer")
		if len(header) > authSchemeLen {
			tokStr := header[authSchemeLen+1:]
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				return getJWTKey(), nil
			}
			token = new(jwt.Token)
			token, err = jwt.Parse(tokStr, keyFunc)
		} else {
			err = fmt.Errorf("Unexpected auth scheme used to JWT")
		}
	}
	return token, err
}

//RetrieveSessionInfo - retrieves session information from JWT token
func retrieveUserId(ctx echo.Context) (string, error) {
	token, err := getToken(ctx)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims in JWT")
	}

	userId, ok := claims["userID"].(string)
	if !ok {
		return "", fmt.Errorf("couldnt find userId in token")
	}

	return userId, nil
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(etx echo.Context) (err error) {
		userId, err := retrieveUserId(etx)
		if err != nil {
			err = &echo.HTTPError{
				Code:     http.StatusForbidden,
				Message:  "invalid JWT information",
				Internal: err,
			}
		}
		etx.Set("userID", userId)
		return next(etx)
	}
}
