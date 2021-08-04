package logger

import (
	"PoDownloader/util"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	stdLog  *log.Logger
	fileLog *log.Logger
	file    *os.File
}

func getLogFileNameByDate() string {
	return fmt.Sprintf("%s.log", time.Now().Format("20060102"))
}

func getLogFilePathByDate(logFolder string) string {
	return filepath.Join(logFolder, getLogFileNameByDate())
}

func NewLogger(logFolder string) (*Logger, error) {
	stdLogger := log.New(os.Stderr, "", log.LstdFlags)
	// Returns only stdLogger when logFolder parameter is not specified
	if logFolder == "" {
		return &Logger{
			stdLog:  stdLogger,
			fileLog: nil,
			file:    nil,
		}, nil
	}
	err := util.EnsureDirAll(logFolder)
	if err != nil {
		return &Logger{
			stdLog:  stdLogger,
			fileLog: nil,
			file:    nil,
		}, err
	}
	f, err := os.OpenFile(getLogFilePathByDate(logFolder), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return &Logger{
			stdLog:  stdLogger,
			fileLog: nil,
			file:    nil,
		}, err
	}
	fileLogger := log.New(f, "", log.LstdFlags)
	return &Logger{
		stdLog:  stdLogger,
		fileLog: fileLogger,
		file:    f,
	}, nil
}

func (l *Logger) CloseFile() {
	if l.file != nil {
		l.file.Close()
	}
}

func (l *Logger) Println(v ...interface{}) {
	l.PrintlnToStd(v...)
	l.PrintlnToFile(v...)
}

func (l *Logger) PrintlnToStd(v ...interface{}) {
	if l.stdLog != nil {
		l.stdLog.Println(v...)
	}
}

func (l *Logger) PrintlnToFile(v ...interface{}) {
	if l.fileLog != nil {
		l.fileLog.Println(v...)
	}
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.FatallnToStd(v...)
	l.FatallnToFile(v...)
}

func (l *Logger) FatallnToStd(v ...interface{}) {
	if l.stdLog != nil {
		l.stdLog.Fatalln(v...)
	}
}

func (l *Logger) FatallnToFile(v ...interface{}) {
	if l.fileLog != nil {
		l.fileLog.Fatalln(v...)
	}
}
