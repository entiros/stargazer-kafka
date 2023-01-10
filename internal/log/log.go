package log

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type MyLogger struct {
	zap.SugaredLogger
}

var Logger MyLogger

var debugLevel zap.AtomicLevel

func (l MyLogger) Level() kgo.LogLevel {
	return kgo.LogLevelError
}

func (l MyLogger) Log(level kgo.LogLevel, msg string, keyvals ...interface{}) {
	l.Debugf(msg, keyvals)
}

// runs when package is loaded to init log package
func init() {

	debugLevel = zap.NewAtomicLevel()
	debugLevel.SetLevel(zap.DebugLevel)

	clusterEnv := os.Getenv("env")
	var encoderCfg zapcore.EncoderConfig
	if clusterEnv == "prd" {
		encoderCfg = zap.NewProductionEncoderConfig()
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}

	l := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		debugLevel,
	))

	Logger = MyLogger{*l.Sugar()}

}

func SetDebugLevel(level string) {

	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		debugLevel.SetLevel(zap.DebugLevel)
	case "error":
		debugLevel.SetLevel(zap.ErrorLevel)
	case "info":
		debugLevel.SetLevel(zap.InfoLevel)
	case "warn":
		debugLevel.SetLevel(zap.WarnLevel)
	case "fatal":
		debugLevel.SetLevel(zap.FatalLevel)
	default:
		debugLevel.SetLevel(zap.DebugLevel)

	}

}
