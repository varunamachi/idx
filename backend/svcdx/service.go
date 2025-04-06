package svcdx

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
func (c *Client) GetPerms(
	gtx context.Context, serviceId, userId int) ([]string, error) {
	return nil, nil
}

func (c *Client) CreateService(
	gtx context.Context, srv *core.Service) (int64, error) {

	apiRes := c.build().Path("/api/v1/service").Post(gtx, srv)
	res := map[string]int64{"serviceId": int64(-1)}
	if err := apiRes.LoadClose(&res); err != nil {
		return -1, errx.Errf(err, "failed to cratre service: '%s'", srv.Name)
	}
	return res["serviceId"], nil
}

func (c *Client) UpdateService(
	gtx context.Context, srv *core.Service) error {
	apiRes := c.build().Path("/api/v1/service").Put(gtx, srv)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to update service: '%s'", srv.Name)
	}
	return nil
}

func (c *Client) GetService(
	gtx context.Context, id int64) (*core.Service, error) {
	apiRes := c.build().Path("/api/v1/service", id).Get(gtx)
	var service core.Service
	if err := apiRes.LoadClose(&service); err != nil {
		return nil, errx.Errf(err, "failed to get service: '%d'", id)
	}
	return &service, nil
}

func (c *Client) GetServices(
	gtx context.Context, params *data.CommonParams) ([]*core.Service, error) {

	// TODO - use common params as JSON query param
	apiRes := c.build().Path("/api/v1/service").CmnParam(params).Get(gtx)
	services := make([]*core.Service, 0, params.PageSize)
	if err := apiRes.LoadClose(&services); err != nil {
		return nil, errx.Errf(err, "failed to get list of services")
	}
	return services, nil
}

func (c *Client) RemoveService(gtx context.Context, id int64) error {
	// apiRes := c.Get(gtx, "/api/v1/service", strconv.FormatInt(id, 10))
	apiRes := c.build().Path("/api/v1/service", id).Delete(gtx)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to remove service: '%d'", id)
	}
	return nil
}

func (c *Client) ServiceExists(
	gtx context.Context, name string) (bool, error) {
	// apiRes := c.Get(gtx, "/api/v1/service/exists", name)
	apiRes := c.build().Path("/api/v1/service/exists", name).Get(gtx)
	exists := map[string]bool{
		"exists": false,
	}
	if err := apiRes.LoadClose(&exists); err != nil {
		return false, errx.Errf(
			err, "failed to check if service exists: '%s'", name)
	}
	return exists["exists"], nil
}

func (c *Client) NumServices(
	gtx context.Context, filter *data.Filter) (int64, error) {
	// apiRes := c.Get(gtx, "/api/v1/service/count")
	apiRes := c.build().Path("/api/v1/service/count").Filter(filter).Get(gtx)
	num := map[string]int64{
		"count": 0,
	}
	if err := apiRes.LoadClose(&num); err != nil {
		return 0, errx.Errf(
			err, "failed to get service count for a filter")
	}
	return num["count"], nil
}

func (c *Client) GetServiceByName(
	gtx context.Context, name string) (*core.Service, error) {
	// apiRes := c.Get(gtx, "/api/v1/service/named", name)
	apiRes := c.build().Path("/api/v1/service/named", name).Get(gtx)
	var service core.Service
	if err := apiRes.LoadClose(&service); err != nil {
		return nil, errx.Errf(err, "failed to get service: '%s'", name)
	}
	return &service, nil
}

func (c *Client) GetServicesForOwner(
	gtx context.Context, ownerId string) ([]*core.Service, error) {
	// apiRes := c.Get(gtx, "/api/v1/service/owner", ownerId)
	apiRes := c.build().Path("/api/v1/service/owner", ownerId).Get(gtx)
	services := make([]*core.Service, 0, 10)
	if err := apiRes.LoadClose(&services); err != nil {
		return nil, errx.Errf(
			err, "failed to get services owned by '%s'", ownerId)
	}
	return services, nil
}

func (c *Client) AddAdminToService(
	gtx context.Context, serviceId, userId int64) error {
	// apiRes := c.Put(gtx, nil,
	// 	"/api/v1/service/", strconv.FormatInt(serviceId, 10),
	// 	"admin", strconv.FormatInt(userId, 10))

	apiRes := c.build().
		Path("/api/v1/service/", serviceId, "admin", userId).
		Put(gtx, nil)

	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to add admin '%d' to service '%d'",
			serviceId, userId)
	}

	return nil
}

func (c *Client) GetServiceAdmins(
	gtx context.Context, serviceId int64) ([]*core.User, error) {
	// apiRes := c.Get(gtx, "/api/v1/service/", strconv.FormatInt(serviceId, 10),
	// 	"admin")

	apiRes := c.build().Path(gtx, "/api/v1/service/", serviceId).Get(gtx)
	admins := make([]*core.User, 0, 10)
	if err := apiRes.LoadClose(&admins); err != nil {
		return nil, errx.Errf(
			err, "failed to get admins of service '%d'", serviceId)
	}
	return admins, nil
}

func (c *Client) RemoveAdminFromService(
	gtx context.Context, serviceId, userId int64) error {
	// apiRes := c.Delete(gtx,
	// 	"/api/v1/service/", strconv.FormatInt(serviceId, 10),
	// 	"admin", strconv.FormatInt(userId, 10))
	apiRes := c.build().
		Path("/api/v1/service/", serviceId, "admin", userId).
		Delete(gtx)
	if err := apiRes.Close(); err != nil {
		return errx.Errf(err, "failed to delete admin '%d' to service '%d'",
			serviceId, userId)
	}

	return nil
}

func (c *Client) IsServiceAdmin(
	gtx context.Context, serviceId, userId int64) (bool, error) {
	// apiRes := c.Get(gtx,
	// 	"/api/v1/service/", strconv.FormatInt(serviceId, 10),
	// 	"admin", strconv.FormatInt(userId, 10),
	// 	"exists")

	apiRes := c.build().
		Path("/api/v1/service/", serviceId, "admin", userId, "exists").
		Get(gtx)
	out := map[string]bool{
		"isAdmin": false,
	}
	if err := apiRes.LoadClose(&out); err != nil {
		return false, errx.Errf(
			err, "failed to get admins of service '%d'", serviceId)
	}
	return out["isAdmin"], nil
}

func (c *Client) GetUserPermsForService(
	gtx context.Context, serviceId, userId int64) ([]string, error) {
	// TODO - implement
	return nil, nil
}
