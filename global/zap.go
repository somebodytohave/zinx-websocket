package global

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

var _level zapcore.Level

var Glog *zap.Logger

func InitZap() {
	if Object == nil {
		panic("plz init global object by global.InitObject()")
	}
	zapConfig := Object.ZapConfig
	dir := zapConfig.Director
	if ok, _ := PathExists(dir); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", dir)
		_ = os.Mkdir(dir, os.ModePerm)
	}

	switch zapConfig.Level { // 初始化配置文件的Level
	case "debug":
		_level = zap.DebugLevel
	case "info":
		_level = zap.InfoLevel
	case "warn":
		_level = zap.WarnLevel
	case "error":
		_level = zap.ErrorLevel
	case "dpanic":
		_level = zap.DPanicLevel
	case "panic":
		_level = zap.PanicLevel
	case "fatal":
		_level = zap.FatalLevel
	default:
		_level = zap.InfoLevel
	}

	if _level == zap.DebugLevel || _level >= zap.WarnLevel {
		Glog = zap.New(getEncoderCore(), zap.AddStacktrace(_level))
	} else {
		Glog = zap.New(getEncoderCore())
	}
	if zapConfig.ShowLine {
		Glog = Glog.WithOptions(zap.AddCaller())
	}
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  Object.ZapConfig.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case Object.ZapConfig.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case Object.ZapConfig.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case Object.ZapConfig.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case Object.ZapConfig.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if Object.ZapConfig.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := rotatelogs.New(
		path.Join(Object.ZapConfig.Director, "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if Object.ZapConfig.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	writer, err := GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(), writer, _level)
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(Object.ZapConfig.Prefix + "2006/01/02 - 15:04:05.000"))
}
