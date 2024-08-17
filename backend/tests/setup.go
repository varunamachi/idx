package tests

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/pg/schema"
	"github.com/varunamachi/idx/tests/smsrv"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/netx"
	"github.com/varunamachi/libx/rt"
)

type TestConfig struct {
	SkipAppServer     bool
	SkipDockerCompose bool
	SkipDataClean     bool
	DestroySchema     bool
}

const pgUrl = "postgresql://idx:idxp@localhost:5432/test-data?sslmode=disable"

func Setup(gtx context.Context, testConfig *TestConfig) (*os.Process, error) {

	if !testConfig.SkipDockerCompose {
		err := RunDockerCompose("up", MustGetPgDockerComposePath())
		if err != nil {
			return nil, err
		}
	} else {
		// Wait for externally launched postgres server
		log.Info().Msg("waiting for postgres at port 5432 (max wait = 2m)")
		err := netx.WaitForPorts(gtx, "localhost:5432", 2*time.Minute)
		if err != nil {
			return nil, errx.Errf(err, "could not connect to postgres server")
		}
	}

	// init is done during server start
	// if err := schema.Init(gtx, "test"); err != nil {
	// 	return errx.Wrap(err)
	// }
	if err := ConnectToTestDB(gtx); err != nil {
		return nil, errx.Wrap(err)
	}

	smsrv.GetMailService().Start(gtx)
	if !testConfig.SkipAppServer {
		process, err := BuildAndRunServer()
		if err != nil {
			return nil, err
		}
		return process, nil
	} else {
		// Wait for externally launched server
		log.Info().Msg("waiting for app server at port 8888 (max wait = 2m)")
		err := netx.WaitForPorts(gtx, "localhost:8888", 2*time.Minute)
		if err != nil {
			return nil, errx.Errf(err, "could not connect to app server")
		}
	}

	return nil, nil
}

func Destroy(
	gtx context.Context,
	testConfig *TestConfig,
	serverProcess *os.Process) error {

	if !testConfig.SkipDockerCompose {
		err := RunDockerCompose("down", MustGetPgDockerComposePath())
		if err != nil {
			return errx.Wrap(err)
		}
		log.Info().Msg("docker-compose shutdown complete")
	}

	if !testConfig.SkipAppServer {
		if err := serverProcess.Signal(os.Interrupt); err != nil {
			return errx.Errf(err, "failed to send SIGINT to server process")
		}
		if _, err := serverProcess.Wait(); err != nil {
			return errx.Errf(err, "waiting for server process failed")
		}
		log.Info().Msg("app-server shutdown complete")
	}

	if !testConfig.SkipDataClean {
		if err := schema.CleanData(gtx); err != nil {
			fmt.Println(err)
			return errx.Wrap(err)
		}
		log.Info().Msg("data clean complete")
	}

	if testConfig.DestroySchema {
		if err := schema.Destroy(gtx); err != nil {
			return errx.Wrap(err)
		}
		log.Info().Msg("schema delete complete")
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
		WithEnv("IDX_SIMPLE_SRV_SEND_URL", "http://localhost:9999/api/v1/send").
		WithEnv("IDX_ROLE_MAPPING", "super:Super")

	process, err := builder.Start()
	if err != nil {
		return nil, errx.Errf(err, "server exited with an error")
	}
	log.Info().Msg("server exited normally")

	return process, nil
}

func ConnectToTestDB(gtx context.Context) error {
	purl, err := url.Parse(pgUrl)
	if err != nil {
		return errx.Errf(err, "invalid pg-url in setup")
	}
	db, err := pg.Connect(gtx, purl, "")
	if err != nil {
		return errx.Wrap(err)
	}
	pg.SetDefaultConn(db)
	return nil
}
