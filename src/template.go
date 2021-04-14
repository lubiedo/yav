package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/lubiedo/yav/src/utils"
)

/* default definitions */
const (
	templdir    = "template"
	templateext = ".html" /* load files with this extension only */
)

/* extra template functions */
var extrafmap = map[string]interface{}{
	"dirls":    TplFuncDirListing,
	"toTplExt": TplFuncToTplExt,
}

func InitTemplate() (tpls *template.Template) {
	if !utils.FileExist(templdir) {
		config.Log.Fatal("Directory \"%s\" does not exist", templdir)
	}

	/* add extra functions */
	fmap := sprig.FuncMap()
	for k, v := range extrafmap {
		fmap[k] = v
	}

	tpls = template.Must(template.New("templates").Funcs(fmap).ParseGlob(templdir + "/*" + templateext))
	if config.Verbose {
		config.Log.Info("Templates loaded%s", tpls.DefinedTemplates())
	}
	return
}

func ToTemplateExt(s string) string {
	dir, name := filepath.Split(s)
	return filepath.Join(dir, strings.Replace(name, markdownext, templateext, -1))
}

func FromTemplateExt(s string) string {
	dir, name := filepath.Split(s)
	if name == "" {
		return ""
	}

	return filepath.Join(dir, strings.Replace(name, templateext, markdownext, -1))
}

func errorTplFunc(name string) { config.Log.Error("Error executing template function \"s\"", name) }
func TplFuncDirListing(p string) (files []string) {
	fullpath := sitedir + "/" + p

	err := filepath.Walk(fullpath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path != fullpath {
				files = append(files, strings.Replace(path, "site/", "", 1))
				if info.IsDir() {
					return filepath.SkipDir
				}
			}
			return nil
		})
	if err != nil {
		errorTplFunc("dirls")
	}
	return
}
func TplFuncToTplExt(p string) string { return ToTemplateExt(p) }
