-- Схема gateway. Застосовується автоматично на старті (store.Open), цей файл —
-- довідковий/для ручних міграцій.
CREATE TABLE IF NOT EXISTS chats (
    chat_id          BIGINT PRIMARY KEY,
    lat              DOUBLE PRECISION NOT NULL DEFAULT 0,
    lon              DOUBLE PRECISION NOT NULL DEFAULT 0,
    place_name       TEXT NOT NULL DEFAULT '',
    tz               TEXT NOT NULL DEFAULT 'Europe/Kyiv',
    interval_minutes INT NOT NULL DEFAULT 90,
    enabled          BOOLEAN NOT NULL DEFAULT TRUE,
    last_sent        TIMESTAMPTZ NOT NULL DEFAULT to_timestamp(0)
);
