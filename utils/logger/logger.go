package logger

import (
	"github.com/sdlp99/sdpkg/utils/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"sync"
)

var (
	gLogger *zap.Logger
	logMap  = struct {
		sync.RWMutex
		m map[string]*zap.Logger
	}{m: make(map[string]*zap.Logger)}
)

func init() {
	//println("logger init")
}

func InitGlobalLogger() {

	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	//defer gLogger.Sync()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   customLevelEncoder, // 小写编码器
		EncodeTime:    zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		//EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}

	logLevel := config.GetConfig("logger.logLevel", "DEBUG")
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)

	switch logLevel {
	case "DEBUG":
		atom = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "INFO":
		atom = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "WARN":
		atom = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "ERROR":
		atom = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "DPANIC":
		atom = zap.NewAtomicLevelAt(zapcore.DPanicLevel)
	case "PANIC":
		atom = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case "FATAL":
		atom = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	}

	// 设置日志级别
	//config := zap.Config{
	//	Level:            atom,                                                // 日志级别
	//	Development:      false,                                                // 开发模式，堆栈跟踪
	//	Encoding:         "console",                                              // 输出格式 console 或 json
	//	EncoderConfig:    encoderConfig,                                       // 编码器配置
	//	InitialFields:    nil, // 初始化字段，如：添加一个服务器名称
	//	OutputPaths:      []string{"stdout", "test.log"},         // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
	//	ErrorOutputPaths: []string{"stderr"},
	//}
	logDir := config.GetConfig("logger.logDir", ".")
	fileName := config.GetConfig("logger.fileName", "log.log")
	if !strings.HasSuffix(logDir, "/") {
		fileName = "/" + fileName
	}
	maxSize := config.GetConfigInt("logger.maxSize", 10)
	maxAge := config.GetConfigInt("logger.maxAge", 10)
	maxBackups := config.GetConfigInt("logger.maxBackups", 100)

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename: logDir + fileName, // ⽇志⽂件路径
		MaxSize:  maxSize,
		// 历史日志文件保留天数
		MaxAge: maxAge,
		// 最大保留历史日志数量
		MaxBackups: maxBackups,
		LocalTime:  true,  // 采用本地时间
		Compress:   false, // 是否压缩日志
	})

	//writer2:=zapcore.AddSync(&lumberjack.Logger{
	//	Filename: "./log/test2.log", // ⽇志⽂件路径
	//	MaxSize:    100,
	//	// 历史日志文件保留天数
	//	MaxAge:     30,
	//	// 最大保留历史日志数量
	//	MaxBackups: 100,
	//	LocalTime: false,                            // 采用本地时间
	//	Compress: false,                          // 是否压缩日志
	//})

	//zapCore := zapcore.NewTee(
	//	zapcore.NewCore(
	//		zapcore.NewConsoleEncoder(encoderConfig),
	//		zapcore.NewMultiWriteSyncer( zapcore.AddSync(writer),zapcore.AddSync(os.Stdout)),
	//		atom,
	//	),
	//)

	zapCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout)),
		atom,
	)

	gLogger = zap.New(zapCore, zap.AddCaller())

	// 构建日志
	//gLogger, _ = config.Build( )

}

func GetLogger() *zap.Logger {
	return gLogger
}

func ExitLogger() {
	defer gLogger.Sync()
}

func InitLogger(logKey string) {
	logMap.Lock()
	defer logMap.Unlock()
	if logMap.m[logKey] != nil {
		return
	}
	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	//defer gLogger.Sync()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   customLevelEncoder, // 小写编码器
		EncodeTime:    zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		//EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}

	logLevel := "INFO"
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)

	switch logLevel {
	case "DEBUG":
		atom = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "INFO":
		atom = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "WARN":
		atom = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "ERROR":
		atom = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "DPANIC":
		atom = zap.NewAtomicLevelAt(zapcore.DPanicLevel)
	case "PANIC":
		atom = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case "FATAL":
		atom = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	}

	// 设置日志级别
	//config := zap.Config{
	//	Level:            atom,                                                // 日志级别
	//	Development:      false,                                                // 开发模式，堆栈跟踪
	//	Encoding:         "console",                                              // 输出格式 console 或 json
	//	EncoderConfig:    encoderConfig,                                       // 编码器配置
	//	InitialFields:    nil, // 初始化字段，如：添加一个服务器名称
	//	OutputPaths:      []string{"stdout", "test.log"},         // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
	//	ErrorOutputPaths: []string{"stderr"},
	//}
	logDir := config.GetConfig("logger.logDir", ".")
	fileName := logKey + ".log"
	if !strings.HasSuffix(logDir, "/") {
		fileName = "/" + fileName
	}

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename: logDir + fileName, // ⽇志⽂件路径
		MaxSize:  10000,
		// 历史日志文件保留天数
		MaxAge: 0,
		// 最大保留历史日志数量
		MaxBackups: 0,
		LocalTime:  true,  // 采用本地时间
		Compress:   false, // 是否压缩日志
	})

	zapCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout)),
		atom,
	)

	gLogger = zap.New(zapCore, zap.AddCaller())

	logMap.m[logKey] = zap.New(zapCore, zap.AddCaller())
	// 构建日志
	//gLogger, _ = config.Build( )

}

func GetLogByName(name string) *zap.Logger {
	if logMap.m[name] == nil {
		InitLogger(name)
	}
	return logMap.m[name]
}
