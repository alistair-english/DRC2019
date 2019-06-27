package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	l    *logger
	once sync.Once
)

const (
	STD_LOG_PATH = "/src/github.com/alistair-english/DRC2019/logs/std.log"
)

// LoggerType is the type of the logger
type LoggerType int

type logger struct {
	stream io.Writer
	title  string
	mu     sync.RWMutex
}

func (l *logger) Init() {
	// Inits on a default stream :)
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(os.Getenv("GOPATH")+STD_LOG_PATH, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	l.stream = f
	l.title = "DEFAULT_LOG"
}

func (l *logger) SetMode(stream io.Writer, title string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.stream = stream
	l.title = title
	fmt.Printf("Changing Log Stream to %v", title)
	l.Log("---Changing Log Mode To %v ---", title)
}

func (l *logger) GetMode() (io.Writer, string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.stream, l.title
}

func (l *logger) Log(format string, v ...interface{}) {
	// I would like to replace fmt.Sprintf with custom function but this will do
	l.stream.Write([]byte(fmt.Sprintf(format, v...)))
}

func (l *logger) Logln(format string, v ...interface{}) {
	// I would like to replace fmt.Sprintf with custom function but this will do
	l.stream.Write([]byte(fmt.Sprintf(format+"\n", v...)))
}

// Logger creates a new logger, if one doesnt already exist
func Logger() *logger {
	once.Do(func() {
		l = &logger{}
	})
	return l
}
