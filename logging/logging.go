package logging

import (
	"io"
	"log"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func InitLoggers(infoHandle io.Writer, errorHandle io.Writer) {

	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(infoHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

}
