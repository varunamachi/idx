package client

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

type IdxClient struct {
	*httpx.Client
}

func New(address string) *IdxClient {
	return &IdxClient{
		Client: httpx.NewClient(address, ""),
	}
}

func (c *IdxClient) CurrentUser() *core.User {
	return c.User().(*core.User)
}

func (c *IdxClient) Register(
	gtx context.Context, user *core.User, password string) (int64, error) {
	up := core.UserWithPassword{User: user, Password: password}

	res := map[string]int64{"userId": int64(-1)}
	apiRes := c.Build().Path("/api/v1/user").Post(gtx, up)
	if err := apiRes.LoadClose(&res); err != nil {
		return -1, errx.Errf(err, "failed to register user")
	}
	return res["userId"], nil

}

func (c *IdxClient) Verify(
	gtx context.Context, userId, verifyId string) error {
	apiRes := c.Build().
		Path("/api/v1/user/verify", userId, verifyId).
		Post(gtx, nil)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to verify user '%s'", userId)
	}
	return nil
}

func (c *IdxClient) UpdatePassword(
	gtx context.Context, userId, oldPwd, newPwd string) error {
	data := map[string]string{
		"userId":      userId,
		"oldPassword": oldPwd,
		"newPassword": newPwd,
	}
	apiRes := c.Build().Path("/user/password").Put(gtx, data)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to update password for user '%s'", userId)
	}
	return nil
}

func (c *IdxClient) InitResetPassword(
	gtx context.Context, userId string) error {
	apiRes := c.Build().
		Path("/api/v1/user", userId, "password/reset/init").
		Post(gtx, nil)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err,
			"failed to initialize password reset for user '%s'", userId)
	}
	// TODO - think about handling mail
	return nil
}

func (c *IdxClient) ResetPassword(
	gtx context.Context, userId, token, newPwd string) error {
	data := map[string]string{
		"userId":      userId,
		"token":       token,
		"newPassword": newPwd,
	}
	apiRes := c.Build().Path("/api/v1/user/password/reset").Post(gtx, data)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to reset password of user '%s'", userId)
	}
	return nil
}

func (c *IdxClient) Approve(
	gtx context.Context,
	userId string,
	role auth.Role,
	groupIds ...int64) error {
	apiRes := c.Build().
		Path("/api/v1/user", userId, "approve", string(role)).
		Post(gtx, groupIds)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to reset password of user '%s'", userId)
	}
	return nil
}

func (c *IdxClient) Login(
	gtx context.Context, userId, password string) (*core.User, error) {

	creds := auth.AuthData{}
	apiRes := c.Build().Path("/api/v1/authenticate").Post(gtx, creds)

	authResult := struct {
		User  *core.User `json:"user"`
		Token string     `json:"token"`
	}{}

	if err := apiRes.LoadClose(&authResult); err != nil {
		return nil, errx.Errf(err, "failed to authenticate user '%s'", userId)
	}
	c.SetUser(authResult.User).SetToken(authResult.Token)

	return authResult.User, nil
}

func (c *IdxClient) UpdateUser(gtx context.Context, user *core.User) error {
	apiRes := c.Build().Path("/user").Put(gtx, user)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to update user: '%s'", user.UserId)
	}
	return nil
}

func (c *IdxClient) GetUser(
	gtx context.Context, id int64) (*core.User, error) {
	var user core.User
	apiRes := c.Build().Path("/user", id).Get(gtx)
	if err := apiRes.LoadClose(&user); err != nil {
		return nil, errx.Errf(err, "failed to get user '%d'", id)
	}
	return &user, nil
}

func (c *IdxClient) GetByUserId(
	gtx context.Context, id string) (*core.User, error) {
	var user core.User
	apiRes := c.Build().Path("/user/strid/", id).Get(gtx)
	if err := apiRes.LoadClose(&user); err != nil {
		return nil, errx.Errf(err, "failed to get user '%d'", id)
	}
	return &user, nil
}

func (c *IdxClient) SetUserState(
	gtx context.Context, id int64, state core.UserState) error {
	apiRes := c.Build().Path("/user", id, "state", state).Put(gtx, nil)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err,
			"failed to set state '%s' for user '%d'", id, state)
	}
	return nil
}

func (c *IdxClient) RemoveUser(gtx context.Context, id int64) error {
	apiRes := c.Build().Path("/user", id).Delete(gtx)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to delete user '%d'", id)
	}
	return nil
}

func (c *IdxClient) GetUsers(
	gtx context.Context, params *data.CommonParams) ([]*core.User, error) {
	apiRes := c.Build().Path("/user").CmnParam(params).Get(gtx)
	users := make([]*core.User, 0, params.PageSize)
	if err := apiRes.LoadClose(&users); err != nil {
		return nil, errx.Errf(err, "failed to get user list")
	}
	return users, nil
}

func (c *IdxClient) UserExists(gtx context.Context, id string) (bool, error) {
	apiRes := c.Build().Path("/user", id, "exists").Get(gtx)
	res := map[string]bool{"exists": false}
	if err := apiRes.LoadClose(&res); err != nil {
		return false, errx.Errf(err, "failed to check if user exists: '%s'", id)
	}
	return res["exists"], nil
}

func (c *IdxClient) UserCount(
	gtx context.Context, filter *data.Filter) (int64, error) {
	apiRes := c.Build().Path("/user").Filter(filter).Get(gtx)
	res := map[string]int64{"count": 0}
	if err := apiRes.LoadClose(&res); err != nil {
		return 0, errx.Errf(err, "failed to get user count for filter")
	}
	return res["count"], nil
}
