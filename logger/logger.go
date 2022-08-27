package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"taleteller/constants"
)

var sugaredLogger *zap.SugaredLogger

func SetupLogger(env string) (*zap.SugaredLogger, error) {
	logger, err := getLoggerbyEnv(env)
	if err != nil {
		return nil, err
	}

	sugaredLogger = logger.Sugar()
	return sugaredLogger, nil
}

func getLoggerbyEnv(env string) (logger *zap.Logger, err error) {
	option := zap.AddCallerSkip(1)

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config.Build(option)
}

func Errorw(ctx context.Context, message string, args ...interface{}) {
	args = appendRequestIDIntoArgs(ctx, args)
	sugaredLogger.Errorw(message, args...)
}

func Errorf(ctx context.Context, message string, args ...interface{}) {
	sugaredLogger.Errorf(message, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	sugaredLogger.Error(args...)
}

func Infow(ctx context.Context, message string, args ...interface{}) {
	args = appendRequestIDIntoArgs(ctx, args)
	sugaredLogger.Infow(message, args...)
}

func Infof(ctx context.Context, message string, args ...interface{}) {
	sugaredLogger.Infof(message, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	sugaredLogger.Info(args...)
}

func Warnw(ctx context.Context, message string, args ...interface{}) {
	args = appendRequestIDIntoArgs(ctx, args)
	sugaredLogger.Warnw(message, args...)
}

func Warnf(ctx context.Context, message string, args ...interface{}) {
	sugaredLogger.Warnf(message, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	sugaredLogger.Warn(args...)
}

func Debugw(ctx context.Context, message string, args ...interface{}) {
	args = appendRequestIDIntoArgs(ctx, args)
	sugaredLogger.Debugw(message, args...)
}

func Debugf(ctx context.Context, message string, args ...interface{}) {
	sugaredLogger.Debugf(message, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	sugaredLogger.Debug(args...)
}

func InfoAudit(ctx context.Context, auditServiceName string, message string, args ...interface{}) {
	message = adjustAuditMessage(message, auditServiceName)
	args = appendRequestIDIntoArgs(ctx, args)
	sugaredLogger.Infow(message, args...)
}

func WarnAudit(ctx context.Context, auditServiceName string, message string, args ...interface{}) {
	message = adjustAuditMessage(message, auditServiceName)
	args = appendRequestIDIntoArgs(ctx, args)
	sugaredLogger.Warnw(message, args...)
}

// This should be rarely used
func ErrorAudit(ctx context.Context, auditServiceName string, message string, args ...interface{}) {
	message = adjustAuditMessage(message, auditServiceName)
	args = appendRequestIDIntoArgs(ctx, args)
	sugaredLogger.Errorw(message, args...)
}

func appendRequestIDIntoArgs(ctx context.Context, args []interface{}) []interface{} {
	ridValue, ok := ctx.Value(constants.RequestIDKey).(string)
	if !ok {
		return args
	}
	args = append(args, "request-id")
	args = append(args, ridValue)
	return args
}

func adjustAuditMessage(message string, serviceName string) (auditMessage string) {
	return ("[" + serviceName + "-audit] " + message)
}
