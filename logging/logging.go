package logging

import (
	//"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"git.dsrt-int.net/actionmc/actionmc-site-go/config"
)

// alex moment

var site_configuration *config.Configuration = config.InitialiseConfig()

type Logger struct {
	Debug *log.Logger
	Log   *log.Logger
	Warn  *log.Logger
	Err   *log.Logger
	Fatal *log.Logger
}

func New() *Logger {
	var result Logger

	logFile, err := os.OpenFile(filepath.Join(config.GetDataDir(), "log_stdout.txt"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	errLogFile, err := os.OpenFile(filepath.Join(config.GetDataDir(), "log_stderr.txt"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	mw_err := io.MultiWriter(os.Stderr, errLogFile)

	result.Debug = log.New(mw, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	if site_configuration.IsPreproduction {
		result.Debug = log.New(mw, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		result.Debug = log.New(io.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	result.Log = log.New(mw, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	result.Warn = log.New(mw, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	result.Err = log.New(mw_err, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	result.Fatal = log.New(mw_err, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &result
}
