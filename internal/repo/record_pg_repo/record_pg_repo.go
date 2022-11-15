package recordpgrepo

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

type RecordPGRepo struct {
	db *sqlx.DB
}

func New(dsn string) (*RecordPGRepo, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &RecordPGRepo{db: db}, nil
}

func (repo *RecordPGRepo) Save(ctx context.Context, rec *domain.Record) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx, "record_pg_repo: save")
	defer span.Finish()

	tx, err := repo.db.Begin()
	if err != nil {
		return -1, err
	}
	sqlQuery := `
	INSERT INTO records (user_id, amount, category, date)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	var id int64 = 0
	rows, err := tx.Query(sqlQuery, rec.UserId, rec.Amount, rec.Category, rec.Date)
	if err != nil {
		return -1, err
	}
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			logger.Fatal(err)
			return -1, err
		} else {
			logger.Debug("saved id = ", id)
			break
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Error("Postgres transaction error", err)
	} else {
		logger.Debug("Postgres transaction ok")
		rec.Id = id
	}
	return rec.Id, err
}

func (repo *RecordPGRepo) Find(ctx context.Context, id int64) (*domain.Record, error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx, "record_pg_repo: find")
	defer span.Finish()

	rec := domain.Record{}
	sqlQuery := `
	SELECT * FROM person WHERE id = ($1) LIMIT 1
	`
	err := repo.db.Select(&rec, sqlQuery, id)
	return &rec, err
}

func (repo *RecordPGRepo) ThisWeek(ctx context.Context, userId int64) ([]domain.Record, error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx, "record_pg_repo: this_week")
	defer span.Finish()

	records := []domain.Record{}
	sqlQuery := `
	select * from records
	where user_id = ($1) and
	date_trunc('week', records.date) = date_trunc('week', current_date)
	`
	err := repo.db.Select(&records, sqlQuery, userId)
	if err != nil {
		logger.Error("Select error", err)
	}
	return records, err
}

func (repo *RecordPGRepo) ThisMonth(ctx context.Context, userId int64) ([]domain.Record, error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx, "record_pg_repo: this_month")
	defer span.Finish()

	records := []domain.Record{}
	sqlQuery := `
	select * from records
	where user_id = ($1) and
	date_trunc('month', records.date) = date_trunc('month', current_date)
	`
	err := repo.db.Select(&records, sqlQuery, userId)
	if err != nil {
		logger.Error("Select error", err)
	}
	return records, err
}

func (repo *RecordPGRepo) ThisYear(ctx context.Context, userId int64) ([]domain.Record, error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx, "record_pg_repo: this_year")
	defer span.Finish()

	records := []domain.Record{}
	sqlQuery := `
	select * from records
	where user_id = ($1) and
	date_trunc('year', records.date) = date_trunc('year', current_date)
	`
	err := repo.db.Select(&records, sqlQuery)
	if err != nil {
		logger.Error("Select error", err)
	}
	return records, err
}

func (repo *RecordPGRepo) ThisMonthSum(ctx context.Context, userId int64) (float64, error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx, "record_pg_repo: this_month_sum")
	defer span.Finish()

	sqlQuery := `
	select sum(amount) from records
	where user_id = ($1) and
	date_trunc('month', records.date) = date_trunc('month', current_date)
	`
	var total float64
	err := repo.db.Get(&total, sqlQuery, userId)
	if err != nil {
		logger.Error("ThisMonthSum() error", err)
	}

	return total, nil
}

func (repo *RecordPGRepo) Count(ctx context.Context, userId int64) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(
		ctx, "record_pg_repo: count")
	defer span.Finish()

	var count int64

	sqlQuery := `
	select count(*) from records
	where user_id = ($1)`

	err := repo.db.Get(&count, sqlQuery, userId)
	if err != nil {
		logger.Error("Select error", err)
	}
	return count, err
}
