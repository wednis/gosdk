package gosdk

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO 增加日志切割功能 gopkg.in/natefinch/lumberjack
// TODO 测试运行时更改日志级别是否有用

/*
func InitLogger(logPath string, level zapcore.Level) *zap.Logger {
    // 日志轮转配置
    lumberjackLogger := &lumberjack.Logger{
        Filename:   logPath,
        MaxSize:    100,
        MaxBackups: 3,
        MaxAge:     28,
        Compress:   true,
    }

    // 编码器
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "time",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    // 核心
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(encoderConfig),
        zapcore.AddSync(lumberjackLogger),
        level,
    )

    // 创建 logger
    return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}
*/

// 获取zap日志记录器
//   - debug 指定是否需要debug级别（开发环境时使用）
//   - logdir 指定日志文件目录
func NewZapLogger(debug bool, logdir string) *zap.Logger {
	// 控制台打印设置
	consoleConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000 -0700"),
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 大写 带颜色
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// 日志文件打印设置
	fileconfig := zap.NewDevelopmentEncoderConfig()
	fileconfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000 -0700")

	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	fileEncoder := zapcore.NewJSONEncoder(fileconfig)

	// 创建日志目录及目录文件
	err := os.MkdirAll(logdir, 0755)
	if err != nil {
		panic(err.Error())
	}
	file, err := os.Create(filepath.Join(logdir, "app-"+time.Now().Format("2006-01-02--15-04-05")+".log"))
	if err != nil {
		panic(err.Error())
	}

	// stderr 和 file
	cores := []zapcore.Core{
		zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stderr),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel // 开发/生产环境中 stderr 打印ERROR级别及以上日志
			})),
		zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.InfoLevel // 开发/生产环境中 file 输出INFO级别及以上信息
			})),
	}

	// 根据是否需要debug判断 stdout 日志处理
	if debug {
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.DebugLevel // 开发环境中控制台打印DEBUG级别及以上日志
			})))
	} else {
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl == zapcore.WarnLevel // 生产环境中 stdout 仅打印WARN级别日志
			})))
	}
	return zap.New(zapcore.NewTee(cores...))
}
