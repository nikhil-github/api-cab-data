package wiring

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	Config *Config
}

func (a App) Run() {
	cfg := a.Config
	err := godotenv.Load()
	if err == nil {
		log.Println("Loaded .env file")
	}

	err = envconfig.Init(&cfg)
	if err != nil {
		log.Fatal("Error loading config", err)
	}

	logger, err := configureLogger(cfg.LOG.Level)
	if err != nil {
		log.Fatalf("Failed to create zap logger: %s", err.Error())
	}

	if err := Start(cfg, logger); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func configureLogger(logLevel string) (*zap.Logger, error) {
	var level zapcore.Level
	if logLevel == "INFO" {
		level = zapcore.InfoLevel
	} else {
		level = zapcore.ErrorLevel
	}

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
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
