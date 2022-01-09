package httpx

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/sause/libx"
	"github.com/varunamachi/sause/libx/auth"
	"github.com/varunamachi/sause/libx/errx"
)

var (
	ErrHttpParam = errors.New("sause: http param error")
)

type ParamGetter struct {
	etx  echo.Context
	errs map[string]string
}

func NewParamGetter(etx echo.Context) *ParamGetter {
	return &ParamGetter{
		etx: etx,
	}
}

func (pm *ParamGetter) Int(name string) int {
	param := pm.etx.Param(name)
	val, err := strconv.Atoi(param)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return val
}

func (pm *ParamGetter) Int64(name string) int64 {
	param := pm.etx.Param(name)
	val, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return val
}

func (pm *ParamGetter) UInt(name string) uint {
	param := pm.etx.Param(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return uint(val)
}

func (pm *ParamGetter) UInt64(name string) uint64 {
	param := pm.etx.Param(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return val
}

func (pm *ParamGetter) Float64(name string) float64 {
	param := pm.etx.Param(name)
	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return val
}

func (pm *ParamGetter) BoolParam(name string) bool {
	param := pm.etx.Param(name)
	if libx.EqFold(param, "true", "on") {
		return true
	} else if libx.EqFold(param, "false", "off") {
		return false
	}
	pm.errs[name] = "invalid string for bool param"
	return false
}

func (pm *ParamGetter) QueryInt(name string) int {
	param := pm.etx.QueryParam(name)
	val, err := strconv.Atoi(param)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return val
}

func (pm *ParamGetter) QueryInt64(name string) int64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return val
}

func (pm *ParamGetter) QueryUInt(name string) uint {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return uint(val)
}

func (pm *ParamGetter) QueryUInt64Param(name string) uint64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return val
}

func (pm *ParamGetter) QueryFloat64(name string) float64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		pm.errs[name] = err.Error()
	}
	return val
}

func (pm *ParamGetter) QueryBool(name string) bool {
	param := pm.etx.QueryParam(name)
	if libx.EqFold(param, "true", "on") {
		return true
	} else if libx.EqFold(param, "false", "off") {
		return false
	}
	pm.errs[name] = "invalid string for bool param"
	return false
}

func (pm *ParamGetter) QueryIntOr(name string, def int) int {
	param := pm.etx.QueryParam(name)
	val, err := strconv.Atoi(param)
	if err != nil {
		return def
	}
	return val
}

func (pm *ParamGetter) QueryInt64Or(name string, def int64) int64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return def
	}
	return val
}

func (pm *ParamGetter) QueryUIntOr(name string, def uint) uint {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return def
	}
	return uint(val)
}

func (pm *ParamGetter) QueryUInt64Or(name string, def uint64) uint64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return def
	}
	return val
}

func (pm *ParamGetter) QueryFloat64Or(name string, def float64) float64 {
	param := pm.etx.QueryParam(name)
	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return def
	}
	return val
}

func (pm *ParamGetter) QueryBoolOr(name string, def bool) bool {
	param := pm.etx.QueryParam(name)
	if libx.EqFold(param, "true", "on") {
		return true
	} else if libx.EqFold(param, "false", "off") {
		return false
	}
	return def
}

func (pm *ParamGetter) QueryJSON(name string, out interface{}) {
	val := pm.etx.QueryParam(name)
	if len(val) == 0 {
		pm.errs[name] = "could not find json param"
		return
	}
	decoded, err := url.PathUnescape(val)
	if err != nil {
		pm.errs[name] = err.Error()
		return
	}
	if err = json.Unmarshal([]byte(decoded), out); err != nil {
		pm.errs[name] = err.Error()
		return
	}
}

func (pm *ParamGetter) HasError() bool {
	return len(pm.errs) != 0
}

func (pm *ParamGetter) Error() error {
	return errx.Errf(ErrHttpParam, "")
}

func MustGetUser(etx echo.Context) *auth.User {
	obj := etx.Get("user")
	user, ok := obj.(*auth.User)
	if !ok {
		panic("failed to get user info from echo.Context")
	}
	return user
}

func GetUserId(etx echo.Context) string {
	obj := etx.Get("user")
	user, ok := obj.(*auth.User)
	if !ok {
		return ""
	}
	return user.Id
}

func MustGetEndpoint(etx echo.Context) *Endpoint {
	obj := etx.Get("endpoint")
	ep, ok := obj.(*Endpoint)
	if !ok {
		panic("failed to get endpoint info from echo.Context")
	}
	return ep
}

func StrMsg(err *echo.HTTPError) string {
	msg, ok := err.Message.(string)
	if !ok {
		return ""
	}
	return msg
}
