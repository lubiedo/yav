package models

import "github.com/lubiedo/yav/src/utils"

type Config struct {
	Port    string
	Addr    string
	Verbose bool

	UseHTTPS bool
	CertPath string
	KeyPath  string

	LogFile     string
	Log         *utils.Log
	TplVars     string
	TplErroPage string
}
