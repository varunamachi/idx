package grpdx

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/idx/userdx"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/utils/rest"
)

func GroupEndpoints(gtx context.Context) []*httpx.Endpoint {
	gs := core.GroupCtlr(gtx)
	return []*httpx.Endpoint{
		createGroupEp(gs),
		updateGroupEp(gs),
		getGroupEp(gs),
		getGroupsEp(gs),
		deleteGroupEp(gs),
		groupExistsEp(gs),
		getNumGroupsEp(gs),
		setGroupPermissionsEp(gs),
		getGroupPermissionsEp(gs),
		addUserToGroups(gs),
		removeUserFromGroup(gs),
	}
}

func createGroupEp(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var group core.Group
		if err := etx.Bind(&group); err != nil {
			return errx.BadReq("failed to read group info from request", err)
		}

		id, err := gs.Save(etx.Request().Context(), &group)
		if err != nil {
			return errx.Wrap(err)
		}
		return httpx.SendJSON(etx, data.M{
			"groupId": id,
		})
	}

	return &httpx.Endpoint{
		Method:      echo.POST,
		Path:        "/group",
		Category:    "idx.group",
		Desc:        "Create a group",
		Version:     "v1",
		Permissions: []string{userdx.PermCreateUser},
		Handler:     handler,
	}
}

func updateGroupEp(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var group core.Group
		if err := etx.Bind(&group); err != nil {
			return errx.BadReq("failed to read group info from request", err)
		}

		if err := gs.Update(etx.Request().Context(), &group); err != nil {
			return errx.Wrap(err)
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.PUT,
		Path:        "/group",
		Category:    "idx.group",
		Desc:        "Update a group",
		Version:     "v1",
		Permissions: []string{PermCreateGroup},
		Handler:     handler,
	}
}

func getGroupEp(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int64("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		group, err := gs.GetOne(etx.Request().Context(), id)
		if err != nil {
			return errx.Wrap(err)
		}

		return httpx.SendJSON(etx, group)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/group/:id",
		Category:    "idx.group",
		Desc:        "Get info for a group",
		Version:     "v1",
		Permissions: []string{PermGetGroup},
		Handler:     handler,
	}
}

func getGroupsEp(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		cmnParams, err := rest.GetCommonParams(etx)
		if err != nil {
			return errx.Wrap(err)
		}

		groups, err := gs.Get(etx.Request().Context(), cmnParams)
		if err != nil {
			return errx.Wrap(err)
		}

		return httpx.SendJSON(etx, groups)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/group",
		Category:    "idx.group",
		Desc:        "Get group list",
		Version:     "v1",
		Permissions: []string{PermGetGroup},
		Handler:     handler,
	}
}

func deleteGroupEp(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int64("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		err := gs.Remove(etx.Request().Context(), id)
		if err != nil {
			return errx.Wrap(err)
		}

		return etx.String(http.StatusOK, strconv.FormatInt(id, 10))
	}

	return &httpx.Endpoint{
		Method:      echo.DELETE,
		Path:        "/group/:id",
		Category:    "idx.group",
		Desc:        "Delete a group",
		Version:     "v1",
		Permissions: []string{PermDeleteGroup},
		Handler:     handler,
	}
}

func groupExistsEp(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int64("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		group, err := gs.Exists(etx.Request().Context(), id)
		if err != nil {
			return errx.Wrap(err)
		}

		return httpx.SendJSON(etx, group)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/group/:id/exists",
		Category:    "idx.group",
		Desc:        "Check if a group exists",
		Version:     "v1",
		Permissions: []string{PermGetGroup},
		Handler:     handler,
	}
}

func getNumGroupsEp(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		filter, err := rest.GetFilter(etx)
		if err != nil {
			return errx.BadReqX(err, "failed decode filter from request")
		}

		count, err := gs.Count(etx.Request().Context(), filter)
		if err != nil {
			return errx.Wrap(err)
		}

		return httpx.SendJSON(etx, map[string]int64{
			"count": count,
		})
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/group/count",
		Category:    "idx.group",
		Desc:        "Get number of groups, selected by given filter",
		Version:     "v1",
		Permissions: []string{PermGetGroup},
		Handler:     handler,
	}
}

func setGroupPermissionsEp(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {

		pmg := httpx.NewParamGetter(etx)
		groupId := pmg.Int64("groupId")
		if pmg.HasError() {
			return pmg.BadReqError()
		}

		perms := make([]string, 0, 100)
		if err := etx.Bind(&perms); err != nil {
			return errx.BadReq("failed to read group info from request", err)
		}

		err := gs.SetPermissions(etx.Request().Context(), groupId, perms)
		if err != nil {
			return errx.Wrap(err)
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.PUT,
		Path:        "/group/:groupId/perm",
		Category:    "idx.group",
		Desc:        "Set permissions for a group",
		Version:     "v1",
		Permissions: []string{PermModifyGroupPerm},
		Handler:     handler,
	}
}

func getGroupPermissionsEp(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int64("groupId")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		perms, err := gs.GetPermissions(etx.Request().Context(), id)
		if err != nil {
			return errx.Wrap(err)
		}

		return httpx.SendJSON(etx, perms)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/group/:groupId/perm",
		Category:    "idx.group",
		Desc:        "Get permissions of a group",
		Version:     "v1",
		Permissions: []string{PermGetGroup},
		Handler:     handler,
	}
}

func addUserToGroups(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {

		pmg := httpx.NewParamGetter(etx)
		userId := pmg.Int64("userId")
		if pmg.HasError() {
			return pmg.BadReqError()
		}

		groupIds := make([]int64, 0, 100)
		if err := etx.Bind(&groupIds); err != nil {
			return errx.BadReq("failed to read group info from request", err)
		}

		err := gs.AddToGroups(etx.Request().Context(), userId, groupIds...)
		if err != nil {
			return errx.Wrap(err)
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.PUT,
		Path:        "/group/:userId",
		Category:    "idx.group",
		Desc:        "Add user to ne or more groups",
		Version:     "v1",
		Permissions: []string{PermModifyGroupPerm},
		Handler:     handler,
	}
}

func removeUserFromGroup(gs core.GroupController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		groupId := prmg.Int64("groupId")
		userId := prmg.Int64("userId")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		err := gs.RemoveFromGroup(etx.Request().Context(), userId, groupId)
		if err != nil {
			return errx.Wrap(err)
		}

		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.DELETE,
		Path:        "/group/:groupId/:userId",
		Category:    "idx.group",
		Desc:        "Delete a group",
		Version:     "v1",
		Permissions: []string{PermDeleteGroup},
		Handler:     handler,
	}
}
