//日志相关
package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"ggs/conf"
)

//错误级别
const (
	debugLevel = 0
	infoLevel  = 1
	warnLevel  = 2
	errorLevel = 3
	fatalLevel = 4
)

//打印对应的错误级别
const (
	printDebugLevel = "[debug] "
	printInfoLevel  = "[info ] "
	printWarnLevel  = "[warn ] "
	printErrorLevel = "[error] "
	printFatalLevel = "[fatal] "
)

type Logger struct {
	level      int         //错误级别
	baseLogger *log.Logger //log包的Logger类型引用
	file       *os.File    //os包的File类型引用
}

var gLogger *Logger //定义一个Logger对象

//初始化
func init() {
	var level int
	switch strings.ToLower(conf.Env.LogLevel) { //判断错误级别
	case "debug":
		level = debugLevel
	case "info":
		level = infoLevel
	case "warn":
		level = warnLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		panic("unknown logger level: " + conf.Env.LogLevel)
	}

	var baseLogger *log.Logger
	var file *os.File
	if conf.Env.LogPath != "" { //日志路径不为空
		now := time.Now()

		filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second()) //确定日志文件名

		file, err := os.Create(path.Join(conf.Env.LogPath, filename)) //创建日志文件
		if err != nil {
			panic("cannot create log file")
		}

		baseLogger = log.New(file, "", log.LstdFlags) //创建一个记录日志的对象，将日志写入到文件
	} else {
		baseLogger = log.New(os.Stdout, "", log.LstdFlags) //创建一个记录日志的对象，将日志写入到控制台
	}

	gLogger = new(Logger) //将当前日志信息保存到gLogger中
	gLogger.level = level
	gLogger.baseLogger = baseLogger
	gLogger.file = file
}

func (logger *Logger) Close() {
	if logger.file != nil { //表示日志记录在文件中
		logger.file.Close() //关闭日志文件
	}

	logger.baseLogger = nil
	logger.file = nil
}

//打印日志信息(内部)
func (logger *Logger) printf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Printf(format, a...) //打印日志信息

	if level == fatalLevel { //如果为严重错误级别则终止程序
		os.Exit(1)
	}
}

//打印Debug级别的日志信息
func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.printf(debugLevel, printDebugLevel, format, a...)
}

//打印Info级别的日志信息
func (logger *Logger) Info(format string, a ...interface{}) {
	logger.printf(infoLevel, printInfoLevel, format, a...)
}

//打印Warn级别的日志信息
func (logger *Logger) Warn(format string, a ...interface{}) {
	logger.printf(warnLevel, printWarnLevel, format, a...)
}

//打印Error级别的日志信息
func (logger *Logger) Error(format string, a ...interface{}) {
	logger.printf(errorLevel, printErrorLevel, format, a...)
}

//打印Fatal级别的日志信息
func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.printf(fatalLevel, printFatalLevel, format, a...)
}

/* 如果日志记录在文件中，则将对应级别的日志信息写入在文件中。否则在控制台输出 */
func Debug(format string, a ...interface{}) {
	gLogger.Debug(format, a...)
}

func Info(format string, a ...interface{}) {
	gLogger.Info(format, a...)
}

func Warn(format string, a ...interface{}) {
	gLogger.Warn(format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.Error(format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.Fatal(format, a...)
}

func Close() {
	gLogger.Close()
}
