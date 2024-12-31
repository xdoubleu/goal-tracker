package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

const ProgressDateFormat = "2006-01-02"

type Progress struct {
	TypeID int64            `json:"typeId"`
	Date   pgtype.Timestamp `json:"date"`
	Value  int64            `json:"value"`
} //	@name	Progress
