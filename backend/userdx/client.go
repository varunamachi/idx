package userdx

import (
	"context"
	"net/url"
	"time"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

type Client struct {
	*httpx.Client
	Timeout time.Duration
}

func (c *Client) CurrentUser() *core.User {
	return c.User().(*core.User)
}

func (c *Client) build() *httpx.RequestBuilder {
	builder := c.Build()
	if c.Timeout != 0 {
		builder = builder.WithTimeout(c.Timeout)
	}
	return builder
}

func (c *Client) Register(
	gtx context.Context, user *core.User, password string) (int64, error) {
	up := core.UserWithPassword{User: user, Password: password}

	res := map[string]int64{"userId": int64(-1)}
	apiRes := c.build().Path("/api/v1/user").Post(gtx, up)
	if err := apiRes.LoadClose(&res); err != nil {
		return -1, errx.Errf(err, "failed to register user")
	}
	return res["userId"], nil

}

func (c *Client) Verify(
	gtx context.Context, userId, verifyId string) error {
	apiRes := c.build().
		Path("/api/v1/user/verify", userId, verifyId).
		Post(gtx, nil)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to verify user '%s'", userId)
	}
	return nil
}

func (c *Client) VerifyWithUrl(
	gtx context.Context, fullUrl string) error {

	u, err := url.Parse(fullUrl)
	if err != nil {
		return errx.Errf(err, "failed to parse verify url: %s", fullUrl)
	}

	apiRes := c.build().
		Path(u.Path).
		Post(gtx, nil)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to verify user with path '%s'", u.Path)
	}
	return nil
}

func (c *Client) UpdatePassword(
	gtx context.Context, userId, oldPwd, newPwd string) error {
	data := map[string]string{
		"username":    userId,
		"oldPassword": oldPwd,
		"newPassword": newPwd,
	}
	apiRes := c.build().Path("/api/v1/user/password").Put(gtx, data)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to update password for user '%s'", userId)
	}
	return nil
}

func (c *Client) InitResetPassword(
	gtx context.Context, userId string) error {
	apiRes := c.build().
		Path("/api/v1/user", userId, "password/reset/init").
		Post(gtx, nil)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err,
			"failed to initialize password reset for user '%s'", userId)
	}
	// TODO - think about handling mail
	return nil
}

func (c *Client) ResetPassword(
	gtx context.Context, userId, token, newPwd string) error {
	data := map[string]string{
		"userId":      userId,
		"token":       token,
		"newPassword": newPwd,
	}
	apiRes := c.build().Path("/api/v1/user/password/reset").Post(gtx, data)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to reset password of user '%s'", userId)
	}
	return nil
}

func (c *Client) Approve(
	gtx context.Context,
	userId int64,
	role auth.Role,
	groupIds ...int64) error {
	apiRes := c.build().
		Path("/api/v1/user", userId, "approve", string(role)).
		Patch(gtx, groupIds)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to set approve user '%d'", userId)
	}
	return nil
}

func (c *Client) Login(
	gtx context.Context, userId, password string) (*core.User, error) {

	creds := core.Creds{
		UniqueName: userId,
		Password:   password,
		Type:       core.AuthUser,
	}
	apiRes := c.build().Path("/api/v1/authenticate").Post(gtx, creds)

	authResult := struct {
		User  *core.User `json:"user"`
		Token string     `json:"token"`
	}{}

	if err := apiRes.LoadClose(&authResult); err != nil {
		return nil, errx.Errf(err, "failed to authenticate user '%d'", userId)
	}
	c.SetUser(authResult.User).SetToken(authResult.Token)

	return authResult.User, nil
}

func (c *Client) UpdateUser(gtx context.Context, user *core.User) error {
	apiRes := c.build().Path("/api/v1/user").Put(gtx, user)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to update user: '%s'", user.Username())
	}
	return nil
}

func (c *Client) GetUser(
	gtx context.Context, id int64) (*core.User, error) {
	var user core.User
	apiRes := c.build().Path("/api/v1/user", id).Get(gtx)
	if err := apiRes.LoadClose(&user); err != nil {
		return nil, errx.Errf(err, "failed to get user '%d'", id)
	}
	return &user, nil
}

func (c *Client) GetByUserId(
	gtx context.Context, id string) (*core.User, error) {
	var user core.User
	apiRes := c.build().Path("/api/v1/user/strid/", id).Get(gtx)
	if err := apiRes.LoadClose(&user); err != nil {
		return nil, errx.Errf(err, "failed to get user '%d'", id)
	}
	return &user, nil
}

func (c *Client) SetUserState(
	gtx context.Context, id int64, state core.UserState) error {
	apiRes := c.build().Path("/api/v1/user", id, "state", state).Put(gtx, nil)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err,
			"failed to set state '%s' for user '%d'", id, state)
	}
	return nil
}

func (c *Client) RemoveUser(gtx context.Context, id int64) error {
	apiRes := c.build().Path("/api/v1/user", id).Delete(gtx)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to delete user '%d'", id)
	}
	return nil
}

func (c *Client) GetUsers(
	gtx context.Context, params *data.CommonParams) ([]*core.User, error) {
	apiRes := c.build().Path("/api/v1/user").CmnParam(params).Get(gtx)
	users := make([]*core.User, 0, params.PageSize)
	if err := apiRes.LoadClose(&users); err != nil {
		return nil, errx.Errf(err, "failed to get user list")
	}
	return users, nil
}

func (c *Client) UserExists(gtx context.Context, id string) (bool, error) {
	apiRes := c.build().Path("/api/v1/user", id, "exists").Get(gtx)
	res := map[string]bool{"exists": false}
	if err := apiRes.LoadClose(&res); err != nil {
		return false, errx.Errf(err, "failed to check if user exists: '%s'", id)
	}
	return res["exists"], nil
}

func (c *Client) UserCount(
	gtx context.Context, filter *data.Filter) (int64, error) {
	apiRes := c.build().Path("/api/v1/user").Filter(filter).Get(gtx)
	res := map[string]int64{"count": 0}
	if err := apiRes.LoadClose(&res); err != nil {
		return 0, errx.Errf(err, "failed to get user count for filter")
	}
	return res["count"], nil
}
