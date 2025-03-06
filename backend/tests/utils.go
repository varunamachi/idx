package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/iox"
)

func execCmd(cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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

func MustGetSourceRoot() string {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		const msg = "This command must be executed in source, path error"
		log.Fatal().Msg(msg)
	}

	if !strings.Contains(filename, "idx/backend") {
		const msg = "This command must be executed in source, path error"
		log.Fatal().Str("path", filename).Msg(msg)
	}

	idx := strings.Index(filename, "idx/backend")
	fxRootPath := filename[:idx+len("idx")]

	subdirs := []string{"backend/core", "_scripts", "_local"}
	for _, sd := range subdirs {
		if !iox.ExistsAsDir(filepath.Join(fxRootPath, sd)) {
			panic(fmt.Errorf("could not find expected dir '%s' in source root",
				sd))
		}
	}

	return fxRootPath
}

func MustGetPgComposePath() string {
	return filepath.Join(MustGetSourceRoot(), "backend/tests/pg-dc.yaml")

}
