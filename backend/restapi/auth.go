package restapi

import (
	"github.com/labstack/echo/v4"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

func AuthEndpoints(athr auth.Authenticator) []*httpx.Endpoint {
	return []*httpx.Endpoint{
		authenticateEp(athr),
		logout(athr),
	}
}

func authenticateEp(athr auth.Authenticator) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var creds auth.AuthData
		if err := etx.Bind(&creds); err != nil {
			return errx.BadReqX(err, "failed to decode credentials")
		}
		err := athr.Authenticate(etx.Request().Context(), creds)
		if err != nil {
			return err
		}
		return nil
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
		var creds auth.AuthData
		if err := etx.Bind(&creds); err != nil {
			return errx.BadReqX(err, "failed to decode credentials")
		}
		err := athr.Authenticate(etx.Request().Context(), creds)
		if err != nil {
			return err
		}
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
