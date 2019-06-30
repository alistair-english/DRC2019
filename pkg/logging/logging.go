package logging

import (
	"fmt"
	"github.com/alistair-english/DRC2019/pkg/services/seriallogservice"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
	"sync"
)

var (
	l    *logger
	once sync.Once
)

const (
	stdLog = "default.log"
	path   = "/src/github.com/alistair-english/DRC2019/logs/"
)

// Flags for logging
const (
	LogFile   = 1
	LogStd    = 2
	LogSerial = 4
	All       = LogFile | LogStd | LogSerial
)

const TAG = ""

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
	serLog     *seriallogservice.SerialLogService
}

func (l *logger) Init() {
	// Inits on a default stream :)
	l.AddStream(stdLog, "DEFAULT_LOG")
	l.cTitle = l.streamList[0].title
	l.cStream = l.streamList[0].stream
	l.cIndex = 0
	l.ChangeStream("DEFAULT_LOG")
}

func (l *logger) AddSerialLogService(serLog *seriallogservice.SerialLogService) {
	l.serLog = serLog
	l.Logln(TAG, All, "Got Serial Log Service!")
}

func (l *logger) AddStream(fileName string, title string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(os.Getenv("GOPATH")+path+fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	l.streamList = append(l.streamList, Stream{f, title})
}

func (l *logger) ChangeStream(title string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logln(TAG, LogFile, "---Changing Log Stream: [%v -> %v]---", l.cTitle, title)
	for i := range l.streamList {
		if l.streamList[i].title == title {
			l.cStream = l.streamList[i].stream
			log.SetOutput(l.cStream)
			if title != l.cTitle {
				l.Logln(TAG, LogFile, "---Changing Log Stream: [%v -> %v]---", l.cTitle, title)
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
	fmt.Printf("Stream List:\n[\n")
	for i := range l.streamList {
		fmt.Printf("%v\n", l.streamList[i].title)
	}
	fmt.Printf("]\n")
}

func (l *logger) Log(tag string, flags int, format string, v ...interface{}) {
	// Make some colors
	tagColour := color.New(color.FgCyan).Add(color.Bold)
	logColour := color.New(color.FgWhite)
	// I would like to replace fmt.Sprintf with custom function but this will do
	if tag != "" {
		tag = tag + ": "
	}
	if (flags & LogFile) != 0 {
		// l.cStream.Write([]byte(fmt.Sprintf(tag+format, v...)))
		log.Printf(tag+format, v...)
	}
	if (flags & LogStd) != 0 {
		// fmt.Printf(tag+format, v...)
		tagColour.Print(tag)
		logColour.Printf(format, v...)
	}
	if (flags&LogSerial) != 0 && l.serLog != nil {
		l.serLog.LogToSerial(fmt.Sprintf(tag+format, v...))
	}
}

func (l *logger) Logln(tag string, flags int, format string, v ...interface{}) {
	// Make some colors
	tagColour := color.New(color.FgCyan).Add(color.Bold)
	logColour := color.New(color.FgWhite)
	// I would like to replace fmt.Sprintf with custom function but this will do
	if tag != "" {
		tag = tag + ": "
	}
	if (flags & LogFile) != 0 {
		// l.cStream.Write([]byte(fmt.Sprintf(tag+format+"\n", v...)))
		log.Printf(tag+format+"\n", v...)
	}
	if (flags & LogStd) != 0 {
		// fmt.Printf(tag+format+"\n", v...)
		tagColour.Print(tag)
		logColour.Printf(format+"\n", v...)
	}
	if (flags&LogSerial) != 0 && l.serLog != nil {
		l.serLog.LogToSerial(fmt.Sprintf(tag+format+"\n", v...))
	}
}

// Logger creates a new logger, if one doesnt already exist
func Logger() *logger {
	once.Do(func() {
		l = &logger{}
		l.streamList = []Stream{}
	})
	return l
}

// L is an alias to Logger to make more concise calls to logger
func L() *logger {
	return Logger()
}
