package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"

	"net/http"
)

var prefix string = fmt.Sprintf("(%d) Y∆V ", os.Getpid())

type Log struct {
	Logger *log.Logger
	OutFD  *os.File
}

func NewLog(out io.Writer) *Log {
	return &Log{
		Logger: log.New(out, prefix, log.Ldate|log.Lmicroseconds),
	}
}

func (L *Log) Fatal(s string, args ...interface{}) { L.Logger.Fatalf("[FATAL] "+s, args...) }
func (L *Log) Panic(s string, args ...interface{}) { L.Logger.Panicf("[PANIC] "+s, args...) }
func (L *Log) Error(s string, args ...interface{}) {
	L.Logger.Printf("[ERROR] "+s, args...)
	debug.PrintStack()
}
func (L *Log) Info(s string, args ...interface{}) { L.Logger.Printf("[INFO] "+s, args...) }
func (L *Log) Access(r *http.Request) {
	L.Logger.Printf("[REQ] %s - %s %s %s", r.RemoteAddr, r.Method, r.URL, r.Proto)
}
