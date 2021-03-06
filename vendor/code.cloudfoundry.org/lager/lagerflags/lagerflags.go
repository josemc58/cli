package lagerflags

import (
	"flag"
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

type LagerConfig struct {
	LogLevel string `json:"log_level,omitempty"`
}

func DefaultLagerConfig() LagerConfig {
	return LagerConfig{
		LogLevel: string(INFO),
	}
}

var minLogLevel string

func AddFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(
		&minLogLevel,
		"logLevel",
		string(INFO),
		"log level: debug, info, error or fatal",
	)
}

func New(component string) (lager.Logger, *lager.ReconfigurableSink) {
	return newLogger(component, minLogLevel)
}

func NewFromConfig(component string, config LagerConfig) (lager.Logger, *lager.ReconfigurableSink) {
	return newLogger(component, config.LogLevel)
}

func newLogger(component, minLogLevel string) (lager.Logger, *lager.ReconfigurableSink) {
	var minLagerLogLevel lager.LogLevel

	switch minLogLevel {
	case DEBUG:
		minLagerLogLevel = lager.DEBUG
	case INFO:
		minLagerLogLevel = lager.INFO
	case ERROR:
		minLagerLogLevel = lager.ERROR
	case FATAL:
		minLagerLogLevel = lager.FATAL
	default:
		panic(fmt.Errorf("unknown log level: %s", minLogLevel))
	}

	logger := lager.NewLogger(component)

	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), minLagerLogLevel)
	logger.RegisterSink(sink)

	return logger, sink
}
