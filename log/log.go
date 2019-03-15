package log

import (
	"sync"
	"time"

	iEnv "github.com/g-airport/go-infra/env"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//-------------------------------------------------*
// 基于 zap 的 log level
//-------------------------------------------------*
var LEVELS = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"err":   zap.ErrorLevel,
}

type Logger struct {
	path  string
	rLog  *lumberjack.Logger
	log   *zap.Logger
	sugar *zap.SugaredLogger

	level zapcore.Level
	pid   []interface{}

	rolling        bool
	lastRotateTime time.Time
	lastRotateRW   sync.Mutex
}

func NewLogger(path string, level string) (*Logger, error) {
	out := new(Logger)
	out.rLog = new(lumberjack.Logger)

	out.path = path
	out.lastRotateTime = time.Now()
	out.level = LEVELS[level]
	out.pid = []interface{}{iEnv.Pid}

	// config lumberjack
	out.rLog.Filename = path
	// automatic rolling file on it increment than 2GB
	out.rLog.MaxSize = 0x1000 * 5
	out.rLog.LocalTime = true
	out.rLog.Compress = true
	// reserve last 60 day logs
	out.rLog.MaxBackups = 60

	// config encoder config
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeLevel = zapcore.CapitalLevelEncoder
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	// config core
	c := zapcore.AddSync(out.rLog)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(ec), c, out.level)
	out.log = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	).
		With(zap.Int("pid", iEnv.Pid))

	// default enable daily rotate
	out.rolling = true

	out.sugar = out.log.Sugar()
	return out, nil
}

func (tLog *Logger) checkRotate() {
	if !tLog.rolling {
		return
	}

	n := time.Now()

	tLog.lastRotateRW.Lock()
	defer tLog.lastRotateRW.Unlock()

	last := tLog.lastRotateTime
	y, m, d := last.Year(), last.Month(), last.Day()
	if y != n.Year() || m != n.Month() || d != n.Day() {
		go tLog.rLog.Rotate()
		tLog.lastRotateTime = n
	}
}

func (tLog *Logger) EnableDailyFile() {
	tLog.rolling = true
}

func (tLog *Logger) Debug(format string, v ...interface{}) {
	tLog.checkRotate()
	if !tLog.level.Enabled(zap.DebugLevel) {
		return
	}

	defer tLog.log.Sync()
	tLog.sugar.Debugf(format, v...)
}

func (tLog *Logger) Info(format string, v ...interface{}) {
	tLog.checkRotate()
	if !tLog.level.Enabled(zap.InfoLevel) {
		return
	}

	defer tLog.log.Sync()
	tLog.sugar.Infof(format, v...)
}

func (tLog *Logger) Warn(format string, v ...interface{}) {
	tLog.checkRotate()
	if !tLog.level.Enabled(zap.WarnLevel) {
		return
	}

	defer tLog.log.Sync()
	tLog.sugar.Warnf(format, v...)
}

func (tLog *Logger) Err(format string, v ...interface{}) {
	tLog.checkRotate()
	if !tLog.level.Enabled(zap.ErrorLevel) {
		return
	}

	defer tLog.log.Sync()
	tLog.sugar.Errorf(format, v...)
}

//-------------------------------------------------*
// json format
//-------------------------------------------------*
func (tLog *Logger) Debugw(format string, v ...interface{}) {
	tLog.checkRotate()
	if !tLog.level.Enabled(zap.DebugLevel) {
		return
	}

	defer tLog.log.Sync()
	tLog.sugar.Debugw(format, v...)
}

func (tLog *Logger) Infow(format string, v ...interface{}) {
	tLog.checkRotate()
	if !tLog.level.Enabled(zap.InfoLevel) {
		return
	}

	defer tLog.log.Sync()
	tLog.sugar.Infow(format, v...)
}

func (tLog *Logger) Warnw(format string, v ...interface{}) {
	tLog.checkRotate()
	if !tLog.level.Enabled(zap.WarnLevel) {
		return
	}

	defer tLog.log.Sync()
	tLog.sugar.Warnw(format, v...)
}

func (tLog *Logger) Errw(format string, v ...interface{}) {
	tLog.checkRotate()
	if !tLog.level.Enabled(zap.ErrorLevel) {
		return
	}

	defer tLog.log.Sync()
	tLog.sugar.Errorw(format, v...)
}

var _logger *Logger

func GetDefault() *Logger {
	return _logger
}

func SetDefault(l *Logger) {
	_logger = l
}

func Stdout() {
	l, _ := NewLogger("stdout", "debug")
	SetDefault(l)
}

func Emergency(format string, v ...interface{}) {
	_logger.Err(format, v...)
}

func Alert(format string, v ...interface{}) {
	_logger.Err(format, v...)
}

func Critical(format string, v ...interface{}) {
	_logger.Err(format, v...)
}

func Notice(format string, v ...interface{}) {
	_logger.Info(format, v...)
}

// base
func Err(format string, v ...interface{}) {
	_logger.Err(format, v...)
}

func Warn(format string, v ...interface{}) {
	_logger.Warn(format, v...)
}

func Info(format string, v ...interface{}) {
	_logger.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	_logger.Debug(format, v...)
}

// json format
func Errw(format string, v ...interface{}) {
	_logger.Errw(format, v...)
}

func Warnw(format string, v ...interface{}) {
	_logger.Warnw(format, v...)
}

func Infow(format string, v ...interface{}) {
	_logger.Infow(format, v...)
}

func Debugw(format string, v ...interface{}) {
	_logger.Debugw(format, v...)
}
