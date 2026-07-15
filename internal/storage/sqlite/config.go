package sqlite

import (
	"errors"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vvitem/item_all/internal/storage"
)

const (
	defaultBusyTimeout = 5 * time.Second
	defaultPoolSize    = 4
)

// Config controls the local SQLite file and bounded connection pool.
type Config struct {
	Path         string
	BusyTimeout  time.Duration
	MaxOpenConns int
	MaxIdleConns int
}

// DefaultConfig returns the approved production defaults.
func DefaultConfig(path string) Config {
	return Config{
		Path:         path,
		BusyTimeout:  defaultBusyTimeout,
		MaxOpenConns: defaultPoolSize,
		MaxIdleConns: defaultPoolSize,
	}
}

func (c Config) validate() error {
	switch {
	case c.BusyTimeout.Milliseconds() <= 0:
		return storage.Wrap(storage.CodeOpenFailed, "validate-busy-timeout", errors.New("busy timeout must be positive"), false)
	case c.MaxOpenConns <= 0:
		return storage.Wrap(storage.CodeOpenFailed, "validate-max-open-connections", errors.New("max open connections must be positive"), false)
	case c.MaxIdleConns < 0:
		return storage.Wrap(storage.CodeOpenFailed, "validate-max-idle-connections", errors.New("max idle connections must not be negative"), false)
	case c.MaxIdleConns > c.MaxOpenConns:
		return storage.Wrap(storage.CodeOpenFailed, "validate-connection-pool", errors.New("max idle connections exceeds max open connections"), false)
	default:
		return nil
	}
}

func buildDSN(config Config) (string, error) {
	if err := config.validate(); err != nil {
		return "", err
	}
	path := config.Path
	if path == "" {
		var err error
		path, err = DefaultDatabasePath()
		if err != nil {
			return "", err
		}
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", storage.Wrap(storage.CodePathUnavailable, "resolve-absolute-database-path", err, false)
	}

	normalized := filepath.ToSlash(absolute)
	if isWindowsDrivePath(absolute) {
		normalized = "/" + strings.TrimLeft(strings.ReplaceAll(absolute, `\`, "/"), "/")
	}
	uri := url.URL{Scheme: "file", Path: normalized}
	query := url.Values{}
	query.Set("mode", "rwc")
	query.Add("_pragma", "journal_mode(WAL)")
	query.Add("_pragma", "foreign_keys(1)")
	query.Add("_pragma", "busy_timeout("+strconv.FormatInt(config.BusyTimeout.Milliseconds(), 10)+")")
	query.Add("_pragma", "synchronous(NORMAL)")
	query.Set("_txlock", "immediate")
	query.Set("_dqs", "false")
	query.Set("_time_format", "sqlite")
	query.Set("_timezone", "UTC")
	uri.RawQuery = query.Encode()
	return uri.String(), nil
}

func isWindowsDrivePath(path string) bool {
	return len(path) >= 2 && path[1] == ':' && ((path[0] >= 'A' && path[0] <= 'Z') || (path[0] >= 'a' && path[0] <= 'z'))
}
