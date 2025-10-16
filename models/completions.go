package models

import (
	"database/sql"
)

type ChatCompletionService struct {
	DB *sql.DB
}
