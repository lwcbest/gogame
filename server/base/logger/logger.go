package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// levels
const (
	debugLevel   = 0
	releaseLevel = 1
	errorLevel   = 2
	fatalLevel   = 3
)

const (
	printDebugLevel   = "[debug] "
	printReleaseLevel = "[release] "
	printErrorLevel   = "[error] "
	printFatalLevel   = "[fatal] "
)

type Logger struct {
	level      int
	baseLogger *log.Logger
	baseFile   *os.File
	filaname   string
	pathname   string
}

var gLogger *Logger

func init() {
	gLogger, _ = New("debug", "", "default", log.LstdFlags)
}

func ResetLogger(strLevel, pathname, filename string) {
	var err error
	gLogger, err = New(strLevel, pathname, filename, log.LstdFlags)
	if err != nil {
		fmt.Printf("ResetLogger failed err %+v\n", err)
	}
}

func New(strLevel string, pathName string, fileName string, flag int) (*Logger, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "release":
		level = releaseLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown level: " + strLevel)
	}

	// logger
	var baseLogger *log.Logger
	var baseFile *os.File
	//fmt.Printf("log.New called！%+v %+v %t\n", strLevel, pathName, pathName == "")
	if pathName != "" {
		now := time.Now()
		filename := fmt.Sprintf("%d%02d%02d_%s",
			now.Year(),
			now.Month(),
			now.Day(),
			fileName)
		exist, _ := PathExists(path.Join(pathName, filename))
		var file *os.File
		var err error
		if !exist {
			file, err = os.Create(path.Join(pathName, filename))
			if err != nil {
				return nil, err
			}
			baseLogger = log.New(file, "", flag)
			baseFile = file

		} else {
			file, err = os.OpenFile(path.Join(pathName, filename), os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				return nil, err
			}
			baseLogger = log.New(file, "", flag)
			baseFile = file
		}
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	// new
	logger := new(Logger)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile
	logger.filaname = fileName
	logger.pathname = pathName

	return logger, nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {

	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	_, filePath, c, _ := runtime.Caller(2)
	format = fmt.Sprintf("%s[%s][%v] %s", printLevel, filePath, c, format)
	var (
		colPrefix string
		colPosfix string
	)
	switch printLevel {
	case "[debug] ":
		colPrefix = "\x1b[0;34m"
		colPosfix = "\x1b[0m"
		break
	case "[release] ":
		colPrefix = "\x1b[0;32m"
		colPosfix = "\x1b[0m"
		break
	case "[error] ":
		colPrefix = "\x1b[0;31m"
		colPosfix = "\x1b[0m"
		break
	case "[fatal] ":
		colPrefix = "\x1b[0;35m"
		colPosfix = "\x1b[0m"
		break
	}

	if logger.baseFile != nil {
		now := time.Now()
		filename := fmt.Sprintf("%d%02d%02d_%s",
			now.Year(),
			now.Month(),
			now.Day(),
			logger.filaname)
		exist, _ := PathExists(path.Join(logger.pathname, filename))
		if !exist {
			file, _ := os.Create(path.Join(logger.pathname, filename))
			logger.baseLogger = log.New(file, "", log.LstdFlags)
			logger.baseFile = file
		}
	}

	logger.baseLogger.Output(3, colPrefix+fmt.Sprintf(format, a...)+colPosfix)
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Release(format string, a ...interface{}) {
	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Debug(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

//兼容release
func Info(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Close() {
	gLogger.Close()
}
