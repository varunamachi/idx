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

func ServiceEndpoints(us model.ServiceStorage) []*httpx.Endpoint {
	return []*httpx.Endpoint{
		createServiceEp(us),
		updateServiceEp(us),
		getService(us),
		getServices(us),
		deleteService(us),
	}
}

func createServiceEp(us model.ServiceStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var service model.Service
		if err := etx.Bind(&service); err != nil {
			return errx.BadReq("failed to read service info from request", err)
		}

		if err := us.Save(etx.Request().Context(), &service); err != nil {
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
		Permissions: []string{model.PermCreateService},
		Handler:     handler,
	}
}

func updateServiceEp(us model.ServiceStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var service model.Service
		if err := etx.Bind(&service); err != nil {
			return errx.BadReq("failed to read service info from request", err)
		}

		if err := us.Update(etx.Request().Context(), &service); err != nil {
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
		Permissions: []string{model.PermCreateService},
		Handler:     handler,
	}
}

func getService(us model.ServiceStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		service, err := us.GetOne(etx.Request().Context(), id)
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
		Permissions: []string{model.PermGetService},
		Handler:     handler,
	}
}

func getServices(us model.ServiceStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		cmnParams, err := rest.GetCommonParams(etx)
		if err != nil {
			return err
		}

		services, err := us.Get(etx.Request().Context(), cmnParams)
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
		Permissions: []string{model.PermGetService},
		Handler:     handler,
	}
}

func deleteService(us model.ServiceStorage) *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		prmg := httpx.NewParamGetter(etx)
		id := prmg.Int("id")
		if prmg.HasError() {
			return prmg.BadReqError()
		}

		err := us.Remove(etx.Request().Context(), id)
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
		Permissions: []string{model.PermGetService},
		Handler:     handler,
	}
}
