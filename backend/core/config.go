package core

import (
	"fmt"
	"net/url"
	"path"

	"github.com/varunamachi/libx/rt"
)

func ToFullUrl(pathElements ...string) string {
	baseUrl := rt.EnvString("IDX_BASE_URL", "http://localhost:8080")
	url, err := url.Parse(baseUrl)
	if err != nil {
		msg := fmt.Sprintf("unvalid base url '%s' given", baseUrl)
		panic(msg)
	}

	url.Path = path.Join(pathElements...)
	return url.String()
}
