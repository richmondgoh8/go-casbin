package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapSugarLog *zap.SugaredLogger
var trackingFieldKeys []string
var ZapBasicLogger *zap.Logger

func init() {
	Init([]string{})
}

const (
	ENV_LOGGING_LEVEL = "LOGGING_LEVEL"
	ENV_SERVICE_NAME  = "SERVICE_NAME"
)

// Init overrides the base init if required for OS Environment & Exclusion
func Init(trackingFieldKeysReq []string) {
	trackingFieldKeys = trackingFieldKeysReq
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.StacktraceKey = ""
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig

	var globalLoggingLevel zapcore.Level

	switch os.Getenv(ENV_LOGGING_LEVEL) {
	case "DEBUG":
		globalLoggingLevel = zapcore.DebugLevel
	case "INFO":
		globalLoggingLevel = zapcore.InfoLevel
	case "WARN":
		globalLoggingLevel = zapcore.WarnLevel
	case "ERROR":
		globalLoggingLevel = zapcore.ErrorLevel
	default:
		globalLoggingLevel = zapcore.DebugLevel
	}

	config.Level = zap.NewAtomicLevelAt(globalLoggingLevel)

	value, exists := os.LookupEnv(ENV_SERVICE_NAME)
	if exists {
		config.InitialFields = map[string]interface{}{
			"engine": value,
		}
	} else {
		config.InitialFields = map[string]interface{}{
			"engine": "not_specified",
		}
	}

	// Caller Skip allows the proper caller instead of showing pkg caller
	ZapLogger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	// clears buffers if any
	defer ZapLogger.Sync()
	zapSugarLog = ZapLogger.Sugar()
}

func InitBasic() {
	var err error
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.StacktraceKey = ""
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig

	var globalLoggingLevel zapcore.Level

	switch os.Getenv(ENV_LOGGING_LEVEL) {
	case "DEBUG":
		globalLoggingLevel = zapcore.DebugLevel
	case "INFO":
		globalLoggingLevel = zapcore.InfoLevel
	case "WARN":
		globalLoggingLevel = zapcore.WarnLevel
	case "ERROR":
		globalLoggingLevel = zapcore.ErrorLevel
	default:
		globalLoggingLevel = zapcore.DebugLevel
	}

	config.Level = zap.NewAtomicLevelAt(globalLoggingLevel)

	value, exists := os.LookupEnv(ENV_SERVICE_NAME)
	if exists {
		config.InitialFields = map[string]interface{}{
			"engine": value,
		}
	} else {
		config.InitialFields = map[string]interface{}{
			"engine": "not_specified",
		}
	}

	// Caller Skip allows the proper caller instead of showing pkg caller
	ZapBasicLogger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	// clears buffers if any
	defer ZapBasicLogger.Sync()
	zapSugarLog = ZapBasicLogger.Sugar()
}

func Debug(msg string, ctx context.Context, fields map[string]interface{}) {
	var keyValueField []interface{}
	for k, v := range fields {
		keyValueField = append(keyValueField, k, v)
	}

	keyValueField = getTrackingDetails(ctx, keyValueField)

	zapSugarLog.Debugw(msg, keyValueField...)
}

func Info(msg string, ctx context.Context, fields map[string]interface{}) {
	var keyValueField []interface{}
	for k, v := range fields {
		keyValueField = append(keyValueField, k, v)
	}

	keyValueField = getTrackingDetails(ctx, keyValueField)

	zapSugarLog.Infow(msg, keyValueField...)
}

func Warn(msg string, ctx context.Context, fields map[string]interface{}) {
	var keyValueField []interface{}
	for k, v := range fields {
		keyValueField = append(keyValueField, k, v)
	}

	keyValueField = getTrackingDetails(ctx, keyValueField)

	zapSugarLog.Warnw(msg, keyValueField...)
}

func Error(msg string, ctx context.Context, fields map[string]interface{}) {
	var keyValueField []interface{}
	for k, v := range fields {
		keyValueField = append(keyValueField, k, v)
	}

	keyValueField = getTrackingDetails(ctx, keyValueField)

	zapSugarLog.Errorw(msg, keyValueField...)
}

func getTrackingDetails(ctx context.Context, keyValueField []interface{}) []interface{} {
	for _, key := range trackingFieldKeys {
		if val, exist := ctx.Value(key).(interface{}); exist {
			keyValueField = append(keyValueField, key, val)
		}
	}

	return keyValueField
}
