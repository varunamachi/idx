package client

import (
	"time"

	"github.com/varunamachi/idx/grpdx"
	"github.com/varunamachi/idx/svcdx"
	"github.com/varunamachi/idx/userdx"
	"github.com/varunamachi/libx/httpx"
)

type (
	UserClient = userdx.Client
	GrpClient  = grpdx.Client
	SvcClient  = svcdx.Client
)

type Client struct {
	UserClient
	GrpClient
	SvcClient
}

func New(address string) *Client {
	hxClient := httpx.NewClient(address, "")

	return &Client{
		UserClient: userdx.Client{
			Client: hxClient,
		},
		GrpClient: grpdx.Client{
			Client: hxClient,
		},
		SvcClient: svcdx.Client{
			Client: hxClient,
		},
	}
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	// c.timeout = timeout
	c.UserClient.Timeout = timeout
	c.GrpClient.Timeout = timeout
	c.SvcClient.Timeout = timeout
	return c
}
