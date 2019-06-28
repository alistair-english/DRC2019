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
	STD_LOG = "std.log"
	PATH    = "/src/github.com/alistair-english/DRC2019/logs/"
)

// Stream contains a stream and title
type Stream struct {
	stream io.Writer
	title  string
}

type logger struct {
	cStream    io.Writer
	cTitle     string
	cIndex     int
	streamList []Stream
	mu         sync.RWMutex
}

func (l *logger) Init() {
	// Inits on a default stream :)
	l.AddStream(STD_LOG, "DEFAULT_LOG")
	l.cTitle = l.streamList[0].title
	l.cStream = l.streamList[0].stream
	l.cIndex = 0
	l.ChangeStream("DEFAULT_LOG")
}

func (l *logger) AddStream(fileName string, title string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(os.Getenv("GOPATH")+PATH+fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	l.streamList = append(l.streamList, Stream{f, title})
}

func (l *logger) ChangeStream(title string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Printf("Changing Log Stream: [%v -> %v] \n", l.cTitle, title)
	l.Logln("---Changing Log Stream: [%v -> %v]---", l.cTitle, title)
	for i := range l.streamList {
		if l.streamList[i].title == title {
			l.cStream = l.streamList[i].stream
			if title != l.cTitle {
				l.Logln("---Changing Log Stream: [%v -> %v]---", l.cTitle, title)
			}
			l.cTitle = l.streamList[i].title
			l.cIndex = i
		}
	}
}

func (l *logger) GetStream() Stream {
	l.mu.Lock()
	defer l.mu.Unlock()
	return Stream{l.cStream, l.cTitle}
}

func (l *logger) ListStreams() {
	fmt.Printf("Steam List:\n")
	for i := range l.streamList {
		fmt.Printf("%v", l.streamList[i])
	}
}

func (l *logger) Log(format string, v ...interface{}) {
	// I would like to replace fmt.Sprintf with custom function but this will do
	l.cStream.Write([]byte(fmt.Sprintf(format, v...)))
}

func (l *logger) Logln(format string, v ...interface{}) {
	// I would like to replace fmt.Sprintf with custom function but this will do
	l.cStream.Write([]byte(fmt.Sprintf(format+"\n", v...)))
}

func (l *logger) LogOut(format string, v ...interface{}) {
	l.cStream.Write([]byte(fmt.Sprintf(format, v...)))
	fmt.Printf(format, v...)
}

func (l *logger) LogOutln(format string, v ...interface{}) {
	l.cStream.Write([]byte(fmt.Sprintf(format+"\n", v...)))
	fmt.Printf(format+"\n", v...)
}

// Logger creates a new logger, if one doesnt already exist
func Logger() *logger {
	once.Do(func() {
		l = &logger{}
		l.streamList = []Stream{}
	})
	return l
}
