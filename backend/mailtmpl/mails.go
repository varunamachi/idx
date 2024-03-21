package mailtmpl

import (
	"embed"
	"sync"

	"github.com/varunamachi/libx/errx"
)

//go:embed resources/*
var mailDir embed.FS

var cache = struct {
	sync.Mutex
	mp map[string]string
}{
	mp: make(map[string]string),
}

func UserAccountVerificationTemplate() (string, error) {
	return readTemplate("verify_user_account")
}

func UserAccountVerificationSuccessTemplate() (string, error) {
	return readTemplate("user_account_verified")
}

func UserAccountLockedTemplate() (string, error) {
	return readTemplate("user_account_locked")
}

func readTemplate(name string) (string, error) {
	cache.Lock()
	defer cache.Unlock()

	if val, found := cache.mp[name]; found {
		return val, nil
	}

	dt, err := mailDir.ReadFile(name + ".tmpl.html")
	if err != nil {
		return "", errx.Errf(err,
			"failed read embedded mail template: '%s'", name)
	}
	val := string(dt)
	cache.mp[name] = val
	return val, nil
}
