package ctxlog

import (
	"context"
	"net/http"
	"testing"

	mocks "github.com/beatlabs/harvester/mocks/patron"
	"github.com/beatlabs/patron/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolvesSameLogger(t *testing.T) {
	t.Run("with some context values out of patron context", func(t *testing.T) {
		// given
		ctx := AddLoggerForRequest(testRequest())

		// when
		FromContext(ctx).Str("hello", "world")
		loggerFromCtx := log.FromContext(ctx)

		// then
		castedLogger, castedOk := loggerFromCtx.(*CtxLogger)
		assert.True(t, castedOk)
		assert.Equal(t, FromContext(ctx).contextData, castedLogger.contextData)
	})

	t.Run("creates a new logger in context implicitly", func(t *testing.T) {
		// given
		ctx := AddLoggerForRequest(testRequest())

		// when
		loggerFromPatron := log.FromContext(ctx)

		// then
		assert.Same(t, FromContext(ctx), loggerFromPatron)
	})
}

func TestCtxLogger(t *testing.T) {
	tests := []struct {
		Name                    string
		ExpectedFields          []map[string]interface{}
		PrepareExpectedMessage  func(logger *mocks.Logger)
		ExpectedSubCalls        int
		LogCall                 func(context.Context)
		AdditionalContextFields map[string]interface{}
	}{
		{
			Name: "log error",
			PrepareExpectedMessage: func(loggerMock *mocks.Logger) {
				loggerMock.EXPECT().
					Error(gomock.Any()).
					DoAndReturn(func(input ...interface{}) {
						require.Equal(t, 1, len(input))
						assert.Equal(t, "error: 42", input[0])
					}).Times(1)
			},
			ExpectedSubCalls: 1,
			LogCall: func(ctx context.Context) {
				FromContext(ctx).Errorf("error: %d", 42)
			},
		},
		{
			Name: "log info with sub",
			ExpectedFields: []map[string]interface{}{
				{
					"foo": "bar",
				},
			},
			PrepareExpectedMessage: func(loggerMock *mocks.Logger) {
				loggerMock.EXPECT().
					Info(gomock.Any()).
					DoAndReturn(func(input ...interface{}) {
						require.Equal(t, 1, len(input))
						assert.Equal(t, "info: 42", input[0])
					}).Times(1)
			},
			ExpectedSubCalls: 2,
			LogCall: func(ctx context.Context) {
				FromContext(ctx).Sub(map[string]interface{}{"foo": "bar"}).Infof("info: %d", 42)
			},
		},
		{
			Name: "log debug with Str",
			PrepareExpectedMessage: func(loggerMock *mocks.Logger) {
				loggerMock.EXPECT().
					Debug(gomock.Any()).
					DoAndReturn(func(input ...interface{}) {
						require.Equal(t, 1, len(input))
						assert.Equal(t, "error: 42", input[0])
					}).Times(1)
			},
			ExpectedSubCalls: 1,
			LogCall: func(ctx context.Context) {
				FromContext(ctx).Str("hello", "world").Debugf("error: %d", 42)
			},
			AdditionalContextFields: map[string]interface{}{
				"hello": "world",
			},
		},
		{
			Name: "log fatal with Int",
			PrepareExpectedMessage: func(loggerMock *mocks.Logger) {
				loggerMock.EXPECT().
					Fatal(gomock.Any()).
					DoAndReturn(func(input ...interface{}) {
						require.Equal(t, 1, len(input))
						assert.Equal(t, "fatal: 42", input[0])
					}).Times(1)
			},
			ExpectedSubCalls: 1,
			LogCall: func(ctx context.Context) {
				FromContext(ctx).Int("my_key_int", 42).Fatalf("fatal: %d", 42)
			},
			AdditionalContextFields: map[string]interface{}{
				"my_key_int": "42",
			},
		},
		{
			Name: "log panic with both Str and Int",
			PrepareExpectedMessage: func(loggerMock *mocks.Logger) {
				loggerMock.EXPECT().
					Panic(gomock.Any()).
					DoAndReturn(func(input ...interface{}) {
						require.Equal(t, 1, len(input))
						assert.Equal(t, "paniek: 42", input[0])
					}).Times(1)
			},
			ExpectedSubCalls: 1,
			LogCall: func(ctx context.Context) {
				FromContext(ctx).Int("my_key_int", 42).Str("hello", "world").Panicf("paniek: %d", 42)
			},
			AdditionalContextFields: map[string]interface{}{
				"my_key_int": "42",
				"hello":      "world",
			},
		},
		{
			Name: "log warn with SubCtx",
			PrepareExpectedMessage: func(loggerMock *mocks.Logger) {
				loggerMock.EXPECT().
					Warn(gomock.Any()).
					DoAndReturn(func(input ...interface{}) {
						require.Equal(t, 1, len(input))
						assert.Equal(t, "warn: 42", input[0])
					}).Times(1)
			},
			ExpectedSubCalls: 1,
			LogCall: func(ctx context.Context) {
				FromContext(ctx).SubCtx(map[string]interface{}{
					"val_1": 42,
					"val_2": "hello world",
				}).Warnf("warn: %d", 42)
			},
			AdditionalContextFields: map[string]interface{}{
				"val_1": 42,
				"val_2": "hello world",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// given
			fieldCounter := 0
			ctx := AddLoggerForRequest(testRequest())
			ctxLogger := FromContext(ctx)

			loggerMock := mocks.NewLogger(gomock.NewController(t))

			ctxLogger.logger = loggerMock

			// then
			test.PrepareExpectedMessage(loggerMock)

			loggerMock.EXPECT().
				Sub(gomock.Any()).
				DoAndReturn(func(fields map[string]interface{}) log.Logger {
					if fieldCounter+1 == test.ExpectedSubCalls {
						ctxObj, ok := fields[CtxKey]
						require.True(t, ok)
						reqID := FromContext(ctx).RequestID()
						ctxMap, ok := ctxObj.(map[string]interface{})
						require.True(t, ok)
						assert.Equal(t, reqID, ctxMap[RequestID])
						if test.AdditionalContextFields != nil {
							for k, v := range test.AdditionalContextFields {
								assert.Equal(t, v, ctxMap[k])
							}
						}
						return loggerMock
					}
					assert.Equal(t, test.ExpectedFields[fieldCounter], fields)

					fieldCounter++
					return loggerMock
				}).Times(test.ExpectedSubCalls)

			// when
			test.LogCall(ctx)
		})
	}
}

func TestSubCtxChangesTheLoggerInContext(t *testing.T) {
	ctx := AddLoggerForRequest(testRequest())

	FromContext(ctx).SubCtx(map[string]interface{}{
		"foo": "bar",
	})

	assert.Equal(t, "bar", FromContext(ctx).contextData["foo"])
}

func testRequest() *http.Request {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://foo", nil)
	req.Header.Set(RequestID, "my-test-request-ID")
	return req
}
