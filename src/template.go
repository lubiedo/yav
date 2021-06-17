package main

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/lubiedo/yav/src/utils"
	"gopkg.in/yaml.v2"
)

/* default definitions */
const (
	templdir    = "template"
	templateext = ".html" /* load files with this extension only */
)

/* extra template functions */
var extrafmap = map[string]interface{}{
	"readFile":   TplFuncReadFile,
	"dirLs":      TplFuncDirLs,
	"toTplExt":   TplFuncToTplExt,
	"fromTplExt": TplFuncFromTplExt,
	"getFM":      TplFuncGetFM,
	"getSiteMap": TplFuncGetSM,
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

	if config.TplErroPage != "" {
		if t := tpls.Lookup(config.TplErroPage); t == nil {
			config.Log.Fatal("Template \"%s\" does not exist", config.TplErroPage)
		}
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

func errorTplFunc(name string) { config.Log.Error("Error executing template function \"%s\"", name) }
func TplFuncReadFile(p string) (data []byte) {
	fullpath := filepath.Join(sitedir, p)
	if !utils.FileExist(fullpath) || utils.FileIsDir(fullpath) {
		return
	}
	data, err := os.ReadFile(fullpath)
	if err != nil {
		errorTplFunc("readFile")
	}
	return
}
func TplFuncDirLs(p string) (files []string) {
	fullpath := sitedir + "/" + p

	err := filepath.Walk(fullpath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path != fullpath {
				files = append(files, strings.Replace(path, sitedir, "", 1))
				if info.IsDir() {
					return filepath.SkipDir
				}
			}
			return nil
		})
	if err != nil {
		errorTplFunc("dirLs")
	}
	return
}
func TplFuncToTplExt(p string) string   { return ToTemplateExt(p) }
func TplFuncFromTplExt(p string) string { return FromTemplateExt(p) }
func TplFuncGetFM(p string) (fm map[string]interface{}) {
	data := TplFuncReadFile(p)
	if len(data) == 0 {
		return
	}

	fmdata := GetFrontMatter(data)
	if err := yaml.Unmarshal(fmdata, &fm); err != nil {
		errorTplFunc("getFM")
	}
	return
}
func TplFuncGetSM() string {
	data, err := sitemap.ToJSON()
	if err != nil {
		config.Log.Error("Error executing template function \"%s\" (%s)", name, err)
		return ""
	}
	return string(data)
}
