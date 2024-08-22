package svcdx

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

func NewContext(gtx context.Context, services *Services) context.Context {
	return context.WithValue(gtx, servicesKey, services)
}

type Services struct {
	ServiceCtlr ServiceController
}

func ServiceCtlr(gtx context.Context) ServiceController {
	return srvs(gtx).ServiceCtlr
}

func CopyServices(source, target context.Context) context.Context {
	s := srvs(source)
	return context.WithValue(target, servicesKey, s)
}
