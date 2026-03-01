package logger

import (
    "os"
    
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(level string) (*zap.Logger, error) {
    // Parse log level
    var zapLevel zapcore.Level
    err := zapLevel.UnmarshalText([]byte(level))
    if err != nil {
        zapLevel = zapcore.InfoLevel
    }

    // Configure encoder
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "timestamp",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        FunctionKey:    zapcore.OmitKey,
        MessageKey:     "message",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    // Create core
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(encoderConfig),
        zapcore.AddSync(os.Stdout),
        zapLevel,
    )

    // Add options
    options := []zap.Option{
        zap.AddCaller(),
        zap.AddStacktrace(zapcore.ErrorLevel),
    }

    return zap.New(core, options...), nil
}