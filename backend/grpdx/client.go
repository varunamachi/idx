package grpdx

import (
	"context"
	"time"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

type Client struct {
	*httpx.Client
	Timeout time.Duration
}

func (c *Client) build() *httpx.RequestBuilder {
	builder := c.Build()
	if c.Timeout != 0 {
		builder = builder.WithTimeout(c.Timeout)
	}
	return builder
}

func (c *Client) CreateGroup(
	gtx context.Context,
	group *core.Group) (int64, error) {
	// apiRes := c.Post(gtx, group, "/api/v1/group")
	apiRes := c.build().Path("/api/v1/group").Post(gtx, group)

	res := map[string]int64{"groupId": int64(-1)}
	if err := apiRes.LoadClose(&res); err != nil {
		return -1, errx.Errf(err, "failed to create group: '%s'", group.Name)
	}
	return res["groupId"], nil
}

func (c *Client) UpdateGroup(gtx context.Context, group *core.Group) error {
	// apiRes := c.Put(gtx, group, "/api/v1/group")
	apiRes := c.build().Path("/api/v1/group").Put(gtx, group)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to update group: '%s'", group.Name)

	}
	return nil
}

func (c *Client) GetGroup(
	gtx context.Context, id int64) (*core.Group, error) {
	var group core.Group
	// apiRes := c.Get(gtx, "/api/v1/group", strconv.FormatInt(id, 10))
	apiRes := c.build().Path("/api/v1/group", id).Get(gtx)
	if err := apiRes.LoadClose(&group); err != nil {
		return nil, errx.Errf(err, "failed to get group: '%d'", id)
	}
	return &group, nil
}

func (c *Client) RemoveGroup(gtx context.Context, id int64) error {
	// apiRes := c.Delete(gtx, "/api/v1/group", strconv.FormatInt(id, 10))
	apiRes := c.build().Path("/api/v1/group", id).Delete(gtx)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to delete group: '%d'", id)
	}
	return nil
}

func (c *Client) GetGroups(
	gtx context.Context,
	params *data.CommonParams) ([]*core.Group, error) {
	// TODO -ClientV2 for encoding common params

	groups := make([]*core.Group, 0, 25)
	// apiRes := c.Get(gtx, "/api/v1/group")
	apiRes := c.build().Path("/api/v1/group").CmnParam(params).Get(gtx)
	if err := apiRes.LoadClose(&groups); err != nil {
		return nil, errx.Errf(err, "failed to get groups")
	}

	return groups, nil
}

func (c *Client) GroupExists(gtx context.Context, id int64) (bool, error) {
	res := map[string]bool{
		"exists": false,
	}
	// apiRes := c.Get(
	// 		gtx, "/api/v1/group/", strconv.FormatInt(id, 10), "exists")
	apiRes := c.build().Path("/api/v1/group/", id, "exists").Get(gtx)
	if err := apiRes.LoadClose(&apiRes); err != nil {
		return false, errx.Errf(
			err, "failed to check if group exists: '%d'", id)
	}
	return res["exists"], nil
}

func (c *Client) NumGroups(
	gtx context.Context, filter *data.Filter) (int64, error) {
	// TODO - ClientV2 include filter
	res := map[string]int64{
		"count": 0,
	}
	// apiRes := c.Get(gtx, "/api/v1/group/count")
	apiRes := c.build().Path("/api/v1/group/count").Filter(filter).Get(gtx)
	if err := apiRes.LoadClose(&res); err != nil {
		return 0, errx.Errf(err, "failed to get number of group for filter")
	}
	return res["count"], nil
}

func (c *Client) SetGroupPermissions(
	gtx context.Context,
	groupId int64,
	perms []string) error {
	// apiRes := c.Put(gtx,
	// 	perms,
	// 	"/api/v1/group", strconv.FormatInt(groupId, 10), "perm")
	apiRes := c.build().Path("/api/v1/group", groupId, "perm").Put(gtx, perms)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err,
			"failed to set group's permissions: '%d'", groupId)
	}
	return nil
}

func (c *Client) GetGroupPermissions(
	gtx context.Context,
	groupId int64) ([]string, error) {
	perms := make([]string, 0, 25)
	// apiRes := c.Get(gtx,
	// 	"/api/v1/group", strconv.FormatInt(groupId, 10), "perm")
	apiRes := c.build().Path("/api/v1/group", groupId, "perm").Get(gtx)
	if err := apiRes.LoadClose(&perms); err != nil {
		return nil, errx.Errf(
			err, "failed to get permissions for group: '%d'", groupId)
	}
	return perms, nil
}

func (c *Client) AddUserToGroup(
	gtx context.Context, uid, gid int64) error {
	// TODO - implement
	return nil
}

func (c *Client) RemoveUserFromGroup(
	gtx context.Context, uid, gid int64) error {
	// TODO - implement
	return nil
}
