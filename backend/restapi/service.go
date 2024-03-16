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
		getService(ss),
		getServices(ss),
		deleteService(ss),
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

func getService(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int("id")
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

func getServices(ss core.ServiceController) *httpx.Endpoint {
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

func deleteService(ss core.ServiceController) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		err := ss.Remove(etx.Request().Context(), id)
		if err != nil {
			return err
		}

		return etx.String(http.StatusOK, strconv.Itoa(id))
	}

	return &httpx.Endpoint{
		Method:      echo.DELETE,
		Path:        "/service/:id",
		Category:    "idx.service",
		Desc:        "Delete a service",
		Version:     "v1",
		Permissions: []string{core.PermGetService},
		Handler:     handler,
	}
}
