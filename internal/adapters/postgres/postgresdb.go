package postgresdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-analytics/internal/utils"
	"go.uber.org/zap"
)

type PostgresDB struct {
	DB     *pgxpool.Pool
	logger *zap.SugaredLogger
}

func New(ctx context.Context, pgconn string, logger *zap.SugaredLogger) (*PostgresDB, error) {
	config, err := pgxpool.ParseConfig(pgconn)
	if err != nil {
		return nil, fmt.Errorf("postgres connection string parse failed: %s", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool failed: %s", err)
	}
	return &PostgresDB{
		DB:     pool,
		logger: logger,
	}, nil
}

func (db *PostgresDB) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return db.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}

func (db *PostgresDB) CreateEvent(ctx context.Context, event *models.Event) error {
	logger := db.annotatedLogger(ctx)

	_, err := db.DB.Exec(ctx, "INSERT INTO event VALUES ($1, $2, $3, $4, $5)",
		event.EventId,
		event.TaskId,
		event.Time,
		event.Type,
		event.Status,
	)
	if err != nil {
		logger.Errorf("failed to exec insert event query")
		return fmt.Errorf("failed to exec insert event query")
	}

	return nil
}

func (db *PostgresDB) CountDoneEvents(ctx context.Context) (int, error) {
	logger := db.annotatedLogger(ctx)

	row := db.DB.QueryRow(ctx, "SELECT count(*) as done_cnt FROM event WHERE status = $1", models.TaskDoneStatus)

	doneCnt := 0
	err := row.Scan(&doneCnt)
	if err != nil {
		logger.Errorf("failed to count done events query: %s", err.Error())
		return 0, fmt.Errorf("failed to count done events query: %s", err.Error())
	}
	return doneCnt, nil
}

func (db *PostgresDB) CountDeclinedEvents(ctx context.Context) (int, error) {
	logger := db.annotatedLogger(ctx)

	row := db.DB.QueryRow(ctx, "SELECT count(*) as done_cnt FROM event WHERE status = $1", models.TaskDeclinedStatus)

	doneCnt := 0
	err := row.Scan(&doneCnt)
	if err != nil {
		logger.Errorf("failed to exec count done events query: %s", err.Error())
		return 0, fmt.Errorf("failed to exec count done events query: %s", err.Error())
	}
	return doneCnt, nil
}

func (db *PostgresDB) GetEventsByTaskID(ctx context.Context, taskId string) ([]models.Event, error) {
	logger := db.annotatedLogger(ctx)

	rows, err := db.DB.Query(ctx, "SELECT event_id, task_id, time, type, status FROM event WHERE task_id = $1", taskId)
	if err != nil {
		logger.Errorf("failed to exec query for all events of task ID %s: %s", taskId, err.Error())
		return []models.Event{}, fmt.Errorf("failed to exec query for all events of task ID %s: %s", taskId, err.Error())
	}

	events := []models.Event{}
	for rows.Next() {
		event := models.Event{}
		err = rows.Scan(&event.EventId, &event.TaskId, &event.Time, &event.Type, &event.Status)
		if err != nil {
			logger.Errorf("failed to exec query for all events of task ID %s: %s", taskId, err.Error())
			return events, fmt.Errorf("failed to exec query for all events of task ID %s: %s", taskId, err.Error())
		}
		events = append(events, event)
	}
	if len(events) == 0 {
		logger.Errorf("no events found in DB for task ID %s", taskId)
		return []models.Event{}, fmt.Errorf("no events found in DB for task ID %s", taskId)
	}
	return events, nil
}
