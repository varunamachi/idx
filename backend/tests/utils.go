package tests

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/iox"
)

func execCmd(cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	// cmd.Stdout =

	if err := cmd.Run(); err != nil {
		return errx.Errf(err, "failed to execute '%s'", cmdName)
	}

	return nil
}

func startCmd(cmdName string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(cmdName, args...)
	// cmd.Stdout =

	if err := cmd.Start(); err != nil {
		return nil, errx.Errf(err, "failed to execute '%s'", cmdName)
	}

	return cmd, nil
}

func mustGetSourceRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		const msg = "Couldnt get main file path"
		log.Fatal().Msg(msg)
	}

	fxRootPath, err := filepath.Abs(filename + "/../..")
	if err != nil {
		panic(err)
	}

	subdirs := []string{"backend", "_scripts", "_local"}
	for _, sd := range subdirs {
		if !iox.ExistsAsDir(sd) {
			panic(fmt.Errorf("could not find expected dir '%s' in source root",
				sd))
		}
	}

	return fxRootPath
}

func mustGetPgDockerComposePath() string {
	return filepath.Join(mustGetSourceRoot(),
		"_scripts/deployment/dev/pg.docker-compose.yml")

}
