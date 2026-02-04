package Log

import (
	"RideMarket-CleanWebApi-GoLang/Config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggingLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

func newZapLogger(config *Config.Config) *ZapLogger {
	logger := &ZapLogger{config: config}
	logger.Init()
	return logger
}

type ZapLogger struct {
	config *Config.Config
	logger *zap.SugaredLogger
}

func (l *ZapLogger) Init() {
	writeSinker := zapcore.AddSync(&lumberjack.Logger{
		Filename:   l.config.Logger.FilePath,
		MaxSize:    10,
		MaxAge:     5,
		MaxBackups: 10,
		Compress:   true,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(config),
		writeSinker,
		l.getLogLevel())

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()

	l.logger = logger

}

func (l *ZapLogger) getLogLevel() zapcore.Level {
	level, exists := loggingLevelMap[l.config.Logger.Level]
	if !exists {
		return zapcore.DebugLevel
	}
	return level
}

func (l *ZapLogger) Info(cat Category, subCat SubCategory, message string,
	extra map[Extrakey]interface{}) {

	params := prepareLogKeys(cat, subCat, extra)

	l.logger.Infow(message, params...)

}

func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.Infof(template, args...)
}

func (l *ZapLogger) Debug(cat Category, subCat SubCategory, message string,
	extra map[Extrakey]interface{}) {

	params := prepareLogKeys(cat, subCat, extra)

	l.logger.Debugw(message, params...)
}

func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	l.Debugf(template, args...)
}

func (l *ZapLogger) Warning(cat Category, subCat SubCategory, message string,
	extra map[Extrakey]interface{}) {
	params := prepareLogKeys(cat, subCat, extra)

	l.logger.Warnw(message, params...)
}

func (l *ZapLogger) Warningf(template string, args ...interface{}) {
	l.Warningf(template, args...)
}

func (l *ZapLogger) Error(cat Category, subCat SubCategory, message string,
	extra map[Extrakey]interface{}) {
	params := prepareLogKeys(cat, subCat, extra)

	l.logger.Errorw(message, params...)
}

func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.Errorf(template, args...)
}

func (l *ZapLogger) Fatal(cat Category, subCat SubCategory, message string,
	extra map[Extrakey]interface{}) {
	params := prepareLogKeys(cat, subCat, extra)

	l.logger.Fatalw(message, params...)
}

func (l *ZapLogger) Fatalf(template string, args ...interface{}) {
	l.Fatalf(template, args...)
}

func prepareLogKeys(cat Category, subCat SubCategory, extra map[Extrakey]interface{}) []interface{} {
	if extra == nil {
		extra = make(map[Extrakey]interface{})
	}

	extra["Category"] = cat
	extra["SubCategory"] = subCat

	params := MapToZapParams(extra)
	return params
}
