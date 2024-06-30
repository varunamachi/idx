package tests

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/auth"
	idxAuth "github.com/varunamachi/idx/auth"
	"github.com/varunamachi/idx/controller"
	"github.com/varunamachi/idx/core"
	idxpg "github.com/varunamachi/idx/pg"
	"github.com/varunamachi/idx/pg/schema"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
)

func Setup(gtx context.Context) error {

	if err := runDockerCompose("up", mustGetPgDockerComposePath()); err != nil {
		return err
	}
	if err := schema.Init(gtx, "test"); err != nil {
		return err
	}

	if err := builbuildAndRunServer(); err != nil {
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

func builbuildAndRunServer() error {

	// go build -ldflags "-s -w" -race -o "$root/_local/bin/picl"
	fxRootPath := mustGetSourceRoot()
	goArch := runtime.GOARCH

	cmdDir := filepath.Join(fxRootPath, "cmd", "idx")
	output := filepath.Join(fxRootPath, "_local", "bin", goArch, "picl")

	cmd := exec.Command(
		"go", "build",
		"-ldflags", "-s -w",
		"-o", output,
		cmdDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// cmd.Env = append(os.Environ(), "GOARCH="+goArch)
	if err := cmd.Run(); err != nil {
		const msg = "failed to run go build"
		log.Error().Err(err).Msg(msg)
		return errx.Errf(err, msg)
	}

	// return output, nil

	// TODO - launch this in background or in a goroutine
	cmd, err := startCmd(output, "serve", "--pg-url",
		"postgres://postgres:postgres@localhost:5432/test-data")
	if err != nil {
		return err
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Error().Err(err).Msg("server exited with an error")
			return
		}
		log.Info().Msg("server exited normally")
	}()

	return nil
}
