package logging

import (
	"log"
	"io"
	"io/ioutil"
	"os"
)

var (
	Trace       *log.Logger
	Info        *log.Logger
	Query       *log.Logger
	Warning     *log.Logger
	Error       *log.Logger
	Integration *log.Logger
	Consuming   *log.Logger
	Socket      *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Integration = log.New(infoHandle,
		"INTEGRATION: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Query = log.New(infoHandle,
		"QUERY: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Consuming = log.New(infoHandle,
		"AMQP: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Socket = log.New(infoHandle,
		"SOCKET: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

}

func init() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}
