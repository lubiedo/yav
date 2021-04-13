package main

import (
	"html/template"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/lubiedo/yav/src/utils"
)

/* default definitions */
const (
	templdir    = "template"
	templateext = ".html" /* load files with this extension only */
)

func InitTemplate() (tpls *template.Template) {
	if !utils.FileExist(templdir) {
		config.Log.Fatal("Directory \"%s\" does not exist", templdir)
	}

	tpls = template.Must(template.New("templates").Funcs(sprig.FuncMap()).ParseGlob(templdir + "/*" + templateext))
	if config.Verbose {
		config.Log.Info("Templates loaded%s", tpls.DefinedTemplates())
	}
	return
}

func RevertTemplateExt(s string) string {
	pos := strings.LastIndex(s, "/")
	if pos == -1 || len(s[pos:]) == 1 {
		return ""
	}

	return strings.Replace(s, templateext, markdownext, -1)
}
