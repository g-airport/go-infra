package log

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	log   *zap.Logger
	rLog  *lumberjack.Logger
	sugar *zap.SugaredLogger
	level zapcore.Level

	rolling    bool
	path       string
	lastRotate time.Time

	rotateMu *sync.Mutex
}

func Init(path string, level zapcore.Level) error {
	l, err := NewLogger(path, level)
	if err != nil {
		return err
	}

	SetDefault(l)
	return nil
}

func Sync() {
	_ = _logger.log.Sync()
}

func NewLogger(path string, level zapcore.Level) (*Logger, error) {
	return NewLoggerWithEncoderConfig(path, level, nil)
}

func NewLoggerWithEncoderConfig(path string, level zapcore.Level, config *zapcore.EncoderConfig) (*Logger, error) {
	out := new(Logger)
	out.rLog = new(lumberjack.Logger)

	out.path = path
	out.lastRotate = time.Now()
	out.level = level
	out.rotateMu = &sync.Mutex{}

	// config lumberjack
	out.rLog.Filename = path
	out.rLog.MaxSize = 0x1000 * 2 // automatic rolling file on it increment than 2GB
	out.rLog.LocalTime = true
	out.rLog.Compress = true
	out.rLog.MaxBackups = 60 // reserve last 60 day logs
	out.rolling = true

	// config encoder
	var cfg zapcore.EncoderConfig
	if config == nil {
		// default config
		cfg = zap.NewProductionEncoderConfig()
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = *config
	}

	// config core
	c := zapcore.AddSync(out.rLog)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), c, out.level)
	out.log = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	)

	out.sugar = out.log.Sugar()
	return out, nil
}

func (log *Logger) checkRotate() {
	if !log.rolling {
		return
	}

	n := time.Now()
	if log.differentDay(n) {
		log.rotateMu.Lock()
		defer log.rotateMu.Unlock()

		if log.differentDay(n) {
			_ = log.rLog.Rotate()
			log.lastRotate = n
		}
	}
}

func (log *Logger) differentDay(t time.Time) bool {
	y, m, d := log.lastRotate.Year(), log.lastRotate.Month(), log.lastRotate.Day()
	return y != t.Year() || m != t.Month() || d != t.Day()
}

func (log *Logger) EnableDailyFile() {
	log.rolling = true
}

func (log *Logger) Err(format string, v ...interface{}) {
	log.checkRotate()
	log.sugar.Errorf(format, v...)
}

func (log *Logger) Errw(format string, v ...interface{}) {
	log.checkRotate()
	log.sugar.Errorw(format, v...)
}

func (log *Logger) Warn(format string, v ...interface{}) {
	log.checkRotate()
	log.sugar.Warnf(format, v...)
}

func (log *Logger) Warnw(format string, v ...interface{}) {
	log.checkRotate()
	log.sugar.Warnw(format, v...)
}

func (log *Logger) Info(format string, v ...interface{}) {
	log.checkRotate()
	log.sugar.Infof(format, v...)
}

func (log *Logger) Infow(format string, v ...interface{}) {
	log.checkRotate()
	log.sugar.Infow(format, v...)
}

func (log *Logger) Debug(format string, v ...interface{}) {
	log.checkRotate()
	log.sugar.Debugf(format, v...)
}

func (log *Logger) Debugw(format string, v ...interface{}) {
	log.checkRotate()
	log.sugar.Debugw(format, v...)
}

func (log *Logger) Print(v ...interface{}) {
	log.checkRotate()
	log.sugar.Info(v...)
}

func (log *Logger) Printf(format string, v ...interface{}) {
	log.checkRotate()
	log.sugar.Info(fmt.Sprintf(format, v...))
}

func (log *Logger) GetUnderlyingLogger() *zap.Logger {
	return log.log
}

var _logger, _ = NewLogger("log/common.log", zap.InfoLevel)

func GetDefault() *Logger {
	return _logger
}

func SetDefault(l *Logger) {
	_logger = l
}

func Stdout() *Logger {
	l, _ := NewLogger("stdout", zap.InfoLevel)
	return l
}

func Err(format string, v ...interface{}) error {
	_logger.Err(format, v...)
	return fmt.Errorf(format, v...)
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
