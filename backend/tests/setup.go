package tests

import (
	"context"

	"github.com/varunamachi/idx/auth"
	idxAuth "github.com/varunamachi/idx/auth"
	"github.com/varunamachi/idx/controller"
	"github.com/varunamachi/idx/core"
	idxpg "github.com/varunamachi/idx/pg"
	"github.com/varunamachi/idx/pg/schema"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/rt"
)

func Setup(gtx context.Context) error {

	if err := runDockerCompose("up", mustGetPgDockerComposePath()); err != nil {
		return err
	}
	if err := schema.Init(gtx, "test"); err != nil {
		return err
	}

	return nil
}

func Destroy(gtx context.Context) error {

	if err := schema.Destroy(gtx); err != nil {
		return err
	}
	err := runDockerCompose("down", mustGetPgDockerComposePath())
	if err != nil {
		return err
	}

	return nil
}

func runDockerCompose(op, dcFilePath string) error {
	args := []string{
		"-p",
		"idx_test",
		op,
		"-f",
		dcFilePath,
	}
	return execCmd("docker-compose", args...)
}

func Gtx() (context.Context, context.CancelFunc) {

	// TODO - will these work if pg database is initialize later?
	gtx, cancel := rt.Gtx()

	emailProvider := email.NewFakeEmailProvider()
	evtSrv := idxpg.NewEventService()

	gd := pg.NewGetterDeleter()
	userStore := idxpg.NewUserStorage(gd)
	groupStore := idxpg.NewGroupStorage(gd)
	serviceStore := idxpg.NewServiceStorage(gd)

	hasher := auth.NewArgon2Hasher()
	credStorage := idxpg.NewCredentialStorage(hasher)

	authr := idxAuth.NewAuthenticator(credStorage)
	uctlr := controller.NewUserController(userStore, credStorage, emailProvider)
	gctlr := controller.NewGroupController(groupStore, serviceStore)
	sctlr := controller.NewServiceController(
		serviceStore, userStore, groupStore)

	return core.NewContext(gtx, &core.Services{
		UserCtlr:      uctlr,
		ServiceCtlr:   sctlr,
		GroupCtlr:     gctlr,
		Authenticator: authr,
		MailProvider:  emailProvider,
		EventService:  evtSrv,
	}), cancel

}
