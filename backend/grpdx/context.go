package grpdx

import "context"

type serviceHolderKey string

const servicesKey = serviceHolderKey("core-services")

func srvs(gtx context.Context) *Services {
	svs := gtx.Value(servicesKey).(*Services)
	if svs == nil {
		panic("failed get service holder from context")
	}
	return svs

}

type Services struct {
	GroupCtlr GroupController
}

func GroupCtlr(gtx context.Context) GroupController {
	return srvs(gtx).GroupCtlr
}

func CopyServices(source, target context.Context) context.Context {
	s := srvs(source)
	return context.WithValue(target, servicesKey, s)
}

func NewContext(gtx context.Context, services *Services) context.Context {
	return context.WithValue(gtx, servicesKey, services)
}
