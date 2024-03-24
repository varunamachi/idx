package core

import (
	"context"

	"github.com/varunamachi/libx/email"
	"github.com/varunamachi/libx/rt"
	"github.com/varunamachi/libx/str"
)

func SendSimpleMail(
	gtx context.Context, to string, template string, data any) error {

	from := rt.EnvString("IDX_SENDER_MAIL_ID", "idx-noreply@example.com")

	content, err := str.SimpleTemplateExpand(&str.TemplateDesc{
		Template: template,
		Data:     data,
		Html:     true,
	})
	if err != nil {
		return err
	}

	err = MailProvider(gtx).Send(&email.Message{
		From: from,
		To: []string{
			to,
		},
		Content: content,
	}, true)

	return err

}
