package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Go → pgtype

func toText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func toFloat4(f float32) pgtype.Float4 {
	return pgtype.Float4{Float32: f, Valid: true}
}

func toInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func toTimestampz(t time.Time) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	ts.Scan(t)
	return ts
}

func toPgUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

// pgtype → Go

func toString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

func toFloat32(f pgtype.Float4) float32 {
	if f.Valid {
		return f.Float32
	}
	return 0
}

func toInt32(i pgtype.Int4) int32 {
	if i.Valid {
		return i.Int32
	}
	return 0
}

func toTime(ts pgtype.Timestamptz) time.Time {
	if ts.Valid {
		return ts.Time
	}
	return time.Time{}
}

func toTimePtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	return &ts.Time
}

func toUUIDPtr(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	id := uuid.UUID(u.Bytes)
	return &id
}
