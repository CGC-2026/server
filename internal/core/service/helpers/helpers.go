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

func NewNullFloat(v *float64) pgtype.Float8 {
	if v == nil {
		return pgtype.Float8{Valid: false}
	}
	return pgtype.Float8{Float64: *v, Valid: true}
}
