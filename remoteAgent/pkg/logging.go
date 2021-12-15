package grpcinverter

import (
	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"
	"go.uber.org/zap/zapcore"
	"sync"

	"go.uber.org/zap"
)

var LOGGER *zap.Logger
var logLock = sync.Mutex{}

func init() {
	var err error
	LOGGER, err = zap.NewProduction(zap.AddCaller(), zap.AddStacktrace(zapcore.FatalLevel))
	if err != nil {
		panic(err)
	}
	ioc.SetLogger(LOGGER)
}

// DisableLogging is useful for making tests a bit quieter if you want to. This is pretty much
// mandatory for benchmark test, though, otherwise you will never see the ouput.
func DisableLogging() {
	logLock.Lock()
	defer logLock.Unlock()
	LOGGER = zap.New(zapcore.NewNopCore())
	ioc.SetLogger(LOGGER)
}
