// Package iologger provides slog-based logging initialization and configuration.
package iologger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gnames/gnfinder/pkg/config"
	"github.com/lmittmann/tint"
)

// LogFile is the name of the log file created in the log directory.
const LogFile = "gnfinder.log"

// logFile is the currently open log file, if any.
var logFile *os.File

// Init initializes the global slog logger with the given configuration.
// Creates log file in logDir if destination is "file".
// If append is true, appends to existing log file; otherwise creates fresh file.
func Init(logDir string, cfg config.LogConfig, append bool) error {
	var writer io.Writer
	var file *os.File

	// Determine output destination
	switch cfg.Destination {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	case "file":
		logPath := filepath.Join(logDir, LogFile)
		var err error

		if append {
			// Append to existing log file (preserve previous logs)
			file, err = os.OpenFile(
				logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644,
			)
		} else {
			// Create fresh log file (truncate if exists)
			file, err = os.Create(logPath)
		}

		if err != nil {
			return CreateLogFileError(logPath, err)
		}
		writer = file
	default:
		writer = os.Stderr
	}

	level := parseLevel(cfg.Level)

	var handler slog.Handler
	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	switch cfg.Format {
	case "text":
		handler = slog.NewTextHandler(writer, handlerOpts)
	case "tint":
		handler = tint.NewTextHandler(writer, &tint.Options{
			Level: level,
			// colors make no sense in a log file
			NoColor: cfg.Destination == "file",
		})
	default:
		// JSON format for "json" and any unrecognized format
		handler = slog.NewJSONHandler(writer, handlerOpts)
	}

	slog.SetDefault(slog.New(handler))

	// close the previous log file, it is not used by the new logger
	if logFile != nil {
		_ = logFile.Close()
	}
	logFile = file

	return nil
}

// Close closes the log file opened by Init, if any.
func Close() error {
	if logFile == nil {
		return nil
	}
	err := logFile.Close()
	logFile = nil
	return err
}

// parseLevel converts string level to slog.Level.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
