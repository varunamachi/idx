package restapi

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

func AuthEndpoints(gtx context.Context) []*httpx.Endpoint {
	athr := core.Authenticator(gtx)
	return []*httpx.Endpoint{
		authenticateEp(athr),
		logout(athr),
	}
}

func authenticateEp(athr auth.UserAuthenticator) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		gtx := etx.Request().Context()

		var creds auth.AuthData
		if err := etx.Bind(&creds); err != nil {
			return errx.BadReqX(err, "failed to decode credentials")
		}

		if err := athr.Authenticate(gtx, creds); err != nil {
			return errx.Errf(err, "failed to authenticate user")
		}

		user, err := athr.GetUser(gtx, creds)
		if err != nil {
			return errx.Errf(err, "failed to retrieve user")
		}

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["userId"] = user.Id()

		// TODO - get from application configuration
		claims["exp"] = time.Now().Add(auth.UserSessionTimeout).Unix()
		claims["type"] = "user"

		signed, err := token.SignedString(auth.GetJWTKey())
		if err != nil {
			return errx.Errf(err, "failed to generate session token")
		}

		return httpx.SendJSON(etx, data.M{
			"user":  user,
			"token": signed,
		})

		// return user, signed, nil
	}
	return &httpx.Endpoint{
		Method:   echo.POST,
		Path:     "/authenticate",
		Category: "idx.auth",
		Desc:     "Authenticate an entity",
		Version:  "v1",
		Handler:  handler,
	}
}

func logout(athr auth.Authenticator) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		// var creds auth.AuthData
		// if err := etx.Bind(&creds); err != nil {
		// 	return errx.BadReqX(err, "failed to decode credentials")
		// }
		// err := athr.(etx.Request().Context(), creds)
		// if err != nil {
		// 	return err
		// }
		return nil
	}
	return &httpx.Endpoint{
		Method:   echo.POST,
		Path:     "/logout",
		Category: "idx.auth",
		Desc:     "Logout an entity",
		Version:  "v1",
		Role:     auth.Normal,
		Handler:  handler,
	}
}
