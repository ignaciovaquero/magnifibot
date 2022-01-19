package controller

import (
	"fmt"
	"log"
)

// Logger defines an interface for logs
type Logger interface {
	// logging facilities
	Debug(...interface{})
	Debugf(string, ...interface{})
	Debugw(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Errorw(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Infow(string, ...interface{})
}

// DefaultLogger is a struct that implements the Logger interface in a very
// basic way
type DefaultLogger struct{}

// Debug prints a message with the Debug level
func (d *DefaultLogger) Debug(v ...interface{}) {
	log.Print(v...)
}

// Debugf prints a formatted message with the Debug level
func (d *DefaultLogger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Debugw prints a message with values, with the Debug level
func (d *DefaultLogger) Debugw(base string, keysAndValues ...interface{}) {
	msg := base
	if msg == "" && len(keysAndValues) > 0 {
		msg = fmt.Sprint(keysAndValues...)
	} else if msg != "" && len(keysAndValues) > 0 {
		msg = fmt.Sprintf(base, keysAndValues...)
	}
	fmt.Println(msg)
}

// Error prints a message with the Error level
func (d *DefaultLogger) Error(v ...interface{}) {
	log.Print(v...)
}

// Errorf prints a formatted message with the Error level
func (d *DefaultLogger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Errorw prints a message with values, with the Error level
func (d *DefaultLogger) Errorw(base string, keysAndValues ...interface{}) {
	d.Debugw(base, keysAndValues...)
}

// Info prints a message with the Info level
func (d *DefaultLogger) Info(v ...interface{}) {
	log.Print(v...)
}

// Infof prints a formatted message with the Info level
func (d *DefaultLogger) Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// Infow prints a message with values, with the Info level
func (d *DefaultLogger) Infow(base string, keysAndValues ...interface{}) {
	d.Debugw(base, keysAndValues...)
}
