package client

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

type IdxClient struct {
	*httpx.Client
	user *core.User
}

func New(address string) *IdxClient {
	return &IdxClient{
		Client: httpx.New(address, ""),
	}
}

func (c *IdxClient) Register(
	gtx context.Context, user *core.User, password string) error {
	up := core.UserWithPassword{User: user, Password: password}
	apiRes := c.Build().Path("/api/v1/user").Post(gtx, up)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to register user")
	}
	return nil

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

func (c *IdxClient) AddGroup(
	gtx context.Context, srv int64, gp *core.Group) error {
	// TODO - implement
	return nil
}

func (c *IdxClient) AddUserToGroup(gtx context.Context, uid, gid int64) error {
	// TODO - implement
	return nil
}
