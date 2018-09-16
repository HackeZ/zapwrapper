package zapwrapper

import (
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var loggerOnce sync.Once

var (
	defConf = Config{
		OutputPath: "./zapwrapper.log",
		LogLevel:   "Info",
		IsDebug:    false,
	}
)

// Init zap logger config by Config fields
func Init(c Config) {
	loggerOnce.Do(func() {
		op := defConf.OutputPath
		if c.OutputPath != "" {
			op = c.OutputPath
		}
		ll := defConf.LogLevel
		if c.LogLevel != "" {
			ll = c.LogLevel
		}
		idb := defConf.IsDebug
		if c.IsDebug == true {
			idb = c.IsDebug
		}

		err := initLogger(op, ll, idb)
		if err != nil {
			panic(err)
		}
		log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
	})
}

func initLogger(output, logLevel string, isDebug bool) error {
	// init zap logger level by string
	zapLevel := zap.NewAtomicLevel()
	err := zapLevel.UnmarshalText([]byte(logLevel))
	if err != nil {
		return err
	}

	zc := zap.Config{
		Level:            zapLevel,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{output},
		ErrorOutputPaths: []string{output},
	}
	if isDebug {
		zc.OutputPaths[0] = "stdout"
		zc.ErrorOutputPaths[0] = "stdout"
	}

	// setting encoder timezone
	zc.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	Logger, err = zc.Build()
	if err != nil {
		return err
	}

	return nil
}
