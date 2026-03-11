package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the global structured logger.
// In production it writes JSON; in development it writes coloured text.
func Init(env string) {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "ts"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	log, err = cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		// Fall back to no-op if build fails (should never happen)
		log = zap.NewNop()
	}
}

// Sync flushes buffered log entries. Call this on application exit.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// L returns the global sugared logger. Panics if Init was not called.
func L() *zap.SugaredLogger {
	if log == nil {
		// Fallback: create a basic logger so callers don't panic
		log, _ = zap.NewDevelopment()
	}
	return log.Sugar()
}

// Named returns a child logger tagged with a component name (e.g. "matchmaking")
func Named(name string) *zap.SugaredLogger {
	return L().Named(name)
}

// Fatal logs a fatal message and exits the process.
func Fatal(msg string, args ...interface{}) {
	L().Fatalf(msg, args...)
	os.Exit(1)
}

// Info logs at INFO level with key-value pairs.
func Info(msg string, keysAndValues ...interface{}) {
	L().Infow(msg, keysAndValues...)
}

// Warn logs at WARN level.
func Warn(msg string, keysAndValues ...interface{}) {
	L().Warnw(msg, keysAndValues...)
}

// Error logs at ERROR level.
func Error(msg string, keysAndValues ...interface{}) {
	L().Errorw(msg, keysAndValues...)
}

// Debug logs at DEBUG level.
func Debug(msg string, keysAndValues ...interface{}) {
	L().Debugw(msg, keysAndValues...)
}
