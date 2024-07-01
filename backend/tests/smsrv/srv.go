package smsrv

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
)

type MailService struct {
	fakeProvider email.FakeEmailProvider
}

func (ms *MailService) sendEp() *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		var msg email.Message
		if err := etx.Bind(&msg); err != nil {
			return errx.BadReqX(err, "failed to read mail msg")
		}

		err := ms.fakeProvider.Send(&msg, false)
		return err
	}

	return &httpx.Endpoint{
		Method:  echo.POST,
		Path:    "/send",
		Version: "v1",
		Handler: handler,
	}
}

func (ms *MailService) getEp() *httpx.Endpoint {
	handler := func(etx echo.Context) error {
		user := etx.Param("user")
		recepType := etx.Param("recepType")
		mailId := etx.Param("mailId")

		msg, err := ms.Get(user, recepType, mailId)
		if err != nil {
			return err
		}

		return httpx.SendJSON(etx, msg)
	}

	return &httpx.Endpoint{
		Method:  echo.GET,
		Path:    "/mail/:user/:recepType/:mailId",
		Version: "v1",
		Handler: handler,
	}
}

func (ms *MailService) Start(gtx context.Context) {

	go func() {
		err := httpx.NewServer(nil, nil).WithAPIs(
			ms.sendEp(),
			ms.getEp(),
		).Start(9999)
		// TODO - context, cancel, exit etc
		fmt.Println(err)
	}()
}

func (ms *MailService) Get(
	user, recepType, mailId string) (*email.Message, error) {

	switch recepType {
	case "to":
		return ms.fakeProvider.Get(user, mailId)
	case "cc":
		return ms.fakeProvider.GetCC(user, mailId)
	case "bcc":
		return ms.fakeProvider.GetBCC(user, mailId)
	default:
		return ms.fakeProvider.GetAny(user, mailId)
	}
}
