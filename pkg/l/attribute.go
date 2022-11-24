package l

import (
	"gotag/pkg/v"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logging struct {
	Level      string
	Path       string
	Console    bool //是否输出打印到控制台
	MaxSize    int  //single log file maxsize before compress
	MaxAge     int  //time to live
	MaxBackups int  // max log files counter
}

type LoggerAttr struct {
	Logging logging
}

func (c *LoggerAttr) InitDefaultLogger() *LoggerAttr {
	c.Logging.Level = v.GetViper().GetString("log.level")
	c.Logging.Path = "./logs/"
	c.Logging.Console = v.GetViper().GetBool("log.console_enabled")
	c.Logging.MaxSize = 1
	c.Logging.MaxAge = 60
	c.Logging.MaxBackups = 30

	return c
}

func (c *LoggerAttr) SetLogPath(model string) {
	if len(c.Logging.Path) == 0 {
		return
	}

	c.Logging.Path = c.Logging.Path + "/" + model + ".log"
}

func (c *LoggerAttr) NewLogger(opt ...zap.Option) (*zap.Logger, error) {
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte(c.Logging.Level))
	if err != nil {
		return nil, err
	}

	hook := lumberjack.Logger{
		Filename:   c.Logging.Path,       // 日志文件路径
		MaxSize:    c.Logging.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: c.Logging.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     c.Logging.MaxAge,     // 文件最多保存多少天
		Compress:   true,                 // 是否压缩
	}

	writeSyncer := make([]zapcore.WriteSyncer, 0)
	if len(c.Logging.Path) > 0 {
		writeSyncer = append(writeSyncer, zapcore.AddSync(&hook))
	}

	if c.Logging.Console {
		writeSyncer = append(writeSyncer, zapcore.AddSync(os.Stdout))
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	_, _ = zap.NewDevelopment()

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),    // 编码器配置
		zapcore.NewMultiWriteSyncer(writeSyncer...), // 打印到控制台和文件
		zap.NewAtomicLevelAt(level),                 // 日志级别
	)

	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	logger := zap.New(core, caller, development)
	return logger, nil
}
