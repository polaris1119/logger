// Copyright 2016 polaris. All rights reserved.
// Use of l source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	polaris@studygolang.com

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/context"
)

const (
	ERROR = iota
	INFO
	DEBUG
)

var levelMap = map[string]int{
	"ERROR": ERROR,
	"INFO":  INFO,
	"DEBUG": DEBUG,
}

var (
	// 日志文件
	infoFile  = ""
	debugFile = ""
	errorFile = ""

	accessFile = ""

	level int
)

// Init Init("", "INFO")
func Init(logPath, tmpLevel string) {

	os.Mkdir(logPath, 0777)

	infoFile = logPath + "/info.log"
	debugFile = logPath + "/debug.log"
	errorFile = logPath + "/error.log"

	accessFile = logPath + "/access.log"

	level = levelMap[strings.ToUpper(tmpLevel)]
}

func Infof(format string, args ...interface{}) {
	if level < INFO {
		return
	}

	file, err := openFile(infoFile)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
}

func Infoln(args ...interface{}) {
	if level < INFO {
		return
	}

	file, err := openFile(infoFile)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Println(args...)
}

func Errorf(format string, args ...interface{}) {
	file, err := openFile(errorFile)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
}

func Errorln(args ...interface{}) {
	file, err := openFile(errorFile)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Println(args...)
}

func Debugf(format string, args ...interface{}) {
	if level < DEBUG {
		return
	}

	file, err := openFile(debugFile)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
}

func Debugln(args ...interface{}) {
	if level < DEBUG {
		return
	}

	file, err := openFile(debugFile)
	if err != nil {
		return
	}
	defer file.Close()
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"文件：", filepath.Base(callerFile), "行号:", line}, args...)
	}
	New(file).Println(args...)
}

func openFile(filename string) (*os.File, error) {
	if filename == "" {
		log.Println("[WARNING] You must call logger.Init function First!")
		return nil, fmt.Errorf("[WARNING] You must call logger.Init function First!")
	}

	filename += "-" + time.Now().Format("060102")

	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
}

type Logger struct {
	*log.Logger

	// TODO:append 数据时，没有加锁，所以，如果同一个 Logger 实例，多个goroutine并发可能顺序会乱
	infoBuf  []interface{}
	errorBuf []interface{}
	debugBuf []interface{}
	ctx      context.Context
}

func New(out io.Writer) *Logger {
	return &Logger{
		Logger: log.New(out, "", log.Lmicroseconds),
	}
}

func NewLoggerContext(ctx context.Context) *Logger {

	objLogger := &Logger{
		ctx: ctx,
	}

	// 第一个元素用于最后 flush 时存 uri 信息
	objLogger.resetBuf()

	return objLogger
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.appendInfo(fmt.Sprintf(format, args...))
}

func (l *Logger) Infoln(args ...interface{}) {
	l.appendInfo(fmt.Sprintln(args...))
}

func (l *Logger) appendInfo(info string) {
	if level < INFO {
		return
	}

	l.infoBuf = append(l.infoBuf, info)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.appendError(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorln(args ...interface{}) {
	l.appendError(fmt.Sprintln(args...))
}

func (l *Logger) appendError(errstr string) {
	l.errorBuf = append(l.errorBuf, errstr)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.appendDebug(fmt.Sprintf(format, args...))
}

func (l *Logger) Debugln(args ...interface{}) {
	l.appendDebug(fmt.Sprintln(args...))
}

func (l *Logger) appendDebug(debugstr string) {
	if level < DEBUG {
		return
	}

	l.debugBuf = append(l.debugBuf, debugstr)
}

func (l *Logger) SetContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *Logger) AccessLog(args ...interface{}) {
	file, err := openFile(accessFile)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Println(args...)
}

func (l *Logger) resetBuf() {
	l.infoBuf = make([]interface{}, 1, 20)
	l.errorBuf = make([]interface{}, 1, 20)
	l.debugBuf = make([]interface{}, 1, 5)
}

func (l *Logger) Flush() {

	var (
		file *os.File
		err  error

		uri = l.ctx.Value("uri")
	)

	if len(l.infoBuf) > 1 {
		file, err = openFile(infoFile)
		if err == nil {
			l.Logger = log.New(file, "", log.Ltime)
			defer file.Close()

			l.infoBuf[0] = uri
			l.Println(l.infoBuf...)
		}
	}

	if len(l.errorBuf) > 1 {
		file, err = openFile(errorFile)
		if err == nil {
			l.Logger = log.New(file, "", log.Ltime)
			defer file.Close()

			l.errorBuf[0] = uri
			l.Println(l.errorBuf...)
		}
	}

	if len(l.debugBuf) > 1 {
		file, err = openFile(debugFile)
		if err == nil {
			l.Logger = log.New(file, "", log.Ltime)
			defer file.Close()

			l.debugBuf[0] = uri
			l.Println(l.debugBuf...)
		}
	}

	l.resetBuf()
}
