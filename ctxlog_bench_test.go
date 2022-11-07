package ctxlog

import (
	"context"
	"testing"
)

// BenchmarkFromContextWithLoggerGetEveryTime benches when we retrieve logger from context every time
// 14850 ns/op
func BenchmarkFromContextWithLoggerGetEveryTime(b *testing.B) {
	ctx := AddLoggerForRequest(testRequest())
	for i := 0; i < b.N; i++ {
		for i := 0; i < 20; i++ {
			FromContext(ctx).Info("hello")
		}
	}
}

// BenchmarkFromContextWithLoggerGetOnlyOnce benches when we retrieve logger only once and keep using it
// 17059 ns/op
func BenchmarkFromContextWithLoggerGetOnlyOnce(b *testing.B) {
	ctx := AddLoggerForRequest(testRequest())
	for i := 0; i < b.N; i++ {
		logger := FromContext(ctx)
		for i := 0; i < 20; i++ {
			logger.Info("hello")
		}
	}
}

// BenchmarkFromContextWithoutLogger benches what happens if logger is not initialized
// 22060 ns/op
func BenchmarkFromContextWithoutLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < 20; i++ {
			FromContext(context.Background()).Info("hello")
		}
	}
}
