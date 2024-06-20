package tests

import (
	"context"

	"github.com/varunamachi/idx/pg/schema"
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
