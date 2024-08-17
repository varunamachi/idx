package restapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/utils/rest"
)

func UserEndpoints(gtx context.Context) []*httpx.Endpoint {

	us := core.UserCtlr(gtx)
	return []*httpx.Endpoint{
		registerUserEp(us),
		verifyUserEp(us),
		updateUserEp(us),
		getUserEp(us),
		getUserByUserIdEp(us),
		getUsersEp(us),
		deleteUserEp(us),
		updatePasswordEp(us),
		initPasswordResetEp(us),
		resetPasswordEp(us),
		approveEp(us),
		setStateEp(us),
		userExistsEp(us),
		userCountEp(us),
	}
}

func registerUserEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var up core.UserWithPassword
		if err := etx.Bind(&up); err != nil {
			return errx.BadReqX(err, "failed to read user info from request")
		}

		userId, err := us.Register(
			etx.Request().Context(), up.User, up.Password)
		if err != nil {
			return errx.Wrap(err)
		}
		return httpx.SendJSON(etx, data.M{"userId": userId})
	}

	return &httpx.Endpoint{
		Method:      echo.POST,
		Path:        "/user",
		Category:    "idx.user",
		Desc:        "Create a user",
		Version:     "v1",
		Permissions: []string{},
		Handler:     handler,
	}
}

func verifyUserEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		userId := etx.Param("userId")
		verToken := etx.Param("toekn")
		err := us.Verify(etx.Request().Context(), userId, verToken)
		if err != nil {
			return errx.Wrap(err)
		}
		return httpx.SendJSON(etx, data.M{
			"userId": userId,
		})
	}

	return &httpx.Endpoint{
		Method:      echo.POST,
		Path:        "/user/verify/:userId/:token",
		Category:    "idx.user",
		Desc:        "Verify user account",
		Version:     "v1",
		Permissions: []string{},
		Handler:     handler,
	}
}

func updateUserEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var user core.User
		if err := etx.Bind(&user); err != nil {
			return errx.BadReqX(err, "failed to read user info from request")
		}

		if err := us.Update(etx.Request().Context(), &user); err != nil {
			return errx.Wrap(err)
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

func getUserEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int64("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		user, err := us.GetOne(etx.Request().Context(), id)
		if err != nil {
			return errx.Wrap(err)
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

func getUserByUserIdEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Str("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		user, err := us.GetByUserId(etx.Request().Context(), id)
		if err != nil {
			return errx.Wrap(err)
		}

		return httpx.SendJSON(etx, user)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/user/strid/:id",
		Category:    "idx.user",
		Desc:        "Get info for a user identified by string id",
		Version:     "v1",
		Permissions: []string{core.PermGetUser},
		Handler:     handler,
	}
}

func getUsersEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		cmnParams, err := rest.GetCommonParams(etx)
		if err != nil {
			return errx.Wrap(err)
		}

		users, err := us.Get(etx.Request().Context(), cmnParams)
		if err != nil {
			return errx.Wrap(err)
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

func deleteUserEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int64("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		err := us.Remove(etx.Request().Context(), id)
		if err != nil {
			return errx.Wrap(err)
		}

		return etx.String(http.StatusOK, strconv.FormatInt(id, 10))
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

func updatePasswordEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {

		credx := struct {
			UserId      string `json:"userId"`
			OldPassword string `json:"oldPassword"`
			NewPassword string `json:"newPassword"`
		}{}

		if err := etx.Bind(&credx); err != nil {
			return errx.BadReqX(err, "invalid cred info given")
		}

		err := us.UpdatePassword(
			etx.Request().Context(),
			credx.UserId,
			credx.OldPassword,
			credx.NewPassword)
		if err != nil {
			return errx.Wrap(err)
		}

		return etx.String(http.StatusOK, credx.UserId)
	}

	return &httpx.Endpoint{
		Method:      echo.PUT,
		Path:        "/user/password",
		Category:    "idx.user",
		Desc:        "Update user password",
		Version:     "v1",
		Permissions: []string{},
		Handler:     handler,
	}
}

func initPasswordResetEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {

		prmg := httpx.NewParamGetter(etx)
		user := prmg.Str("userId")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		err := us.InitResetPassword(etx.Request().Context(), user)
		if err != nil {
			return errx.Wrap(err)
		}

		return etx.String(http.StatusOK, user)
	}

	return &httpx.Endpoint{
		Method:   echo.POST,
		Path:     "/user/:user/password/reset/init",
		Category: "idx.user",
		Desc:     "Update user password",
		Version:  "v1",
		Handler:  handler,
	}
}

func resetPasswordEp(us core.UserController) *httpx.Endpoint {

	handler := func(etx echo.Context) error {

		credx := struct {
			UserId      string `json:"userId"`
			Token       string `json:"token"`
			NewPassword string `json:"newPassword"`
		}{}

		if err := etx.Bind(&credx); err != nil {
			return errx.BadReqX(err, "invalid cred info given")
		}

		err := us.ResetPassword(
			etx.Request().Context(),
			credx.UserId,
			credx.Token,
			credx.NewPassword)
		if err != nil {
			return errx.Wrap(err)
		}

		return etx.String(http.StatusOK, credx.UserId)
	}

	return &httpx.Endpoint{
		Method:   echo.PUT,
		Path:     "/user/password/reset",
		Category: "idx.user",
		Desc:     "Update user password",
		Version:  "v1",
		Handler:  handler,
	}
}

func approveEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {

		pmg := httpx.NewParamGetter(etx)
		user := pmg.Str("user")
		role := auth.Role(pmg.Str("role"))

		groupIds := make([]int64, 0, 100)
		if err := etx.Bind(&groupIds); err != nil {
			return errx.BadReqX(
				err, "failed to get groupIds for user '%s'", user)
		}

		if !data.OneOf(role, auth.ValidRoles...) {
			return errx.BadReq("invalid role '%s' provided", role)
		}

		err := us.Approve(etx.Request().Context(), user, role)
		if err != nil {
			return errx.Wrap(err)
		}

		return etx.String(http.StatusOK, string(role))
	}

	return &httpx.Endpoint{
		Method:      echo.PATCH,
		Path:        "/user/:user/approve/:role",
		Category:    "idx.user",
		Desc:        "Approve a user with appropriate role",
		Version:     "v1",
		Permissions: []string{core.PermSetUserState},
		Handler:     handler,
	}
}

func setStateEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {

		pmg := httpx.NewParamGetter(etx)
		user := pmg.Int64("user")
		state := core.UserState(pmg.Str("state"))

		if !data.OneOf(state, core.ValidUserStates...) {
			return errx.BadReq("invalid user state '%s' provided", state)
		}

		err := us.SetState(etx.Request().Context(), user, state)
		if err != nil {
			return errx.Wrap(err)
		}

		return etx.String(http.StatusOK, string(state))
	}

	return &httpx.Endpoint{
		Method:      echo.PUT,
		Path:        "/user/:user/state/:state",
		Category:    "idx.user",
		Desc:        "Approve a user with appropriate role",
		Version:     "v1",
		Permissions: []string{core.PermSetUserState},
		Handler:     handler,
	}
}

func userExistsEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Str("userId")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		res, err := us.Exists(etx.Request().Context(), id)
		if err != nil {
			return errx.Wrap(err)
		}

		return httpx.SendJSON(etx, map[string]bool{
			"exists": res,
		})
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/user/:userId/exists",
		Category:    "idx.user",
		Desc:        "Check if an user exists",
		Version:     "v1",
		Permissions: []string{core.PermGetUser},
		Handler:     handler,
	}
}

func userCountEp(us core.UserController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		filter, err := rest.GetFilter(etx)
		if err != nil {
			return errx.Wrap(err)
		}

		res, err := us.Count(etx.Request().Context(), filter)
		if err != nil {
			return errx.Wrap(err)
		}

		return httpx.SendJSON(etx, map[string]int64{
			"count": res,
		})
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/user/count",
		Category:    "idx.user",
		Desc:        "Get user list",
		Version:     "v1",
		Permissions: []string{core.PermGetUser},
		Handler:     handler,
	}
}
