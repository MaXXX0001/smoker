-- Схема gateway (SQLite). Застосовується автоматично на старті (store.Open),
-- цей файл — довідковий. enabled: 0/1, last_sent: Unix-секунди.
CREATE TABLE IF NOT EXISTS chats (
    chat_id          INTEGER PRIMARY KEY,
    lat              REAL    NOT NULL DEFAULT 0,
    lon              REAL    NOT NULL DEFAULT 0,
    place_name       TEXT    NOT NULL DEFAULT '',
    tz               TEXT    NOT NULL DEFAULT 'Europe/Kyiv',
    interval_minutes INTEGER NOT NULL DEFAULT 90,
    enabled          INTEGER NOT NULL DEFAULT 1,
    last_sent        INTEGER NOT NULL DEFAULT 0
);
