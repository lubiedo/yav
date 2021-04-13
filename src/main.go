package main

import (
	"flag"
	"html/template"
	"os"

	"github.com/lubiedo/yav/src/models"
	"github.com/lubiedo/yav/src/utils"
)

/* defaults */
const (
	verbosity = false
	usehttps  = true
	port      = "8080"
	addr      = "localhost"

	name    = "Y∆V"
	version = "0.0.1"
)

var (
	config models.Config
	files  []models.SiteFile
	tpls   *template.Template
)

func init() {
	/* cli flags */
	flag.BoolVar(&config.Verbose, "verbose", verbosity, "Be verbose.")
	flag.StringVar(&config.Port, "port", port, "Port number.")
	flag.StringVar(&config.Addr, "addr", addr, "Address to listen.")

	flag.BoolVar(&config.UseHTTPS, "use-https", usehttps,
		"Secure connection via HTTPS.")
	flag.StringVar(&config.CertPath, "cert", "", "Certificate file path.")
	flag.StringVar(&config.KeyPath, "key", "", "Key file path.")

	flag.StringVar(&config.LogFile, "log", "", "Output to log file.")

	flag.Parse()
}

func main() {
	/* logging */
	config.Log = utils.NewLog(os.Stderr)
	if config.LogFile != "" {
		fd, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)

		if err != nil {
			config.Log.Fatal("Unable to open file at \"%s\"", config.LogFile)
		}
		config.Log.OutFD = fd
		defer config.Log.OutFD.Close()

		config.Log.Logger.SetOutput(config.Log.OutFD)
	}

	if config.UseHTTPS {
		if !utils.FileExist(config.CertPath) || !utils.FileExist(config.KeyPath) {
			config.Log.Fatal("Can't load certificate or key files")
		}
	}

	/* process site data */
	if config.Verbose {
		config.Log.Info("Parsing site source")
	}
	files = InitMarkdown()
	if config.Verbose {
		config.Log.Info("Loading templates")
	}
	tpls = InitTemplate()
	Serve()
}
