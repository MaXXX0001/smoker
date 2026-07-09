// Package store — persistence контексту gateway поверх SQLite (чистий Go-драйвер
// modernc.org/sqlite, без CGO). Увесь стан бота — це кілька рядків у одному файлі.
package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"smoker/internal/gateway/domain"

	_ "modernc.org/sqlite"
)

// ErrNotFound — чату немає в БД.
var ErrNotFound = errors.New("chat not found")

const schema = `
CREATE TABLE IF NOT EXISTS chats (
    chat_id          INTEGER PRIMARY KEY,
    lat              REAL    NOT NULL DEFAULT 0,
    lon              REAL    NOT NULL DEFAULT 0,
    place_name       TEXT    NOT NULL DEFAULT '',
    tz               TEXT    NOT NULL DEFAULT 'Europe/Kyiv',
    interval_minutes INTEGER NOT NULL DEFAULT 90,
    enabled          INTEGER NOT NULL DEFAULT 1,
    last_sent        INTEGER NOT NULL DEFAULT 0
);`

// Store — одне з'єднання до SQLite-файлу + операції над чатами.
type Store struct {
	db *sql.DB
}

// Open відкриває файл БД (створює за потреби) і застосовує схему (idempotent).
// path — шлях до файлу, напр. "/data/smoker.db".
func Open(ctx context.Context, path string) (*Store, error) {
	// WAL + busy_timeout — щоб паралельні читання планувальника й запис хендлера
	// не билися за блокування; foreign_keys про запас.
	dsn := path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// SQLite — один writer; тримаємо одне з'єднання, щоб уникнути "database is locked".
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() { s.db.Close() }

// Ensure створює чат із дефолтами, якщо його ще нема.
func (s *Store) Ensure(ctx context.Context, chatID int64, defLat, defLon float64, defPlace, defTZ string, defInterval int) error {
	_, err := s.db.ExecContext(ctx, `
        INSERT OR IGNORE INTO chats (chat_id, lat, lon, place_name, tz, interval_minutes)
        VALUES (?, ?, ?, ?, ?, ?)`,
		chatID, defLat, defLon, defPlace, defTZ, defInterval)
	return err
}

// SetLocation оновлює координати/назву/TZ.
func (s *Store) SetLocation(ctx context.Context, chatID int64, lat, lon float64, place, tz string) error {
	_, err := s.db.ExecContext(ctx, `
        UPDATE chats SET lat=?, lon=?, place_name=?, tz=? WHERE chat_id=?`,
		lat, lon, place, tz, chatID)
	return err
}

// SetInterval оновлює інтервал нагадувань (хвилини).
func (s *Store) SetInterval(ctx context.Context, chatID int64, minutes int) error {
	_, err := s.db.ExecContext(ctx, `UPDATE chats SET interval_minutes=? WHERE chat_id=?`, minutes, chatID)
	return err
}

// SetEnabled вмикає/вимикає чат.
func (s *Store) SetEnabled(ctx context.Context, chatID int64, enabled bool) error {
	_, err := s.db.ExecContext(ctx, `UPDATE chats SET enabled=? WHERE chat_id=?`, boolToInt(enabled), chatID)
	return err
}

// MarkSent фіксує час останнього надісланого нагадування (Unix-секунди).
func (s *Store) MarkSent(ctx context.Context, chatID int64, at time.Time) error {
	_, err := s.db.ExecContext(ctx, `UPDATE chats SET last_sent=? WHERE chat_id=?`, at.Unix(), chatID)
	return err
}

// Get повертає чат.
func (s *Store) Get(ctx context.Context, chatID int64) (domain.Chat, error) {
	row := s.db.QueryRowContext(ctx, `
        SELECT chat_id, lat, lon, place_name, tz, interval_minutes, enabled, last_sent
        FROM chats WHERE chat_id=?`, chatID)
	c, err := scanChat(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Chat{}, ErrNotFound
	}
	return c, err
}

// ListEnabled повертає всі активні чати (для планувальника).
func (s *Store) ListEnabled(ctx context.Context) ([]domain.Chat, error) {
	rows, err := s.db.QueryContext(ctx, `
        SELECT chat_id, lat, lon, place_name, tz, interval_minutes, enabled, last_sent
        FROM chats WHERE enabled=1`)
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

// scanChat читає рядок; enabled/last_sent зберігаються як INTEGER, тому
// мапимо їх у bool/time.Time вручну.
func scanChat(r scanner) (domain.Chat, error) {
	var (
		c        domain.Chat
		enabled  int
		lastSent int64
	)
	if err := r.Scan(&c.ChatID, &c.Lat, &c.Lon, &c.PlaceName, &c.TZ,
		&c.IntervalMinutes, &enabled, &lastSent); err != nil {
		return domain.Chat{}, err
	}
	c.Enabled = enabled != 0
	c.LastSent = time.Unix(lastSent, 0)
	return c, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
