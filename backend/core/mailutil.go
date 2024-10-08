package core

import (
	"context"

	"github.com/varunamachi/idx/mailtmpl"
	"github.com/varunamachi/libx/data"
	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
	"github.com/varunamachi/libx/str"
)

func SendSimpleMail(
	gtx context.Context,
	to string,
	templateName string,
	data data.M) error {

	from := rt.EnvString("IDX_SENDER_MAIL_ID", "idx-noreply@example.com")

	template, err := mailtmpl.ReadTemplate(templateName)
	if err != nil {
		return errx.Wrap(err)
	}

	content, err := str.SimpleTemplateExpand(&str.TemplateDesc{
		Template: template,
		Data:     data,
		Html:     true,
	})
	if err != nil {
		return errx.Wrap(err)
	}

	err = MailProvider(gtx).Send(&email.Message{
		Id:   to + ":" + templateName,
		From: from,
		To: []string{
			to,
		},
		Content: content,
		Data:    data,
	}, true)

	return errx.Wrap(err)
}
