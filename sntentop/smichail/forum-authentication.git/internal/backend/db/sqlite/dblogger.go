package sqlite

import (
	"context"
	"database/sql"
	"forum-authentication/internal/utils"
	"os"
)

type DBlogger struct {
	DB      *sql.DB
	logfile *os.File
}

func (l *DBlogger) LogExecContext(ctx context.Context, q string, args ...any) (sql.Result, error) {

	qtolog := utils.FormatLoggingQuery(q, args...)
	l.logfile.WriteString(qtolog)
	return l.DB.ExecContext(ctx, q, args...)
}
