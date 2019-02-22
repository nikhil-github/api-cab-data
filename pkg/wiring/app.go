package wiring

import (
	"context"
	"log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	Config    Config
	logger 	  *zap.Logger
}

func (a App) Run() {
	ctx, _ := context.WithCancel(context.Background())
	// dev logging is enabled always.
	logger ,err := configureLogger(true)
	if err != nil {
		log.Fatalf("Failed to create zap logger: %s", err.Error())
	}

	if err := StartServer(ctx,"api-cab-data",logger); err != nil {
		a.logger.Fatal("Failed to start server", zap.Error(err))
	}
}


func configureLogger(devFlag bool) (*zap.Logger,error) {
	var level zapcore.Level
	if devFlag {
		level = zapcore.InfoLevel
	} else {
		level = zapcore.WarnLevel
	}

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Unable to build zap logger: %s", err.Error())
	}
	logger.Info("info logging enabled")
	return logger, nil
}
