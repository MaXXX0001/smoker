// Package store — persistence контексту gateway поверх PostgreSQL (pgx).
package store

import (
	"context"
	"errors"
	"time"

	"smoker/internal/gateway/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound — чату немає в БД.
var ErrNotFound = errors.New("chat not found")

const schema = `
CREATE TABLE IF NOT EXISTS chats (
    chat_id          BIGINT PRIMARY KEY,
    lat              DOUBLE PRECISION NOT NULL DEFAULT 0,
    lon              DOUBLE PRECISION NOT NULL DEFAULT 0,
    place_name       TEXT NOT NULL DEFAULT '',
    tz               TEXT NOT NULL DEFAULT 'Europe/Kyiv',
    interval_minutes INT NOT NULL DEFAULT 90,
    enabled          BOOLEAN NOT NULL DEFAULT TRUE,
    last_sent        TIMESTAMPTZ NOT NULL DEFAULT to_timestamp(0)
);`

// Store — пул з'єднань + операції над чатами.
type Store struct {
	pool *pgxpool.Pool
}

// Open відкриває пул і застосовує схему (idempotent).
func Open(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if _, err := pool.Exec(ctx, schema); err != nil {
		pool.Close()
		return nil, err
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() { s.pool.Close() }

// Ensure створює чат із дефолтами, якщо його ще нема.
func (s *Store) Ensure(ctx context.Context, chatID int64, defLat, defLon float64, defPlace, defTZ string, defInterval int) error {
	_, err := s.pool.Exec(ctx, `
        INSERT INTO chats (chat_id, lat, lon, place_name, tz, interval_minutes)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (chat_id) DO NOTHING`,
		chatID, defLat, defLon, defPlace, defTZ, defInterval)
	return err
}

// SetLocation оновлює координати/назву/TZ.
func (s *Store) SetLocation(ctx context.Context, chatID int64, lat, lon float64, place, tz string) error {
	_, err := s.pool.Exec(ctx, `
        UPDATE chats SET lat=$2, lon=$3, place_name=$4, tz=$5 WHERE chat_id=$1`,
		chatID, lat, lon, place, tz)
	return err
}

// SetInterval оновлює інтервал нагадувань (хвилини).
func (s *Store) SetInterval(ctx context.Context, chatID int64, minutes int) error {
	_, err := s.pool.Exec(ctx, `UPDATE chats SET interval_minutes=$2 WHERE chat_id=$1`, chatID, minutes)
	return err
}

// SetEnabled вмикає/вимикає чат.
func (s *Store) SetEnabled(ctx context.Context, chatID int64, enabled bool) error {
	_, err := s.pool.Exec(ctx, `UPDATE chats SET enabled=$2 WHERE chat_id=$1`, chatID, enabled)
	return err
}

// MarkSent фіксує час останнього надісланого нагадування.
func (s *Store) MarkSent(ctx context.Context, chatID int64, at time.Time) error {
	_, err := s.pool.Exec(ctx, `UPDATE chats SET last_sent=$2 WHERE chat_id=$1`, chatID, at)
	return err
}

// Get повертає чат.
func (s *Store) Get(ctx context.Context, chatID int64) (domain.Chat, error) {
	row := s.pool.QueryRow(ctx, `
        SELECT chat_id, lat, lon, place_name, tz, interval_minutes, enabled, last_sent
        FROM chats WHERE chat_id=$1`, chatID)
	c, err := scanChat(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Chat{}, ErrNotFound
	}
	return c, err
}

// ListEnabled повертає всі активні чати (для планувальника).
func (s *Store) ListEnabled(ctx context.Context) ([]domain.Chat, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT chat_id, lat, lon, place_name, tz, interval_minutes, enabled, last_sent
        FROM chats WHERE enabled=TRUE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Chat
	for rows.Next() {
		c, err := scanChat(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanChat(r scanner) (domain.Chat, error) {
	var c domain.Chat
	err := r.Scan(&c.ChatID, &c.Lat, &c.Lon, &c.PlaceName, &c.TZ,
		&c.IntervalMinutes, &c.Enabled, &c.LastSent)
	return c, err
}
