package tests

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/tests/smsrv"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
)

type TestConfig struct {
	SkipAppServer     bool
	SkipDockerCompose bool
	DestroySchema     bool
}

const pgUrl = "postgresql://idx:idxp@localhost:5432/test-data?sslmode=disable"

func Setup(gtx context.Context, testConfig *TestConfig) (*os.Process, error) {

	if !testConfig.SkipDockerCompose {
		err := RunDockerCompose("up", MustGetPgDockerComposePath())
		if err != nil {
			return nil, err
		}

		purl, err := url.Parse(pgUrl)
		if err != nil {
			return nil, errx.Errf(err, "invalid pg-url in setup")
		}
		db, err := pg.Connect(gtx, purl, "")
		if err != nil {
			return nil, err
		}
		pg.SetDefaultConn(db)
	}

	// init is done during server start
	// if err := schema.Init(gtx, "test"); err != nil {
	// 	return err
	// }

	if !testConfig.SkipAppServer {
		smsrv.GetMailService().Start(gtx)
		process, err := BuildAndRunServer()
		if err != nil {
			return nil, err
		}
		return process, nil
	}

	return nil, nil
}

func Destroy(
	gtx context.Context,
	testConfig *TestConfig,
	serverProcess *os.Process) error {

	// if err := schema.Destroy(gtx); err != nil {
	// 	log.Error().Err(err)
	// }

	if !testConfig.SkipDockerCompose {
		err := RunDockerCompose("down", MustGetPgDockerComposePath())
		if err != nil {
			return err
		}
	}

	if !testConfig.SkipAppServer {
		if err := serverProcess.Signal(os.Interrupt); err != nil {
			return errx.Errf(err, "failed to send SIGINT to server process")
		}
		if _, err := serverProcess.Wait(); err != nil {
			return errx.Errf(err, "waiting for server process failed")
		}
	}

	return nil
}

func RunDockerCompose(op, dcFilePath string) error {
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

func BuildAndRunServer() (*os.Process, error) {
	fxRootPath := MustGetSourceRoot()
	goArch := runtime.GOARCH

	cmdDir := filepath.Join(fxRootPath, "backend", "cmd", "idx")
	output := filepath.Join(fxRootPath, "_local", "bin", goArch, "idx")

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
		return nil, errx.Errf(err, msg)
	}

	builder := rt.NewCmdBuilder(output).
		WithArgs("serve", "--pg-url", pgUrl).
		WithEnv("IDX_MAIL_PROVIDER", "IDX_SIMPLE_MAIL_SERVICE_CLIENT_PROVIDER").
		WithEnv("IDX_SIMPLE_SRV_SEND_URL", "http://localhost:9999/send")

	process, err := builder.Start()
	if err != nil {
		return nil, errx.Errf(err, "server exited with an error")
	}
	log.Info().Msg("server exited normally")

	return process, nil
}
