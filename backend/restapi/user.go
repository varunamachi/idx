package restapi

import (
	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

func UserEndpoints(us core.UserStorage) []*httpx.Endpoint {
	return []*httpx.Endpoint{
		createUserEp(us),
	}
}

func createUserEp(us core.UserStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var user core.User
		if err := etx.Bind(&user); err != nil {
			return errx.BadReq("failed to read user info from request", err)
		}

		if err := us.Save(etx.Request().Context(), &user); err != nil {
			return err
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:     echo.POST,
		Path:       "/user",
		Category:   "idx.user",
		Desc:       "Create a user",
		Version:    "v1",
		Permission: core.PermCreateUser,
		Handler:    handler,
	}
}

func updateUserEp(us core.UserStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		return nil
	}

	return &httpx.Endpoint{
		Method:     echo.POST,
		Path:       "/user",
		Category:   "idx.user",
		Desc:       "Create a user",
		Version:    "v1",
		Permission: core.PermCreateUser,
		Handler:    handler,
	}
}

func getUser(us core.UserStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		return nil
	}

	return &httpx.Endpoint{
		Method:     echo.POST,
		Path:       "/user",
		Category:   "idx.user",
		Desc:       "Create a user",
		Version:    "v1",
		Permission: core.PermCreateUser,
		Handler:    handler,
	}
}

func getUsers(us core.UserStorage) *httpx.Endpoint {
	handler := func(c echo.Context) error {
		return nil
	}

	return &httpx.Endpoint{
		Method:     echo.POST,
		Path:       "/user",
		Category:   "idx.user",
		Desc:       "Create a user",
		Version:    "v1",
		Permission: core.PermCreateUser,
		Handler:    handler,
	}
}
