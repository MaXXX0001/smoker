// Package env — мінімальні хелпери читання конфігу з оточення, щоб не тягнути
// зайвих залежностей. Кожен сервіс читає свої *_GRPC_ADDR тощо.
package env

import (
	"log/slog"
	"os"
	"strconv"
)

// String повертає значення змінної або def.
func String(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

// Float повертає float-значення або def.
func Float(key string, def float64) float64 {
	if v, ok := os.LookupEnv(key); ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

// Logger будує JSON slog-логер із рівнем з LOG_LEVEL (debug/info/warn/error).
func Logger(service string) *slog.Logger {
	lvl := slog.LevelInfo
	switch String("LOG_LEVEL", "info") {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return slog.New(h).With("service", service)
}
