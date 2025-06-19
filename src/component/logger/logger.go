package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"MetaFarmBankend/src/component/config"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func InitLogger(cfg *config.Config) error {
	logCfg := cfg.Log

	// 设置日志级别
	var level zapcore.Level
	switch strings.ToLower(logCfg.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 创建日志目录
	if err := os.MkdirAll(logCfg.Path, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// 配置日志输出
	cores := make([]zapcore.Core, 0)

	// 控制台输出
	if logCfg.Mode == "console" || logCfg.Mode == "both" {
		consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	}

	// 文件输出
	if logCfg.Mode == "file" || logCfg.Mode == "both" {
		fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})

		logFile := filepath.Join(logCfg.Path, logCfg.ServiceName+".log")
		lumberJackLogger := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    100, // MB
			MaxBackups: 30,
			MaxAge:     logCfg.LeepDays,
			Compress:   logCfg.Compress,
		}
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberJackLogger), level))
	}

	// 创建Logger
	core := zapcore.NewTee(cores...)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = Logger.Sugar()

	// 启动日志压缩任务
	if logCfg.Compress {
		go compressOldLogs(logCfg.Path, logCfg.LeepDays)
	}

	return nil
}

func compressOldLogs(logDir string, keepDays int) {
	for {
		time.Sleep(24 * time.Hour)

		cutoff := time.Now().AddDate(0, 0, -keepDays)

		filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(path, ".log") || info.ModTime().After(cutoff) {
				return nil
			}

			compressLogFile(path)
			return nil
		})
	}
}

func compressLogFile(src string) {
	dst := src + ".gz"

	srcFile, err := os.Open(src)
	if err != nil {
		Error("Failed to open log file for compression", "file", src, "error", err)
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		Error("Failed to create compressed log file", "file", dst, "error", err)
		return
	}
	defer dstFile.Close()

	gzipWriter := gzip.NewWriter(dstFile)
	defer gzipWriter.Close()

	if _, err := io.Copy(gzipWriter, srcFile); err != nil {
		Error("Failed to compress log file", "file", src, "error", err)
		return
	}

	if err := os.Remove(src); err != nil {
		Error("Failed to remove original log file after compression", "file", src, "error", err)
	}
}

// 日志级别方法
func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

func Dump(value interface{}) {
	spew.Dump(value)
}
