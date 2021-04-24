package main

import (
	"flag"
	"html/template"
	"os"
	"os/signal"
	"syscall"

	"github.com/lubiedo/yav/src/models"
	"github.com/lubiedo/yav/src/utils"
	"gopkg.in/yaml.v2"
)

/* defaults */
const (
	verbosity = false
	usehttps  = true
	port      = "8080"
	addr      = "localhost"

	name    = "Y∆V"
	version = "0.1.0"
)

var (
	config  models.Config
	files   models.Sites
	headers models.Headers
	tpls    *template.Template
	tplvars map[string]interface{}
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
	flag.Var(&headers, "header", "Add HTTP header.")

	flag.StringVar(&config.LogFile, "log", "", "Output to log file.")
	flag.StringVar(&config.TplErroPage, "tpl-error", "", "Use template as error page (filename).")
	flag.StringVar(&config.TplVars, "tpl-vars", "", "Load template variables.")

	flag.Parse()

	config.WriteTimeOut = 10
	config.IdleTimeOut = 20
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

	if config.TplVars != "" {
		data, err := os.ReadFile(config.TplVars)
		if err != nil {
			config.Log.Fatal("Error loading template variables from \"%s\"",
				config.TplVars)
		}

		err = yaml.Unmarshal(data, &tplvars)
		if err != nil {
			config.Log.Fatal("Parsing template variables from \"%s\" failed",
				config.TplVars)
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

	/* define USR1 signal catch for template reload */
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGUSR1)

		for {
			_ = <-sig
			tpls = InitTemplate()
			if config.Verbose {
				config.Log.Info("Templates were reloaded.")
			}
		}
	}()

	Serve()
}
