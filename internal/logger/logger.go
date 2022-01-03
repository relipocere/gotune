package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//NewLogger creates new preconfigured Sugared Zap logger.
func NewLogger(options ...func(*zap.Config)) (*zap.SugaredLogger, error) {
	cfg := &zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.ErrorLevel),
		OutputPaths: []string{"stdout"},

		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",

			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}

	for _, option := range options {
		option(cfg)
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

//SetLevel returns function that sets logger level.
func SetLevel(level string) func(*zap.Config) {
	return func(c *zap.Config) {
		switch level {
		case "DEBUG":
			c.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		case "INFO":
			c.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		case "WARN":
			c.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		case "FATAL":
			c.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
		case "PANIC":
			c.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
		default:
			c.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		}
	}
}
