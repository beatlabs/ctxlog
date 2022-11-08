package ctxlog

import (
	"context"
	"fmt"
	"net/http"

	"github.com/beatlabs/patron/log"
	"github.com/google/uuid"
)

// List of header values that should be forwarded during login
const (
	AmazonTraceHeader = "X-Amzn-Trace-Id"
	ForwardForHeader  = "X-Forwarded-For"
	RequestIDHeader   = "X-REQUEST-ID"
	UserAgentHeader   = "User-Agent"
)

// Constants for context data.
const (
	RequestID      = "request_id"
	AmazonTraceID  = "amazon_trace_id"
	IPForwardedFor = "ip_forwarded_for"
	UserAgent      = "user_agent"
)

// CtxKey json key on which the context data is logged.
const CtxKey = "ctx"

// CtxLogger is a logger implementation that delegates logging fully to Patron and always adds some logging context to a message.
type CtxLogger struct {
	logger log.Logger

	contextData map[string]interface{}
}

// AddLoggerToContext adds a ctx logger to the given context.
func AddLoggerToContext(ctx context.Context) context.Context {
	return log.WithContext(ctx, FromContext(ctx))
}

// AddLoggerForRequest adds a ctx logger to the context of an HTTP request
func AddLoggerForRequest(req *http.Request) context.Context {
	ctx := AddLoggerToContext(req.Context())

	return log.WithContext(ctx, FromContext(ctx).SubCtx(map[string]interface{}{
		RequestID:      getRequestID(req),
		AmazonTraceID:  req.Header.Get(AmazonTraceHeader),
		IPForwardedFor: req.Header.Get(ForwardForHeader),
		UserAgent:      req.Header.Get(UserAgentHeader),
	}))
}

func getRequestID(req *http.Request) string {
	reqIDFromHeader := req.Header.Get(RequestIDHeader)
	if reqIDFromHeader == "" {
		return uuid.New().String()
	}
	return reqIDFromHeader
}

// FromContext associates logger with a context and returns both the context and the logger.
// If logger in this context has been already created, returns that one, otherwise creates a new one and associates it with a context.
// When created, a logger gets a unique request UUID.
func FromContext(ctx context.Context) *CtxLogger {
	logger := log.FromContext(ctx)
	if ctxLogger, ok := logger.(*CtxLogger); ok {
		return ctxLogger
	}

	ctxLogger := &CtxLogger{
		logger: logger,
		contextData: map[string]interface{}{
			RequestID: uuid.New().String(),
		},
	}

	return ctxLogger
}

// RequestID returns a current request ID.
func (l *CtxLogger) RequestID() string {
	reqID, ok := l.contextData[RequestID].(string)
	if !ok {
		requestID := uuid.New().String()
		l.contextData[RequestID] = requestID
		return requestID
	}
	return reqID
}

// Str returns a logger with a key-value added to log context.
func (l *CtxLogger) Str(key, value string) *CtxLogger {
	l.contextData[key] = value
	return l
}

// Int returns a logger with a key-value added to log context.
func (l *CtxLogger) Int(key string, value int) *CtxLogger {
	l.contextData[key] = fmt.Sprintf("%d", value)
	return l
}

// SubCtx appends a list of fields to the context data.
func (l *CtxLogger) SubCtx(fields map[string]interface{}) *CtxLogger {
	for k, v := range fields {
		l.contextData[k] = v
	}
	return l
}

// Sub returns a sub logger with new fields attached.
func (l *CtxLogger) Sub(fields map[string]interface{}) log.Logger {
	return &CtxLogger{
		logger:      l.logger.Sub(fields),
		contextData: l.contextData,
	}
}

// Level returns the current loglevel.
func (l *CtxLogger) Level() log.Level {
	return l.logger.Level()
}

// Fatal adds log context data to the log message.
func (l *CtxLogger) Fatal(i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Fatal(fmt.Sprint(i...))
}

// Fatalf adds log context data to the log message.
func (l *CtxLogger) Fatalf(s string, i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Fatal(fmt.Sprintf(s, i...))
}

// Panic adds log context data to the log message.
func (l *CtxLogger) Panic(i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Panic(fmt.Sprint(i...))
}

// Panicf adds log context data to the log message.
func (l *CtxLogger) Panicf(s string, i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Panic(fmt.Sprintf(s, i...))
}

// Error adds log context data to the log message.
func (l *CtxLogger) Error(i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Error(fmt.Sprint(i...))
}

// Errorf adds log context data to the log message.
func (l *CtxLogger) Errorf(s string, i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Error(fmt.Sprintf(s, i...))
}

// Warn adds log context data to the log message.
func (l *CtxLogger) Warn(i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Warn(fmt.Sprint(i...))
}

// Warnf adds log context data to the log message.
func (l *CtxLogger) Warnf(s string, i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Warn(fmt.Sprintf(s, i...))
}

// Info adds log context data to the log message.
func (l *CtxLogger) Info(i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Info(fmt.Sprint(i...))
}

// Infof adds log context data to the log message.
func (l *CtxLogger) Infof(s string, i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Info(fmt.Sprintf(s, i...))
}

// Debug adds log context data to the log message.
func (l *CtxLogger) Debug(i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Debug(fmt.Sprint(i...))
}

// Debugf adds log context data to the log message.
func (l *CtxLogger) Debugf(s string, i ...interface{}) {
	l.logger.Sub(map[string]interface{}{CtxKey: l.contextData}).Debug(fmt.Sprintf(s, i...))
}

// Logf logs at the provided level, adding log context data to the logs message.
func (l *CtxLogger) Logf(loglevel log.Level, s string, i ...interface{}) {
	switch loglevel {
	case log.DebugLevel:
		l.Debugf(s, i...)
	case log.InfoLevel:
		l.Infof(s, i...)
	case log.WarnLevel:
		l.Warnf(s, i...)
	case log.ErrorLevel:
		l.Errorf(s, i...)
	case log.FatalLevel:
		l.Fatalf(s, i...)
	case log.PanicLevel:
		l.Panicf(s, i...)
	}
}
