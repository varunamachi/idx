package restapi

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/model"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/utils/rest"
)

func GroupEndpoints(gs model.GroupStorage) []*httpx.Endpoint {
	return []*httpx.Endpoint{
		createGroupEp(gs),
		updateGroupEp(gs),
		getGroup(gs),
		getGroups(gs),
		deleteGroup(gs),
	}
}

func createGroupEp(gs model.GroupStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var group model.Group
		if err := etx.Bind(&group); err != nil {
			return errx.BadReq("failed to read group info from request", err)
		}

		if err := gs.Save(etx.Request().Context(), &group); err != nil {
			return err
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.POST,
		Path:        "/group",
		Category:    "idx.group",
		Desc:        "Create a group",
		Version:     "v1",
		Permissions: []string{model.PermCreateUser},
		Handler:     handler,
	}
}

func updateGroupEp(gs model.GroupStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var group model.Group
		if err := etx.Bind(&group); err != nil {
			return errx.BadReq("failed to read group info from request", err)
		}

		if err := gs.Update(etx.Request().Context(), &group); err != nil {
			return err
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.PUT,
		Path:        "/group",
		Category:    "idx.group",
		Desc:        "Update a group",
		Version:     "v1",
		Permissions: []string{model.PermCreateGroup},
		Handler:     handler,
	}
}

func getGroup(gs model.GroupStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		group, err := gs.GetOne(etx.Request().Context(), id)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, group)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/group/:id",
		Category:    "idx.group",
		Desc:        "Get info for a group",
		Version:     "v1",
		Permissions: []string{model.PermGetGroup},
		Handler:     handler,
	}
}

func getGroups(gs model.GroupStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		cmnParams, err := rest.GetCommonParams(etx)
		if err != nil {
			return err
		}

		groups, err := gs.Get(etx.Request().Context(), cmnParams)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, groups)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/group",
		Category:    "idx.group",
		Desc:        "Get group list",
		Version:     "v1",
		Permissions: []string{model.PermGetGroup},
		Handler:     handler,
	}
}

func deleteGroup(gs model.GroupStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		err := gs.Remove(etx.Request().Context(), id)
		if err != nil {
			return err
		}

		return etx.String(http.StatusOK, strconv.Itoa(id))
	}

	return &httpx.Endpoint{
		Method:      echo.DELETE,
		Path:        "/group/:id",
		Category:    "idx.group",
		Desc:        "Delete a group",
		Version:     "v1",
		Permissions: []string{model.PermDeleteGroup},
		Handler:     handler,
	}
}
