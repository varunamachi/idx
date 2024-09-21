package pg

import (
	"context"

	"github.com/lib/pq"
	"github.com/varunamachi/libx/data/event"
	"github.com/varunamachi/libx/data/pg"
	"github.com/varunamachi/libx/errx"
)

type eventSrv struct {
}

func NewEventService() event.Service[int64] {
	return &eventSrv{}
}

func (pes *eventSrv) AddEvent(
	gtx context.Context,
	event *event.Event[int64]) error {
	query := `
			INSERT INTO idx_event(
				op,
				ev_type,
				user_id,
				created_on,
				errors,
				metadata
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6
			)
		`

	_, err := pg.Conn().ExecContext(
		gtx, query,
		event.Op,
		event.Type,
		event.UserId,
		event.CreatedOn,
		pq.Array(event.Errors),
		event.Metadata)
	if err != nil {
		return errx.Errf(err, "failed to add event '%s' to database", event.Op)
	}
	return nil
}
