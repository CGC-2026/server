package helpers

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func NewNullString(s string) pgtype.Text {
	return pgtype.Text{
		Valid:  s != "",
		String: s,
	}
}

func NewNullTime(t *time.Time) pgtype.Timestamptz {
	if t == nil || t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{
		Time:  *t,
		Valid: true,
	}
}
