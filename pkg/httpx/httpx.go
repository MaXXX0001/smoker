// Package httpx — тонкий JSON-HTTP клієнт зі строгими таймаутами та одним
// повтором. Використовується ACL-шарами provider-сервісів для походів у
// зовнішні API. Жодне з цих API не критичне: якщо виклик впав — умову просто
// пропускають, тому помилки тут повертаються, а не паніку.
package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client — обгортка над http.Client з дефолтним таймаутом і заголовками.
type Client struct {
	hc      *http.Client
	headers map[string]string
}

// Option конфігурує Client.
type Option func(*Client)

// WithTimeout задає таймаут на запит.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.hc.Timeout = d }
}

// WithHeader додає заголовок до кожного запиту (напр. Accept, User-Agent).
func WithHeader(k, v string) Option {
	return func(c *Client) { c.headers[k] = v }
}

// New створює клієнт. За замовчуванням таймаут 8с і ввічливий User-Agent —
// деякі API (icanhazdadjoke, Wikimedia) вимагають його.
func New(opts ...Option) *Client {
	c := &Client{
		hc:      &http.Client{Timeout: 8 * time.Second},
		headers: map[string]string{"User-Agent": "smoker-bot/0.1 (+https://github.com/smoker)"},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// GetJSON виконує GET і декодує тіло у dst. Один повтор на мережевих/5xx
// помилках з невеликою паузою.
func (c *Client) GetJSON(ctx context.Context, url string, dst any) error {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(300 * time.Millisecond):
			}
		}
		body, err := c.getRaw(ctx, url)
		if err != nil {
			lastErr = err
			continue
		}
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("decode %s: %w", url, err)
		}
		return nil
	}
	return lastErr
}

// GetBytes повертає сире тіло (для API, що віддають не-JSON або текст).
func (c *Client) GetBytes(ctx context.Context, url string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(300 * time.Millisecond):
			}
		}
		body, err := c.getRaw(ctx, url)
		if err != nil {
			lastErr = err
			continue
		}
		return body, nil
	}
	return nil, lastErr
}

func (c *Client) getRaw(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // максимум 1 МіБ
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}
	return body, nil
}
