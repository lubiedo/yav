package utils

import (
	"encoding/json"
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

type LogRequest struct {
	RemoteAddr    string `json:"remote_addr"`
	RemoteUA      string `json:"remote_ua"`
	RequestMethod string `json:"request_method"`
	RequestProto  string `json:"request_proto"`
	RequestURI    string `json:"request_uri"`
}

func ToLogRequest(r *http.Request) *LogRequest {
	return &LogRequest{
		RemoteAddr:    r.RemoteAddr,
		RemoteUA:      r.Header.Get("User-Agent"),
		RequestMethod: r.Method,
		RequestProto:  r.Proto,
		RequestURI:    r.RequestURI,
	}
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
	d, err := json.Marshal(ToLogRequest(r))
	if err != nil {
		L.Error("%s", err)
		return
	}
	L.Logger.Printf("[REQ] %s", d)
}
