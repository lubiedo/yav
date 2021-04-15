package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/lubiedo/yav/src/models"
	"github.com/lubiedo/yav/src/utils"
	"gopkg.in/yaml.v2"
)

/* default definitions */
const (
	sitedir     = "site"
	markdownext = ".md" /* parse files with this extension only */
)

var renderer *html.Renderer

func InitMarkdown() (p []models.SiteFile) {
	if !utils.FileExist(sitedir) {
		config.Log.Fatal("Directory \"%s\" does not exist", sitedir)
	}

	/* render for HTML */
	renderer = html.NewRenderer(html.RendererOptions{
		Flags:     html.CommonFlags | html.NoreferrerLinks | html.LazyLoadImages,
		Generator: name,
	})

	/* walk on site dir and load site source */
	err := filepath.Walk(sitedir, func(path string, info fs.FileInfo, err error) (e error) {
		if info.IsDir() {
			return nil
		}

		page, e := ProcessSiteFile(path)
		if e != nil {
			return e
		}

		if config.Verbose {
			config.Log.Info("Adding file: %s/%s", page.FileDir, page.FileName)
		}
		p = append(p, page)
		return nil
	})
	if err != nil {
		config.Log.Fatal("%s", err)
	}

	return
}

func ProcessSiteFile(path string) (page models.SiteFile, err error) {
	page = models.SiteFile{ /* defaults */
		FileName:   filepath.Base(path),
		FileDir:    filepath.Dir(path),
		IsMarkdown: false,
		Attrs: models.SiteFileAttr{
			Render: true,
		},
	}

	if len(page.FileName) > len(markdownext) &&
		page.FileName[len(page.FileName)-len(markdownext):] == markdownext {
		if page.Data, err = os.ReadFile(path); err != nil {
			return
		}
		page.IsMarkdown = true
		page.Checksum = utils.FileDataChecksum(page.Data)

		/* replace markdown ext for the template ext */
		urlpath := path[len(sitedir):]
		page.URLPath = urlpath[:len(urlpath)-len(markdownext)] + templateext

		/* parse each file's attributes */
		if config.Verbose {
			config.Log.Info("Parsing file: %s/%s", page.FileDir, page.FileName)
		}

		fm := GetFrontMatter(page.Data)
		if len(fm) == 0 {
			err = fmt.Errorf("Invalid or not found front matter (%s/%s)",
				page.FileDir, page.FileName)
			return
		}

		/*
		   TODO: check the templates actually exist.
		*/
		if err = yaml.Unmarshal(fm, &page.Attrs); err != nil {
			err = fmt.Errorf("%s (%s/%s)", err, page.FileDir, page.FileName)
			return
		}
		if page.Attrs.Template == "" {
			err = fmt.Errorf("Template missing for %s/%s", page.FileDir, page.FileName)
		}

		mdextensions := parser.CommonExtensions | parser.Includes
		mdparser := parser.NewWithExtensions(mdextensions)
		/* size(attrs + (delim*2) + (nl*2) */
		if page.Attrs.Render {
			page.Rendered = markdown.ToHTML(page.Data[len(fm)+8:], mdparser, renderer)
		} else {
			page.Rendered = page.Data[len(fm)+8:]
		}
	} else {
		/* common file */
		page.URLPath = path[len(sitedir):]
	}

	page.MimeType = utils.FileMimeType(page.URLPath)

	return
}

func GetFrontMatter(buf []byte) []byte {
	delim := []byte("---")
	if !utils.CompareSlices(buf[0:3], delim) {
		return []byte{}
	}

	for n := range buf {
		if n <= 3 {
			continue
		}
		if (len(buf) - n) <= 3 {
			break
		}
		if utils.CompareSlices(buf[n:n+3], delim) {
			return buf[3 : n-1]
		}
	}
	return []byte{}
}

func UpdateSiteFile(oldf models.SiteFile) (newfile models.SiteFile, e error) {
	for n, f := range files {
		if f.Checksum != oldf.Checksum {
			continue
		}

		newfile, e = ProcessSiteFile(GetSiteFilePath(f))
		if e != nil {
			break
		}
		files[n] = newfile
		break
	}
	return
}

func RemoveSiteFile(s models.SiteFile) {
	newfiles := make([]models.SiteFile, len(files))
	for n, f := range files {
		if (GetSiteFilePath(f) == GetSiteFilePath(s)) && (f.Checksum == s.Checksum) {
			continue
		}
		newfiles[n] = f
	}
	files = newfiles
}

func GetSiteFilePath(f models.SiteFile) string {
	return f.FileDir + "/" + f.FileName
}
