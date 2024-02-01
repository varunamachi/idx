package restapi

import (
	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/httpx"
)

func UserEndpoints(us core.UserStorage) []*httpx.Endpoint {
	return []*httpx.Endpoint{}
}

func createUserEp(us core.UserStorage) echo.HandlerFunc {
	return func(etx echo.Context) error {
		return nil
	}
}

func updateUserEp(us core.UserStorage) echo.HandlerFunc {
	return func(etx echo.Context) error {
		return nil
	}
}

func getUser(us core.UserStorage) echo.HandlerFunc {
	return func(etx echo.Context) error {
		return nil
	}
}

func getUsers(us core.UserStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
