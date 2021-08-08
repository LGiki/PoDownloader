package logger

import (
	"PoDownloader/util"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger contains two logger, one print log to stderr
// and one write log to a file
type Logger struct {
	stdLog  *log.Logger
	fileLog *log.Logger
	file    *os.File
}

// getLogFileNameByDate returns the log file name with date
func getLogFileNameByDate() string {
	return fmt.Sprintf("%s.log", time.Now().Format("20060102"))
}

// getLogFilePathByDate returns a string that concatenates
// the specified log folder path and the log file name with date
func getLogFilePathByDate(logFolder string) string {
	return filepath.Join(logFolder, getLogFileNameByDate())
}

// NewLogger initializes and returns a logger instance
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

// CloseFile closes the file descriptor
func (l *Logger) CloseFile() {
	if l.file != nil {
		l.file.Close()
	}
}

// Println prints the log message to both stderr and log file
func (l *Logger) Println(v ...interface{}) {
	l.PrintlnToStd(v...)
	l.PrintlnToFile(v...)
}

// PrintlnToStd prints the log message to stderr
func (l *Logger) PrintlnToStd(v ...interface{}) {
	if l.stdLog != nil {
		l.stdLog.Println(v...)
	}
}

// PrintlnToFile prints the log message to log file
func (l *Logger) PrintlnToFile(v ...interface{}) {
	if l.fileLog != nil {
		l.fileLog.Println(v...)
	}
}

// Fatalln prints the log message to both stderr and log file then calls os.Exit(1)
func (l *Logger) Fatalln(v ...interface{}) {
	l.FatallnToStd(v...)
	l.FatallnToFile(v...)
}

// FatallnToStd prints the log message to stderr then calls os.Exit(1)
func (l *Logger) FatallnToStd(v ...interface{}) {
	if l.stdLog != nil {
		l.stdLog.Fatalln(v...)
	}
}

// FatallnToFile prints the log message to log file then calls os.Exit(1)
func (l *Logger) FatallnToFile(v ...interface{}) {
	if l.fileLog != nil {
		l.fileLog.Fatalln(v...)
	}
}
