package client

import (
	"context"

	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/data"
)

func (c *IdxClient) CreateGroup(
	gtx context.Context,
	group *core.Group) error {
	return nil
}

func (c *IdxClient) CreateGroupWithPerms(
	gtx context.Context,
	group *core.Group,
	perms []string) error {
	return nil
}

func (c *IdxClient) UpdateGroup(gtx context.Context, group *core.Group) error {
	return nil
}

func (c *IdxClient) GetGroup(gtx context.Context, id int64) (*core.Group, error) {
	return nil, nil
}

func (c *IdxClient) RemoveGroup(gtx context.Context, id int64) error {
	return nil
}

func (c *IdxClient) GetGroups(
	gtx context.Context,
	params *data.CommonParams) ([]*core.Group, error) {
	return nil, nil
}

func (c *IdxClient) GroupExists(gtx context.Context, id int64) (bool, error) {
	return false, nil
}

func (c *IdxClient) NumGroups(
	gtx context.Context, filter *data.Filter) (int64, error) {
	return 0, nil
}

func (c *IdxClient) SetGroupPermissions(
	gtx context.Context,
	groupId int64,
	perms []string) error {
	return nil
}

func (c *IdxClient) GetGroupPermissions(
	gtx context.Context,
	groupId int64) ([]string, error) {
	return nil, nil
}
