package log

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// RequestIDKey is the key used to store and retrieve the request ID from context
const RequestIDKey contextKey = "requestID"

// UserEmailKey is the key used to store and retrieve the user email from context
const UserEmailKey contextKey = "userEmail"

// UserNameKey is the key used to store and retrieve the user name from context
const UserNameKey contextKey = "userName"

var logger = logrus.New()

// LogLevel represents the severity level of the log
type LogLevel string

// Log levels
const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
	PanicLevel LogLevel = "panic"
	TraceLevel LogLevel = "trace"
)

// InitLogger initializes the logger with the specified configuration
func InitLogger() {
	// Set the formatter to JSON
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Set the output to both file and stdout
	logger.SetOutput(io.MultiWriter(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100, // megabytes
		MaxAge:     5,   // days
		MaxBackups: 2,
		LocalTime:  false,
		Compress:   true,
	}, os.Stdout))

	// Set the log level (default to Info)
	SetLogLevel(InfoLevel)

	// Enable caller reporting by default
	logger.SetReportCaller(true)

	Info("Logger initialized successfully")
}

// SetLogLevel sets the log level based on the provided LogLevel
func SetLogLevel(level LogLevel) {
	switch level {
	case DebugLevel:
		logger.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		logger.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		logger.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		logger.SetLevel(logrus.ErrorLevel)
	case FatalLevel:
		logger.SetLevel(logrus.FatalLevel)
	case PanicLevel:
		logger.SetLevel(logrus.PanicLevel)
	case TraceLevel:
		logger.SetLevel(logrus.TraceLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

// Debug logs a message at the debug level
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf logs a formatted message at the debug level
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info logs a message at the info level
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof logs a formatted message at the info level
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn logs a message at the warn level
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf logs a formatted message at the warn level
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Error logs a message at the error level
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf logs a formatted message at the error level
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal logs a message at the fatal level and then calls os.Exit(1)
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf logs a formatted message at the fatal level and then calls os.Exit(1)
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Panic logs a message at the panic level and then panics
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Panicf logs a formatted message at the panic level and then panics
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// Trace logs a message at the trace level
func Trace(args ...interface{}) {
	logger.Trace(args...)
}

// Tracef logs a formatted message at the trace level
func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

// WithFields creates an entry from the standard logger and adds multiple fields to it
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

// WithField creates an entry from the standard logger and adds a field to it
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

// WithContext creates an entry from the standard logger and adds context fields to it
func WithContext(ctx context.Context) *logrus.Entry {
	entry := logger.WithFields(logrus.Fields{})

	// Add request ID if available
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		entry = entry.WithField("request_id", requestID)
	}

	// Add user email if available
	if userEmail, ok := ctx.Value(UserEmailKey).(string); ok {
		entry = entry.WithField("user_email", userEmail)
	}

	// Add user name if available
	if userName, ok := ctx.Value(UserNameKey).(string); ok {
		entry = entry.WithField("user_name", userName)
	}

	return entry
}

// SetReportCaller enables or disables including the calling method as a field in the log
func SetReportCaller(caller bool) {
	logger.SetReportCaller(caller)
}

// NewRequestID generates a new request ID
func NewRequestID() string {
	// Generate 8 random bytes
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Fallback to timestamp if random generation fails
		return fmt.Sprintf("req-%d", time.Now().UnixNano())
	}

	// Combine timestamp and random bytes for uniqueness
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 16)
	randomHex := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("req-%s-%s", timestamp, randomHex)
}

// ContextWithRequestID adds a request ID to the context
func ContextWithRequestID(ctx context.Context) context.Context {
	requestID := NewRequestID()
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// ContextWithUserInfo adds user information to the context
func ContextWithUserInfo(ctx context.Context, userEmail, userName string) context.Context {
	ctx = context.WithValue(ctx, UserEmailKey, userEmail)
	return context.WithValue(ctx, UserNameKey, userName)
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetCallerInfo returns the file name and line number of the caller
func GetCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}

	// Get just the file name, not the full path
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	return fmt.Sprintf("%s:%d", file, line)
}
