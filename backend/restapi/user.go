package restapi

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/utils/rest"
)

func UserEndpoints(us core.UserStorage) []*httpx.Endpoint {
	return []*httpx.Endpoint{
		createUserEp(us),
		updateUserEp(us),
		getUser(us),
		getUsers(us),
		deleteUser(us),
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
		Method:      echo.POST,
		Path:        "/user",
		Category:    "idx.user",
		Desc:        "Create a user",
		Version:     "v1",
		Permissions: []string{core.PermCreateUser},
		Handler:     handler,
	}
}

func updateUserEp(us core.UserStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var user core.User
		if err := etx.Bind(&user); err != nil {
			return errx.BadReq("failed to read user info from request", err)
		}

		if err := us.Update(etx.Request().Context(), &user); err != nil {
			return err
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.PUT,
		Path:        "/user",
		Category:    "idx.user",
		Desc:        "Update a user",
		Version:     "v1",
		Permissions: []string{core.PermCreateUser},
		Handler:     handler,
	}
}

func getUser(us core.UserStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		user, err := us.GetOne(etx.Request().Context(), id)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, user)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/user/:id",
		Category:    "idx.user",
		Desc:        "Get info for a user",
		Version:     "v1",
		Permissions: []string{core.PermGetUser},
		Handler:     handler,
	}
}

func getUsers(us core.UserStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		cmnParams, err := rest.GetCommonParams(etx)
		if err != nil {
			return err
		}

		users, err := us.Get(etx.Request().Context(), cmnParams)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, users)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/user",
		Category:    "idx.user",
		Desc:        "Get user list",
		Version:     "v1",
		Permissions: []string{core.PermGetUser},
		Handler:     handler,
	}
}

func deleteUser(us core.UserStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		err := us.Remove(etx.Request().Context(), id)
		if err != nil {
			return err
		}

		return etx.String(http.StatusOK, strconv.Itoa(id))
	}

	return &httpx.Endpoint{
		Method:      echo.DELETE,
		Path:        "/user/:id",
		Category:    "idx.user",
		Desc:        "Delete a user",
		Version:     "v1",
		Permissions: []string{core.PermGetUser},
		Handler:     handler,
	}
}
