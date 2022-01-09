package httpx

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/sause/libx/env"
	"github.com/varunamachi/sause/libx/errx"
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
			err = fmt.Errorf("invalid token found in context")
		}
	} else {
		header := ctx.Request().Header.Get("Authorization")
		authSchemeLen := len("Bearer")
		if len(header) > authSchemeLen {
			tokStr := header[authSchemeLen+1:]
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				return getJWTKey(), nil
			}
			token, err = jwt.Parse(tokStr, keyFunc)
		} else {
			err = fmt.Errorf("unexpected auth scheme used to JWT")
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

func getAuthMiddleware(server *Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(etx echo.Context) error {
			userId, err := retrieveUserId(etx)
			if err != nil {
				return &echo.HTTPError{
					Code:     http.StatusForbidden,
					Message:  "invalid JWT information",
					Internal: err,
				}
			}

			ep, ok := etx.Get("endpoint").(Endpoint)
			if !ok {
				return &echo.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: "could not find endpoint information",
				}
			}

			user, err := server.userRetriever(userId)
			if err != nil {
				return err
			}

			if !user.HasPerms(ep.Permission) || !user.HasRole(ep.Role) {

				return &echo.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "permission to access resource is denied",
				}
			}

			etx.Set("endpoint", ep)
			etx.Set("user", user)
			return next(etx)
		}
	}
}

func getAccessMiddleware(ep *Endpoint, server *Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(etx echo.Context) error {
			err := next(etx)
			if err == nil {
				if env.GetConfig().PrintAllAccess {
					status := etx.Response().Status
					log.Debug().
						Int("statusCode", status).
						Str("user", GetUserId(etx)).
						Str("method", ep.Route.Method).
						Str("path", ep.Route.Path).
						Str("perm", ep.Permission).
						Str("role", string(ep.Role)).
						Msg(http.StatusText(status))
				}
				return nil
			}

			httpErr, ok := err.(*echo.HTTPError)
			if !ok {
				httpErr = errx.InternalServerErr(err)
			}

			log.Error().
				Int("statusCode", httpErr.Code).
				Str("user", GetUserId(etx)).
				Str("method", ep.Route.Method).
				Str("path", ep.Route.Path).
				Str("perm", ep.Permission).
				Str("role", string(ep.Role)).
				Msg(StrMsg(httpErr))
			return err
		}
	}
}