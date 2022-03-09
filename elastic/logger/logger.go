package logger

import (
	"path"
	"runtime"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

// init 初始化SugaredLogger
func init() {
	writeSyncer := getLogWriter()
    encoder := getEncoder()
    core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

    logger := zap.New(core)
        
    sugarLogger = logger.Sugar()
}

// log 日志
func log(level string, logTag string, msg string, args ...interface{}) {
	level = strings.ToLower(level)
	msg = logTag + " " + msg
	switch level {
	case "debug":
		sugarLogger.Debugf(msg, args...)
	case "info":
		sugarLogger.Infof(msg, args...)
	case "warning":
		sugarLogger.Warnf(msg, args...)
	case "error":
		sugarLogger.Errorf(msg, args...)
	default:
	}
}

// D debug
func D(logTag string, msg string, args ...interface{}) {
	log("debug", logTag, msg, args...)
}

// I info
func I(logTag string, msg string, args ...interface{}) {
	log("info", logTag, msg, args...)
}

// W warning
func W(logTag string, msg string, args ...interface{}) {
	log("warning", logTag, msg, args...)
}

// E error
func E(logTag string, msg string, args ...interface{}) {
	log("error", logTag, msg, args...)
}

// getEncoder 获取日志编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 获取日志输出句柄
func getLogWriter() zapcore.WriteSyncer {
	_, curFile, _, _ := runtime.Caller(0)
	logPath := path.Join(path.Dir(path.Dir(curFile)), "./log/elastic.log")	
	
	// 日志切分
    lumberJackLogger := &lumberjack.Logger{
        Filename:   logPath,
        MaxSize:    10,
        MaxBackups: 5,
        MaxAge:     30,
        Compress:   false,
    }
    return zapcore.AddSync(lumberJackLogger)
}