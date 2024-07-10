package tests

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/pg/schema"
	"github.com/varunamachi/idx/tests/smsrv"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
)

const pgUrl = "postgres://postgres:postgres@localhost:5432/test-data?sslmode=disable"

func Setup(gtx context.Context) error {

	if err := runDockerCompose("up", mustGetPgDockerComposePath()); err != nil {
		return err
	}

	purl, err := url.Parse(pgUrl)
	if err != nil {
		return errx.Errf(err, "invalid pg-url in setup")
	}
	db, err := pg.Connect(gtx, purl, "Asia/Kolkata")
	if err != nil {
		return err
	}
	pg.SetDefaultConn(db)

	// init is done during server start
	// if err := schema.Init(gtx, "test"); err != nil {
	// 	return err
	// }

	smsrv.GetMailService().Start(gtx)

	if err := buildAndRunServer(); err != nil {
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
		"compose",
		"-p",
		"idx_test",
		"-f",
		dcFilePath,
		op,
		data.Qop(op == "up", "-d", ""),
	}
	return execCmd("docker", args...)
}

func buildAndRunServer() error {
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

	builder := rt.NewCmdBuilder(output).
		WithArgs("serve", "--pg-url", pgUrl).
		WithEnv("IDX_MAIL_PROVIDER", "IDX_SIMPLE_MAIL_SERVICE_CLIENT_PROVIDER").
		WithEnv("IDX_SIMPLE_SRV_SEND_URL", "http://localhost:9999/send")

	go func() {
		if err := builder.Execute(); err != nil {
			log.Error().Err(err).Msg("server exited with an error")
			return
		}
		log.Info().Msg("server exited normally")
	}()

	return nil
}
