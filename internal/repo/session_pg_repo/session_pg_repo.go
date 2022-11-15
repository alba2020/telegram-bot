package sessionpgrepo

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

type SessionPGRepo struct {
	db *sqlx.DB
}

func New(dsn string) (*SessionPGRepo, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &SessionPGRepo{db: db}, nil
}

func (repo *SessionPGRepo) Find(ctx context.Context, userId int64) (*domain.Session, error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx, "session_pg_repo: find")
	defer span.Finish()

	s := make([]domain.Session, 1)
	sqlQuery := `SELECT * FROM sessions WHERE user_id = ($1) LIMIT 1`

	err := repo.db.Select(&s, sqlQuery, userId)
	logger.Debug("find()", s)
	if err != nil {
		logger.Error("Session find error", err)
		return nil, err
	} else if len(s) != 1 {
		return nil, errors.New("session not found")
	} else {
		return &s[0], nil
	}
}

func (repo *SessionPGRepo) GetOrCreate(
	ctx context.Context, userId int64) (*domain.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx, "session_pg_repo: get_or_create")
	defer span.Finish()

	logger.Debug("get or create session")
	existingSession, err := repo.Find(ctx, userId)
	if err == nil {
		logger.Debug("session found")
		return existingSession, nil
	}
	logger.Debug("session not found")

	sqlQuery := `INSERT INTO sessions (user_id, currency, month_limit)
					VALUES ($1, $2, $3)`

	newSession, _ := domain.NewSession(userId)

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(sqlQuery, newSession.UserId, newSession.Currency, newSession.MonthLimit)
	if err != nil {
		logger.Error("Error saving session ", err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (repo *SessionPGRepo) Save(ctx context.Context, s *domain.Session) error {
	span, _ := opentracing.StartSpanFromContext(
		ctx, "session_pg_repo: save")
	defer span.Finish()

	sqlQuery := `
	UPDATE sessions
	SET currency=($1), month_limit=($2)
	WHERE user_id=($3)`

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sqlQuery, s.Currency, s.MonthLimit, s.UserId)
	if err != nil {
		logger.Error("Error saving session", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
