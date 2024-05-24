package restapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/utils/rest"
)

func ServiceEndpoints(gtx context.Context) []*httpx.Endpoint {
	ss := core.ServiceCtlr(gtx)
	return []*httpx.Endpoint{
		createServiceEp(ss),
		updateServiceEp(ss),
		getServiceEp(ss),
		getServicesEp(ss),
		deleteServiceEp(ss),
		serviceExistsEp(ss),
		numServices(ss),
		serviceByNameEp(ss),
		servicesForOwnerEp(ss),
		addAdminToServiceEp(ss),
		removeAdminFromServiceEp(ss),
		getServiceAdminsEp(ss),
		isServiceAdminEp(ss),
		getPermissionsForService(ss),
	}
}

func createServiceEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var service core.Service
		if err := etx.Bind(&service); err != nil {
			return errx.BadReq("failed to read service info from request", err)
		}

		if err := ss.Save(etx.Request().Context(), &service); err != nil {
			return err
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.POST,
		Path:        "/service",
		Category:    "idx.service",
		Desc:        "Create a service",
		Version:     "v1",
		Permissions: []string{core.PermCreateService},
		Handler:     handler,
	}
}

func updateServiceEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var service core.Service
		if err := etx.Bind(&service); err != nil {
			return errx.BadReq("failed to read service info from request", err)
		}

		if err := ss.Update(etx.Request().Context(), &service); err != nil {
			return err
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.PUT,
		Path:        "/service",
		Category:    "idx.service",
		Desc:        "Update a service",
		Version:     "v1",
		Permissions: []string{core.PermCreateService},
		Handler:     handler,
	}
}

func getServiceEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int64("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		service, err := ss.GetOne(etx.Request().Context(), id)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, service)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/service/:id",
		Category:    "idx.service",
		Desc:        "Get info for a service",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}

func getServicesEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		cmnParams, err := rest.GetCommonParams(etx)
		if err != nil {
			return err
		}

		services, err := ss.Get(etx.Request().Context(), cmnParams)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, services)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/service",
		Category:    "idx.service",
		Desc:        "Get service list",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}

func deleteServiceEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int64("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		err := ss.Remove(etx.Request().Context(), id)
		if err != nil {
			return err
		}

		return etx.String(http.StatusOK, strconv.FormatInt(id, 10))
	}

	return &httpx.Endpoint{
		Method:      echo.DELETE,
		Path:        "/service/:id",
		Category:    "idx.service",
		Desc:        "Delete a service",
		Version:     "v1",
		Permissions: []string{core.PermDeleteService},
		Handler:     handler,
	}
}

func serviceExistsEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		name := prmg.Str("name")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		exists, err := ss.Exists(etx.Request().Context(), name)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, map[string]bool{
			"exists": exists,
		})
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/service/exists/:name",
		Category:    "idx.service",
		Desc:        "Check if service exists",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}

func numServices(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		filter, err := rest.GetFilter(etx)
		if err != nil {
			return err
		}

		services, err := ss.Count(etx.Request().Context(), filter)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, services)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/service/count",
		Category:    "idx.service",
		Desc:        "Get service count for the filter",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}

func serviceByNameEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		name := prmg.Str("name")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		srv, err := ss.GetByName(etx.Request().Context(), name)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, srv)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/service/named/:name",
		Category:    "idx.service",
		Desc:        "Get service by name ",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}

func servicesForOwnerEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		owner := prmg.Str("owner")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		srv, err := ss.GetForOwner(etx.Request().Context(), owner)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, srv)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/service/owner/:owner",
		Category:    "idx.service",
		Desc:        "Get services for an owner",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}

func addAdminToServiceEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		pmg := httpx.NewParamGetter(etx)
		service := pmg.Int64("service")
		user := pmg.Int64("user")

		err := ss.AddAdmin(etx.Request().Context(), service, user)
		if err != nil {
			return err
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.PUT,
		Path:        "/service/:service/admin/:user",
		Category:    "idx.service",
		Desc:        "Add an admin to a service",
		Version:     "v1",
		Permissions: []string{core.PermServiceManageAdmins},
		Handler:     handler,
	}
}

func removeAdminFromServiceEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		pmg := httpx.NewParamGetter(etx)
		service := pmg.Int64("service")
		user := pmg.Int64("user")

		err := ss.RemoveAdmin(etx.Request().Context(), service, user)
		if err != nil {
			return err
		}
		return nil
	}

	return &httpx.Endpoint{
		Method:      echo.DELETE,
		Path:        "/service/:service/admin/:user",
		Category:    "idx.service",
		Desc:        "Remove an admin from a service",
		Version:     "v1",
		Permissions: []string{core.PermServiceManageAdmins},
		Handler:     handler,
	}
}

func getServiceAdminsEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		service := prmg.Int64("service")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		srv, err := ss.GetAdmins(etx.Request().Context(), service)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, srv)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/service/:service/admin",
		Category:    "idx.service",
		Desc:        "Get admins of a service",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}

func isServiceAdminEp(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		service := prmg.Int64("service")
		user := prmg.Int64("user")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		is, err := ss.IsAdmin(etx.Request().Context(), service, user)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, map[string]bool{
			"isAdmin": is,
		})
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/service/:service/admin/:user/exists",
		Category:    "idx.service",
		Desc:        "Check if user is an admin of a service",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}

func getPermissionsForService(gs core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		userId := prmg.Int64("userId")
		serviceId := prmg.Int64("serviceId")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		perms, err := gs.GetPermissionForService(
			etx.Request().Context(), userId, serviceId)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, perms)
	}

	return &httpx.Endpoint{
		Method:      echo.GET,
		Path:        "/service/:serviceId/perms/:userId",
		Category:    "idx.service",
		Desc:        "Get permissions of a service for a user",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}
